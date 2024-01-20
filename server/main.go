package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"server/main/analyzer"

	"github.com/gorilla/mux"
)

func main() {
	webDir := "./web"
	var port int
	var serveUI bool
	flag.BoolVar(&serveUI, "ui", true, "Serve the UI or not?")
	flag.IntVar(&port, "port", 8080, "Port for the service")
	flag.Parse()

	r := mux.NewRouter()

	// register routes
	r.HandleFunc("/api/health", HealthEndPoint).Methods("GET")
	r.HandleFunc("/api/analyze", PostAnalyze).Methods("POST")

	// serve UI?
	if serveUI {
		r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(webDir))))
	} else {
		log.Println("[WARN] UI is disabled! Only API endpoint is exposed.")
	}

	// start server
	log.Printf("Service is listening on :%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
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
