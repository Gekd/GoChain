package block

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Struct for block
type Block struct {
	Index    int
	Time     string
	Data     string
	PrevHash string
	Hash     string
	Nonce    int
}

var blockchain []Block

// Hashes a block using SHA-256
func CalculateBlockHash(block Block) (string, error) {
	// Check if input block index is valid
	if block.Index < 0 {
		return "", fmt.Errorf("Input block index cannot be negative")
	}

	// Check if input block time is in correct format
	_, err := time.Parse("2006-01-02 15:04:05", block.Time)
	if err != nil {
		return "", fmt.Errorf("Input block time is in wrong format or is empty")
	}

	// Check if input block data is valid
	if block.Data == "" {
		return "", fmt.Errorf("Block data cannot be empty")
	}

	// Check if input block nonce is valid
	if block.Nonce < 0 {
		return "", fmt.Errorf("Block nonce cannot be negative")
	}

	dataToHash := strconv.Itoa(block.Index) + block.Time + block.Data + block.PrevHash + strconv.Itoa(block.Nonce)
	hash := sha256.Sum256([]byte(dataToHash))

	return hex.EncodeToString(hash[:]), nil
}

// Create first block for the chain
func CreateGenesisBlock() Block {
	return Block{
		Index:    0,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
		Data:     "First block in the chain",
		PrevHash: "",
		Hash:     "",
		Nonce:    0,
	}
}

// Adds mined block to the chain
func AddBlock(data string) {
	lastBlock := blockchain[len(blockchain)-1]
	newBlock := Block{
		Index:    lastBlock.Index + 1,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
		Data:     data,
		PrevHash: lastBlock.PrevHash,
		Hash:     "",
		Nonce:    0,
	}

	newBlock.Hash, newBlock.Nonce = MineBlock(newBlock)
	blockchain = append(blockchain, newBlock)

}

// Number of leading zeroes required for the hash
const difficulty = 4

// Returns mined block hash and nonce
func MineBlock(b Block) (string, int) {
	hashStart := strings.Repeat("0", difficulty)

	for {
		// TODO: Add error handling
		hash, _ := CalculateBlockHash(b)
		if strings.HasPrefix(hash, hashStart) {
			return hash, b.Nonce
		}
		b.Nonce++
	}
}

// Checks if block is valid
func IsValidBlock(newBlock, prevBlock Block) bool {
	// If previous block isn't before new block return false
	if newBlock.Index != prevBlock.Index {
		return false
	}

	// If previous block hash isn't stored in new block return false
	if newBlock.PrevHash != prevBlock.Hash {
		return false
	}

	hashStart := strings.Repeat("0", difficulty)

	// TODO: Add error handling
	calculatedHash, _ := CalculateBlockHash(newBlock)

	// If new block isn't correctly mined return false
	if !strings.HasPrefix(calculatedHash, hashStart) {
		return false
	}

	return true
}
