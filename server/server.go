package server

import (
	"GoChain/block"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// Holds known nodes port
var nodes = make(map[int]struct{})

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

// Checks if incoming request port is recognised or not
// If not recognised adds to known nodes
func checkIfNodeRecognised(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only need port address when working on locally
			_, port, _ := net.SplitHostPort(r.RemoteAddr)

			logger.Printf("Request port: %s", port)
			portInt, _ := strconv.Atoi(port)
			addNode(portInt)
			next.ServeHTTP(w, r)
		})
	}
}

// If not present, adds new node port to the known nodes
func addNode(port int) {
	_, ok := nodes[port]

	if !ok {
		nodes[port] = struct{}{}
	}
}

// Returns entire blockchain as a JSON array.
// Route: GET /chain
func handleGetChain(logger *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Println("GET /chain")
			_ = encode(w, r, http.StatusOK, block.GetBlockchain())
		},
	)
}

// Returns all nodes as JSON array.
// Route: GET /nodes
func handleGetNodes(logger *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Println("GET /nodes")
			_ = encode(w, r, http.StatusOK, slices.Collect(maps.Keys(nodes)))
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

			_ = encode(w, r, http.StatusOK, "Block added to the chain")
		},
	)
}

// Defines the JSON body for POST /receive-block request
type ReceiveBlockData struct {
	Data block.Block `json:"data"`
}

// Adds new block with the provided data to the blockchain.
// Route: POST /receive-block
func handleBlockReceive(logger *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Println("POST /receive-block")

			data, err := decode[ReceiveBlockData](r)

			if err != nil {
				logger.Printf("Failed to decode body: %v", err)
				http.Error(w, "Invalid request body", http.StatusBadRequest)
			}

			err = block.AddMinedBlock(data.Data)

			if err != nil {
				_ = encode(w, r, http.StatusBadRequest, err)
			} else {
				_ = encode(w, r, http.StatusOK, "Block added to the chain")

			}
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
	port := strings.Join(strings.Split(os.Getenv("PORT"), ","), "")
	if len(port) < 1 || port == "" {
		port = "8001"
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
		Addr:    ":" + port,
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
