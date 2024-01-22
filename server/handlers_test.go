package server

import (
	"encoding/json"
	"linklens/analyzer"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestHealth(t *testing.T) {
	w := httptest.NewRecorder()
	r := mux.NewRouter()

	HealthEndPoint("/api").Register(r)
	r.ServeHTTP(w, httptest.NewRequest("GET", "/api/health", nil))

	if w.Code != http.StatusOK {
		t.Error("Did not expected to fail the health end point! Actual:", w.Code)
	}
	var heathRes HealthResponse
	err := json.NewDecoder(w.Body).Decode(&heathRes)
	if err != nil {
		t.Error("Expected to have a valid response!")
	}
	if heathRes.Alive != true {
		t.Error("Expected to return alive=true, but got", heathRes.Alive)
	}
}

func TestAnalyze_BadRequests(t *testing.T) {
	// GIVEN
	r := mux.NewRouter()
	AnalyzeEndPoint("/api").Register(r)

	t.Run("When Malformed Json", func(t *testing.T) {
		// WHEN
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("POST", "/api/analyze", strings.NewReader("{")))

		// THEN
		if w.Code != http.StatusBadRequest {
			t.Error("Expected to have status code 400! Actual:", w.Code, w.Body)
		}
	})

	t.Run("When No Request Body", func(t *testing.T) {
		// WHEN
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("POST", "/api/analyze", nil))

		// THEN
		if w.Code != http.StatusBadRequest {
			t.Error("Expected to have status code 400! Actual:", w.Code, w.Body)
		}
	})

	t.Run("Json, but No URL", func(t *testing.T) {
		// WHEN
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("POST", "/api/analyze", strings.NewReader("{}")))

		// THEN
		if w.Code != http.StatusBadRequest {
			t.Error("Expected to have status code 400! Actual:", w.Code, w.Body)
		}
	})

}

func TestAnalyze_400_InvalidURL(t *testing.T) {
	// GIVEN
	r := mux.NewRouter()
	AnalyzeEndPoint("/api").Register(r)

	t.Run("Invalid URL", func(t *testing.T) {
		// WHEN
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("POST", "/api/analyze", strings.NewReader(`{ "url": ":this is an invalid url" }`)))

		// THEN
		if w.Code != http.StatusInternalServerError {
			t.Error("Expected to have status code 500! Actual:", w.Code, w.Body)
		}
		var errObj analyzer.AnalysisError
		if err := json.NewDecoder(w.Body).Decode(&errObj); err != nil {
			t.Errorf("Expected to return an AnlaysisError object! Received: %s", err.Error())
		}
		if errObj.ErrorCode != analyzer.ErrorInvalidUrl {
			t.Errorf("Expected to return invalidUrl error code, but got %s", errObj.ErrorCode)
		}
	})
}
