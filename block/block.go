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

func GetBlockchain() []Block {
	return blockchain
}

// Sets new chain if old one is smaller
func SetBlockchain(chain []Block) {
	if len(chain) > len(blockchain) {
		blockchain = chain
	}
}

// Checks if block is in correct format
// Doesn't check PrevHash and Hash
func IsBlockValid(block Block) (bool, error) {
	// Check if input block index is valid
	if block.Index < 0 {
		return false, fmt.Errorf("Input block index cannot be negative")
	}

	// Check if input block time is in correct format
	_, err := time.Parse("2006-01-02 15:04:05", block.Time)
	if err != nil {
		return false, fmt.Errorf("Input block time is in wrong format or is empty")
	}

	// Check if input block data is valid
	if block.Data == "" {
		return false, fmt.Errorf("Block data cannot be empty")
	}

	// Check if input block nonce is valid
	if block.Nonce < 0 {
		return false, fmt.Errorf("Block nonce cannot be negative")
	}

	return true, nil

}

// Checks if block is valid and Hash is correct
// Doesn't check previous hash
func IsBlockCorrect(block Block) (bool, error) {

	_, err := IsBlockValid(block)

	if err != nil {
		return false, fmt.Errorf("Check is block correct IsBlockValid failed: %v", err)
	}

	calculatedHash, err := CalculateBlockHash(block)

	if err != nil {
		return false, fmt.Errorf("Check is block correct CalculateBlockHash failed: %v", err)
	}

	if calculatedHash != block.Hash {
		return false, fmt.Errorf("Check is block correct calculated hash and block hash doesn't match.")
	}

	return true, nil
}

// Hashes a block using SHA-256
func CalculateBlockHash(block Block) (string, error) {

	// Check if block is valid
	_, err := IsBlockValid(block)

	if err != nil {
		return "", fmt.Errorf("Can't calculate block hash: %v", err)
	}

	dataToHash := strconv.Itoa(block.Index) + block.Time + block.Data + block.PrevHash + strconv.Itoa(block.Nonce)
	hash := sha256.Sum256([]byte(dataToHash))

	return hex.EncodeToString(hash[:]), nil
}

// Create first block for the chain
func CreateGenesisBlock() error {
	genesisBlock := Block{
		Index:    0,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
		Data:     "First block in the chain",
		PrevHash: "",
		Hash:     "",
		Nonce:    0,
	}
	var err error
	genesisBlock.Hash, genesisBlock.Nonce, err = MineBlock(genesisBlock)

	if err != nil {
		return fmt.Errorf("Genesis block greation failed: %v", err)
	}
	blockchain = append(blockchain, genesisBlock)
	return nil
}

// Greates new block and adds to the chain
func AddBlock(data string) {
	lastBlock := blockchain[len(blockchain)-1]
	newBlock := Block{
		Index:    lastBlock.Index + 1,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
		Data:     data,
		PrevHash: lastBlock.Hash,
		Hash:     "",
		Nonce:    0,
	}
	// TODO: Add error handling
	newBlock.Hash, newBlock.Nonce, _ = MineBlock(newBlock)
	blockchain = append(blockchain, newBlock)

}

// Adds mined block to the chain
func AddMinedBlock(data Block) error {
	lastBlock := blockchain[len(blockchain)-1]

	_, err := IsBlockCorrect(data)

	if err == nil && data.PrevHash == lastBlock.Hash {
		blockchain = append(blockchain, data)

	}
	return err
}

// Number of leading zeroes required for the hash
const difficulty = 4

// Returns mined block hash and nonce
func MineBlock(b Block) (string, int, error) {
	hashStart := strings.Repeat("0", difficulty)

	for {
		hash, err := CalculateBlockHash(b)

		if err != nil {
			return "", 0, fmt.Errorf("CalculateBlockHash() error: %v, in MineBlock()", err)
		}

		if strings.HasPrefix(hash, hashStart) {
			return hash, b.Nonce, nil
		}
		b.Nonce++
	}
}
