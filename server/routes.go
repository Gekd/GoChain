package server

import (
	"log"
	"net/http"
)

// All routes of the server
func addRoutes(mux *http.ServeMux, logger *log.Logger) {
	mux.Handle("GET /chain", handleGetChain(logger))
	mux.Handle("POST /add", handleAddBlock(logger))
}
