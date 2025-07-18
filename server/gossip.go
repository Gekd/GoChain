package server

import (
	"GoChain/block"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"math"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"strings"
)

func encodeRequest[T any](v T) (io.Reader, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("encode json: %w", err)
	}
	return bytes.NewReader(data), nil
}

func decodeResponse[T any](body io.ReadCloser) (T, error) {
	defer body.Close()
	var v T
	if err := json.NewDecoder(body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

// Checks nodes to make sure they are alive
func checkNodes(logger *log.Logger) {

	// Initialise local address
	localAddr := strings.Join(strings.Split(os.Getenv("LOCAL_ADDR"), ","), "")
	if len(localAddr) < 1 || localAddr == "" {
		localAddr = "0.0.0.0:8001"
		log.Println("Local address setup from .env failed, using default value")
	}

	knownNodes := slices.Collect(maps.Keys(nodes))
	checkLimit := int(math.RoundToEven(math.Sqrt(float64(len(knownNodes)))))

	nodesToCheck := []string{}

	for {
		// Loop as long is needed
		if len(nodesToCheck) == checkLimit {
			break
		}

		randInt := rand.Intn(len(knownNodes))

		if !slices.Contains(nodesToCheck, knownNodes[randInt]) {
			nodesToCheck = append(nodesToCheck, knownNodes[randInt])
		}
	}
	logger.Printf("Nodes to check: %v", nodesToCheck)

	for _, s := range nodesToCheck {
		func() {
			url := "http://" + s + "/ping"

			req, err := http.NewRequest("GET", url, nil)

			if err != nil {
				logger.Printf("Failed to create request for node %v: %v\n", s, err)
				return
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Node-Addr", localAddr)

			client := &http.Client{}

			resp, err := client.Do(req)

			if err != nil {
				logger.Printf("Error connecting to host: %v, %v", s, err)
				removeNode(s)
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				logger.Printf("Unexpected response: %v, from %v", resp.StatusCode, s)
			}

			/*data, decodeErr := decodeResponse[GetPingData](resp.Body)

			if decodeErr != nil {
				logger.Printf("Error decoding %v GET /ping: %v", s, decodeErr)
			}*/

		}()
	}

}

func getNodes(bootstrapNode string) error {

	// Initialise local address
	localAddr := strings.Join(strings.Split(os.Getenv("LOCAL_ADDR"), ","), "")
	if len(localAddr) < 1 || localAddr == "" {
		localAddr = "0.0.0.0:8001"
		log.Println("Local address setup from .env failed, using default value")
	}

	url := "http://" + bootstrapNode + "/nodes"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return fmt.Errorf("Failed to create request for node %v: %v\n", bootstrapNode, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Node-Addr", localAddr)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("Error connecting to host: %v, %v", bootstrapNode, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected response: %v, from %v", resp.StatusCode, bootstrapNode)
	}

	addNode(bootstrapNode)

	data, decodeErr := decodeResponse[GetNodesData](resp.Body)

	if decodeErr != nil {
		return fmt.Errorf("Error decoding GET /nodes: %v,", decodeErr)
	}

	if len(data.Data) > 0 {
		for _, node := range data.Data {
			addNode(node)
		}
	}

	return nil
}

func getChain(bootstrapNode string) error {

	// Initialise local address
	localAddr := strings.Join(strings.Split(os.Getenv("LOCAL_ADDR"), ","), "")
	if len(localAddr) < 1 || localAddr == "" {
		localAddr = "0.0.0.0:8001"
		log.Println("Local address setup from .env failed, using default value")
	}

	url := "http://" + bootstrapNode + "/chain"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return fmt.Errorf("Failed to create request for node %v: %v\n", bootstrapNode, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Node-Addr", localAddr)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("Error connecting to host: %v, %v", bootstrapNode, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected response: %v, from %v", resp.StatusCode, bootstrapNode)
	}

	data, decodeErr := decodeResponse[GetChainData](resp.Body)

	if decodeErr != nil {
		return fmt.Errorf("Error decoding GET /chain: %v,", decodeErr)
	}

	if len(data.Data) > 0 {
		block.SetBlockchain(data.Data)
	}

	return nil
}

// Synchronises current node chain and nodes list with bootstrap node
func syncNode(bootstrapNode string) error {

	nodeErr := getNodes(bootstrapNode)
	if nodeErr != nil {
		return nodeErr
	}

	chainErr := getChain(bootstrapNode)
	if chainErr != nil {
		return chainErr
	}

	return nil

}

// Distributes mined block amongst known peers
func shareMinedBlock(logger *log.Logger, block block.Block) {

	// Initialise local address
	localAddr := strings.Join(strings.Split(os.Getenv("LOCAL_ADDR"), ","), "")
	if len(localAddr) < 1 || localAddr == "" {
		localAddr = "0.0.0.0:8001"
		log.Println("Local address setup from .env failed, using default value")
	}

	for node := range nodes {
		body, encodeErr := encodeRequest(ReceiveBlockData{Data: block})

		if encodeErr != nil {
			logger.Printf("Failed to encode payload: %v\n", encodeErr)
			continue
		}

		url := "http://" + node + "/receive-block"

		req, err := http.NewRequest("POST", url, body)

		if err != nil {
			logger.Printf("Failed to create request for node %v: %v\n", node, err)
			continue

		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Node-Addr", localAddr)

		client := &http.Client{}

		resp, err := client.Do(req)

		if err != nil {
			logger.Printf("Block sharing on node %v failed: %v\n", node, err)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("HTTP request code on node %v: %v\n", node, resp.StatusCode)
		}
	}

}
