package server

import (
	"log"
	"net/http"
)

// All routes of the server
func addRoutes(mux *http.ServeMux, logger *log.Logger) {
	mux.Handle("GET /chain", checkIfNodeRecognised(logger)(handleGetChain(logger)))
	mux.Handle("GET /nodes", checkIfNodeRecognised(logger)(handleGetNodes(logger)))
	mux.Handle("POST /add", checkIfNodeRecognised(logger)(handleAddBlock(logger)))
	mux.Handle("POST /receive-block", checkIfNodeRecognised(logger)(handleBlockReceive(logger)))
}
