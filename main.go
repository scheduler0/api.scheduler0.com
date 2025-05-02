package main

import (
	"fmt"
	httpLogger "github.com/go-http-utils/logger"
	"github.com/gorilla/mux"
	"log"
	"log/slog"
	"net/http"
	"os"
)

// Start this will start the http server
func main() {
	router := mux.NewRouter()

	router.PathPrefix("/api-docs/").Handler(http.StripPrefix("/api-docs/", http.FileServer(http.Dir("./api-docs/"))))

	slog.Info("Openapi-doc is running on", ":port", 3002)
	err := http.ListenAndServe(fmt.Sprintf(":%v", 3002), httpLogger.Handler(router, os.Stderr, httpLogger.CombineLoggerType))
	if err != nil {
		log.Fatal("failed to start openapi-doc-server", err)
	}
}
