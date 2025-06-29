package server

import (
	"log"
	"net/http"
)

// All routes of the server
func addRoutes(mux *http.ServeMux, logger *log.Logger) {
	mux.Handle("GET /chain", checkIfNodeRecognised(logger)(handleGetChain(logger)))
	mux.Handle("GET /nodes", handleGetNodes(logger))
	mux.Handle("POST /add", handleAddBlock(logger))
	mux.Handle("POST /receive-block", handleBlockReceive(logger))
}
