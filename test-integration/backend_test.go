package integration_test

import (
	"context"
	"net/http"
	"path/filepath"
	"phasor-backend/testutil"
	"runtime"
	"testing"

	"github.com/monkescience/testastic"
)

func TestBackendInstanceAPI(t *testing.T) {
	t.Parallel()

	t.Run("returns instance info with version", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a backend server with version 1.2.3
		server := testutil.NewTestServer("1.2.3", testutil.NewTestLogger(t))
		defer server.Close()

		// WHEN: requesting instance info
		resp := httpGet(t, server.URL+"/instance/info")
		defer resp.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: response matches expected JSON structure
		testastic.Equal(t, http.StatusOK, resp.StatusCode)
		testastic.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		testastic.AssertJSON(t, testdataPath("backend_instance_info", "expected_response.json"), resp.Body)
	})

	t.Run("returns consistent hostname across requests", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a backend server
		server := testutil.NewTestServer("1.0.0", testutil.NewTestLogger(t))
		defer server.Close()

		// WHEN: requesting instance info twice
		resp1 := httpGet(t, server.URL+"/instance/info")
		defer resp1.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		resp2 := httpGet(t, server.URL+"/instance/info")
		defer resp2.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: both responses match expected JSON structure
		testastic.AssertJSON(t, testdataPath("backend_consistent_hostname", "expected_response.json"), resp1.Body)
		testastic.AssertJSON(t, testdataPath("backend_consistent_hostname", "expected_response.json"), resp2.Body)
	})

	t.Run("health live endpoint responds OK", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a backend server
		server := testutil.NewTestServer("test-version", testutil.NewTestLogger(t))
		defer server.Close()

		// WHEN: requesting the live health endpoint
		resp := httpGet(t, server.URL+"/health/live")
		defer resp.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: response matches expected JSON structure
		testastic.Equal(t, http.StatusOK, resp.StatusCode)
		testastic.AssertJSON(t, testdataPath("backend_health_live", "expected_response.json"), resp.Body)
	})

	t.Run("health ready endpoint responds OK", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a backend server
		server := testutil.NewTestServer("test-version", testutil.NewTestLogger(t))
		defer server.Close()

		// WHEN: requesting the ready health endpoint
		resp := httpGet(t, server.URL+"/health/ready")
		defer resp.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: response matches expected JSON structure
		testastic.Equal(t, http.StatusOK, resp.StatusCode)
		testastic.AssertJSON(t, testdataPath("backend_health_ready", "expected_response.json"), resp.Body)
	})
}

// testdataPath returns the path to a testdata file for the given test case.
//
//nolint:unparam // filename varies in different test scenarios.
func testdataPath(testcase, filename string) string {
	//nolint:dogsled // runtime.Caller returns 4 values, we only need filename.
	_, callerFile, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(callerFile), "testdata", testcase, filename)
}

// httpGet performs an HTTP GET request with context.
func httpGet(t *testing.T, url string) *http.Response {
	t.Helper()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	testastic.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	testastic.NoError(t, err)

	return resp
}
