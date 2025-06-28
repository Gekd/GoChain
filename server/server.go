package server

import (
	"GoChain/block"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

// NewServer initializes the HTTP multiplexer and attaches all routes.
// It returns an http.Handler to be passed into the server.
func NewServer(logger *log.Logger) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, logger)
	return mux
}

// Writes JSON response with the given status code and payload into ResponseWriter.
func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

// Parses JSON request body into a chosen value.
func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

// Returns entire blockchain as a JSON array.
// Route: GET /chain
func handleGetChain(logger *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Println("GET /chain")
			_ = encode(w, r, 200, block.GetBlockchain())
		},
	)
}

// Defines the JSON body for POST /add request
type AddBlockData struct {
	Data string `json:"data"`
}

// Adds new block with the provided data to the blockchain.
// Route: POST /add
func handleAddBlock(logger *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Println("POST /add")

			data, _ := decode[AddBlockData](r)

			chain := block.GetBlockchain()

			if len(chain) < 1 {
				block.CreateGenesisBlock()
			}
			block.AddBlock(data.Data)

			_ = encode(w, r, 200, "Block added to the chain")
		},
	)
}

// Launches the HTTP server
// Has graceful termination and runs initialization logic before startup.
func Run(ctx context.Context, w io.Writer, args []string) error {
	// Listen for CTRL+C or termination signal
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialise port
	port := strings.Split(os.Getenv("PORT"), ",")
	if len(port) < 1 || port[0] == "" {
		port = []string{"8", "0", "0", "1"}
		log.Println("Port setup from .env failed, using default value")
	}

	// Start new logger
	logger := log.New(w, "", log.LstdFlags)

	// Blockchain initialisation
	// TODO: Add sync between nodes logic
	if err := block.CreateGenesisBlock(); err != nil {
		return fmt.Errorf("Failed to generate genesis block: %w", err)
	}

	// HTTP server setup
	srv := NewServer(logger)
	httpServer := &http.Server{
		Addr:    ":" + strings.Join(port, ""),
		Handler: srv,
	}

	// Server start in goroutine
	go func() {
		logger.Printf("Listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	// Graceful shutdown
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Block until termination signal
		<-ctx.Done()
		logger.Println("Shutting down HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down HTTP server: %s\n", err)
		}
	}()

	wg.Wait()

	return nil
}
