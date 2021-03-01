package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApp_rate(t *testing.T) {
	// 1. Initialise app
	a := InitialiseApp("redis://localhost:6379")

	// 2. Create recorder and request
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/suggestion/u1", nil)

	// 3. Server the request
	a.router.ServeHTTP(rw, req)

	// 4. Assert the status code and the body
	if rw.Result().StatusCode != http.StatusOK {
		t.Errorf("server responded with the wrong error code: %d", rw.Result().StatusCode)
	}
}
