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

func TestAnalyze_Errors(t *testing.T) {
	testcases := map[string]struct {
		requestBody interface{}
		statusCode  int
	}{
		"With Malformed Json": {
			requestBody: "{",
			statusCode:  400,
		},
		"No Request Body": {
			requestBody: nil,
			statusCode:  400,
		},
		"Empty Request Body": {
			requestBody: "",
			statusCode:  400,
		},
		"Json, but No URL": {
			requestBody: "{}",
			statusCode:  400,
		},
	}

	// GIVEN
	r := mux.NewRouter()
	AnalyzeEndPoint("/api").Register(r)

	for name, test := range testcases {
		t.Run(name, func(t *testing.T) {
			// WHEN
			w := httptest.NewRecorder()
			if test.requestBody == nil {
				r.ServeHTTP(w, httptest.NewRequest("POST", "/api/analyze", nil))
			} else {
				r.ServeHTTP(w, httptest.NewRequest("POST", "/api/analyze", strings.NewReader(test.requestBody.(string))))
			}

			// THEN
			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected to have status code %d! Actual: %d", test.statusCode, w.Code)
			}
		})
	}
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

func TestAnalyze_200_Success(t *testing.T) {
	// GIVEN
	w := httptest.NewRecorder()
	r := mux.NewRouter()
	AnalyzeEndPoint("/api").Register(r)

	url := "https://www.google.com"
	// WHEN
	r.ServeHTTP(w, httptest.NewRequest("POST", "/api/analyze", strings.NewReader(`{ "url": "https://www.google.com" }`)))

	// THEN
	if w.Code != http.StatusOK {
		t.Error("Not expected to throw an error! Actual:", w.Code)
	}
	var res analyzer.AnalysisData
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Errorf("Expected to return an AnlaysisData object! Received: %s", err.Error())
	}

	if res.SourceUrl != url {
		t.Errorf("Expected SourceUrl to be %s, but got %s", url, res.SourceUrl)
	} else if res.HtmlVersion != "5" {
		t.Errorf("Expected HtmlVersion to be 5, but got %s", res.HtmlVersion)
	} else if res.Title == "" {
		t.Errorf("Expected Title to be non-empty, but got %s", res.Title)
	} else if res.PageType != analyzer.Unknown {
		t.Errorf("Expected Title to be Unknown, but got %s", res.PageType)
	} else if res.LinkStats.ExternalLinkCount <= 0 && res.LinkStats.InternalLinkCount <= 0 && res.LinkStats.InvalidLinkCount <= 0 {
		t.Error("Expected at least to have one link, but got all zeros")
	}
}
