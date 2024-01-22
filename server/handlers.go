package server

import (
	"encoding/json"
	"fmt"
	"linklens/analyzer"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type RouteHandler struct {
	RouteDef func(r *mux.Route) string

	Handler http.HandlerFunc
}

func (h RouteHandler) Register(r *mux.Router) {
	ep := h.RouteDef(r.NewRoute().HandlerFunc(h.Handler))
	log.Printf(" - %s", ep)
}

func HealthEndPoint(contextPath string) RouteHandler {
	return RouteHandler{
		RouteDef: func(r *mux.Route) string {
			r.Path(contextPath + "/health").Methods("GET")
			return fmt.Sprintf("%s: %s%s", "GET", contextPath, "/health")
		},
		Handler: func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(HealthResponse{Alive: true})
		},
	}
}

func AnalyzeEndPoint(contextPath string) RouteHandler {
	return RouteHandler{
		RouteDef: func(r *mux.Route) string {
			r.Path(contextPath + "/analyze").Methods("POST")
			return fmt.Sprintf("%s: %s%s", "POST", contextPath, "/analyze")
		},
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req AnalyzeRequest
			err := json.NewDecoder(r.Body).Decode(&req)

			if err != nil {
				log.Printf("[ERROR] Error decoding request: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else if req.Url == "" {
				log.Println("[ERROR] Analyze URL cannot be empty! Url!")
				http.Error(w, "Empty URL", http.StatusBadRequest)
				return
			}

			result, err := analyzer.AnalyzeUrl(req.Url)
			if err != nil {
				handleAnalysisError(err, w)
				return
			}

			content, _ := json.Marshal(result)
			w.WriteHeader(http.StatusOK)
			w.Write(content)
		},
	}
}

func handleAnalysisError(err error, w http.ResponseWriter) {
	log.Println("[ERROR] " + err.Error())

	w.WriteHeader(http.StatusInternalServerError)
	var errObj []byte
	e, ok := err.(*analyzer.AnalysisError)

	if ok {
		errObj, _ = json.Marshal(map[string]interface{}{
			"errorCode": e.ErrorCode,
			"message":   e.Cause.Error(),
		})
	} else {
		errObj, _ = json.Marshal(err)
	}
	w.Write([]byte(errObj))
}