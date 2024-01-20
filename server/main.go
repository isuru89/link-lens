package main

import (
	"encoding/json"
	"log"
	"net/http"
	"server/main/analyzer"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/health", HealthEndPoint).Methods("GET")
	r.HandleFunc("/api/analyze", PostAnalyze).Methods("POST")
	//analyzer.AnalyzeUrl("https://www.w3.org/TR/1998/REC-CSS2-19980512/sample.html")

	err := http.ListenAndServe(":8090", r)
	if err != nil {
		log.Fatalf("Error occurred while loading server: %v", err)
	}
}

func HealthEndPoint(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func PostAnalyze(w http.ResponseWriter, r *http.Request) {
	var req AnalyzeRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		log.Fatalf("Error decoding request %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := analyzer.AnalyzeUrl(req.Url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errObj, _ := json.Marshal(err)
		w.Write([]byte(errObj))
		return
	}

	content, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
