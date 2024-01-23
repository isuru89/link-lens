package main

import (
	"flag"
	"fmt"
	"linklens/server"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	var webDir string
	var port int
	var serveUI bool
	flag.BoolVar(&serveUI, "ui", true, "Serve the UI or not?")
	flag.StringVar(&webDir, "webDir", "./web/build", "Directory path to the web artifacts")
	flag.IntVar(&port, "port", 8080, "Port for the service")
	flag.Parse()

	r := mux.NewRouter()

	slog.Info("Registering end points:")
	contextPath := "/api"
	// register routes
	server.HealthEndPoint(contextPath).Register(r)
	server.AnalyzeEndPoint(contextPath).Register(r)

	// serve UI?
	if serveUI {
		r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(webDir))))
	} else {
		slog.Warn("UI is disabled! Only API endpoint is exposed.")
	}

	// start server
	slog.Info("Service is listening on ", "port", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		slog.Error(fmt.Sprintf("Error occurred while loading server: %v", err))
	}
}
