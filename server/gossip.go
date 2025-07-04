package server

import (
	"GoChain/block"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func decodeResponse[T any](body io.ReadCloser) (T, error) {
	defer body.Close()
	var v T
	if err := json.NewDecoder(body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func getNodes(bootstrapNode string) error {
	resp, err := http.Get(bootstrapNode + "/nodes")
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
	resp, err := http.Get(bootstrapNode + "/chain")
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
