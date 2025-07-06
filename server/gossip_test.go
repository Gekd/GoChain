package server

import (
	"io"
	"strings"
	"testing"
)

// Defines the JSON body for testing decode function
type testDecodeData struct {
	Data string `json:"data"`
}

// Calls gossip.decodeResponse with a correct JSON data format, checking
// for a valid return value
func TestDecodeResponse(t *testing.T) {

	testMessage := "Test message"
	JSONToDecode := `{ "data": "` + testMessage + `" }`

	ioReadCloserJSON := io.NopCloser(strings.NewReader(JSONToDecode))

	resp, err := decodeResponse[testDecodeData](ioReadCloserJSON)

	if err != nil {
		t.Errorf("decodeResponse() returned an error: %v", err)
	}

	if resp.Data != testMessage {
		t.Errorf("decodeResponse() returned an error: %v", err)
	}
}

// Calls gossip.decodeResponse with a wrong JSON data format, checking if there is error message
func TestDecodeResponseWrongJSONFormat(t *testing.T) {

	testMessage := "Test message"
	JSONToDecode := `{ "data": "` + testMessage + `" }`

	ioReadCloserJSON := io.NopCloser(strings.NewReader(JSONToDecode))

	_, err := decodeResponse[string](ioReadCloserJSON)

	if err == nil {
		t.Error("decodeResponse() didn't return an error with wrong JSON format.")
	}

}

// Calls gossip.decodeResponse with a empty string, checking if there is error message
func TestDecodeResponseEmptyInput(t *testing.T) {

	testMessage := ""

	ioReadCloserJSON := io.NopCloser(strings.NewReader(testMessage))

	_, err := decodeResponse[string](ioReadCloserJSON)

	if err == nil {
		t.Error("decodeResponse() didn't return an error with empty input.")
	}

}

// Calls gossip.getNodes with a wrong hostname, checking if there is error message
func TestGetNodesCantConnect(t *testing.T) {

	bootstrapNode := "http://127.0.0.1:8888"
	err := getNodes(bootstrapNode)

	if err == nil {
		t.Error("getNodes() didn't return any errors")
	}
}

// Calls gossip.getChain with a wrong hostname, checking if there is error message
func TestGetChainCantConnect(t *testing.T) {

	bootstrapNode := "http://127.0.0.1:8888"
	err := getChain(bootstrapNode)

	if err == nil {
		t.Error("getChain() didn't return any errors")
	}
}

// Calls gossip.syncNode with a wrong hostname, checking if there is error message
func TestSyncNodeCantConnect(t *testing.T) {

	bootstrapNode := "http://127.0.0.1:8888"
	err := getChain(bootstrapNode)

	if err == nil {
		t.Error("syncNode() didn't return any errors")
	}
}
