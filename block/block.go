package block

import (
	"crypto/sha256"
	"encoding/hex"
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
func CalculateBlockHash(block Block) string {
	dataToHash := strconv.Itoa(block.Index) + block.Time + block.Data + block.PrevHash + strconv.Itoa(block.Nonce)
	hash := sha256.Sum256([]byte(dataToHash))
	return hex.EncodeToString(hash[:])
}

// Create first block for the chain
func CreateGenesisBlock() Block {
	return Block{
		Index:    0,
		Time:     time.Now().String(),
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
		Time:     time.Now().String(),
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
		hash := CalculateBlockHash(b)
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
	// If new block isn't correctly mined return false
	if !strings.HasPrefix(CalculateBlockHash(newBlock), hashStart) {
		return false
	}

	return true
}
