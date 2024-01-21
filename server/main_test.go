package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestAnalyze(t *testing.T) {
	// w := httptest.NewRecorder()
	// r := mux.NewRouter()

	// AnalyzeEndPoint("/api").Register(r)
	// r.ServeHTTP(w, httptest.NewRequest("POST", "/api/analyze", nil))

	// if w.Code != http.StatusBadRequest {
	// 	t.Error("Expected to have status code 400! Actual:", w.Code)
	// }
}
