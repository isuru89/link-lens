package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	webDir := "./web/build"
	var port int
	var serveUI bool
	flag.BoolVar(&serveUI, "ui", true, "Serve the UI or not?")
	flag.IntVar(&port, "port", 8080, "Port for the service")
	flag.Parse()

	r := mux.NewRouter()

	log.Println("Registering end points:")
	contextPath := "/api"
	// register routes
	HealthEndPoint(contextPath).Register(r)
	AnalyzeEndPoint(contextPath).Register(r)

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
