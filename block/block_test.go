package block

import (
	"strings"
	"testing"
)

// Calls block.CalculateBlockHash with a block, checking
// for a valid return value
func TestCalculateBlockHash(t *testing.T) {
	block := Block{
		Index:    1,
		Time:     "2025-01-01 12:00:00",
		Data:     "Testing block",
		PrevHash: "",
		Hash:     "",
		Nonce:    0,
	}
	want := "b9955df591916ea5b250e59a6ef58303a39d986b9300df82de82e16bdba48a7b"

	calculatedHash, err := CalculateBlockHash(block)

	if calculatedHash != want || err != nil {
		t.Errorf(`CalculateBlockHash() = %q, %v, want match for %#q`, calculatedHash, err, want)
	}
}

// Calls block.CalculateBlockHash with an empty block, checking if there is error message
func TestCalculateBlockHashEmpty(t *testing.T) {
	block := Block{}

	_, err := CalculateBlockHash(block)

	if err == nil {
		t.Errorf(`CalculateBlockHash() didn't return error %v`, err)
	}
}

// Calls block.CalculateBlockHash with a wrong index, checking if there is error message
func TestCalculateBlockHashWrongIndex(t *testing.T) {
	block := Block{
		Index:    -10,
		Time:     "2025-01-01 12:00:00",
		Data:     "Testing block",
		PrevHash: "",
		Hash:     "",
		Nonce:    0,
	}

	_, err := CalculateBlockHash(block)

	if err == nil {
		t.Errorf(`CalculateBlockHash() didn't return wrong index error %v`, err)
	}
}

// Calls block.CalculateBlockHash with a wrong time format, checking if there is error message
func TestCalculateBlockHashWrongTime(t *testing.T) {
	block := Block{
		Index:    12,
		Time:     "202501-01 12:00:00",
		Data:     "Testing block",
		PrevHash: "",
		Hash:     "",
		Nonce:    0,
	}

	_, err := CalculateBlockHash(block)

	if err == nil {
		t.Errorf(`CalculateBlockHash() didn't return wrong time error %v`, err)
	}
}

// Calls block.CalculateBlockHash with a wrong nonce, checking if there is error message
func TestCalculateBlockHashWrongNonce(t *testing.T) {
	block := Block{
		Index:    12,
		Time:     "2025-01-01 12:00:00",
		Data:     "Testing block",
		PrevHash: "",
		Hash:     "",
		Nonce:    -10,
	}

	_, err := CalculateBlockHash(block)

	if err == nil {
		t.Errorf(`CalculateBlockHash() didn't return wrong nonce error %v`, err)
	}
}

// Calls block.MineBlock with a block, checking
// for a valid return value
func TestMineBlock(t *testing.T) {
	block := Block{
		Index:    1,
		Time:     "2025-01-01 12:00:00",
		Data:     "Testing block",
		PrevHash: "",
		Hash:     "",
		Nonce:    0,
	}
	hashStart := "0000"

	calculatedHash, nonce, err := MineBlock(block)

	if !strings.HasPrefix(calculatedHash, hashStart) || nonce < 0 || err != nil {
		t.Errorf(`Block maining failed: %q, %q, %v`, calculatedHash, nonce, err)
	}
}

// Calls block.MineBlock with an empty block, checking if there is error message
func TestMineBlockEmpty(t *testing.T) {
	block := Block{}

	_, _, err := MineBlock(block)

	if err == nil {
		t.Errorf(`MineBlock() didn't return error %v`, err)
	}
}

// Calls block.MineBlock with a wrong index, checking if there is error message
func TestMineBlockWrongIndex(t *testing.T) {
	block := Block{
		Index:    -10,
		Time:     "2025-01-01 12:00:00",
		Data:     "Testing block",
		PrevHash: "",
		Hash:     "",
		Nonce:    0,
	}

	_, _, err := MineBlock(block)

	if err == nil {
		t.Errorf(`MineBlock() didn't return wrong index error %v`, err)
	}
}

// Calls block.MineBlock with a wrong time format, checking if there is error message
func TestMineBlockWrongTime(t *testing.T) {
	block := Block{
		Index:    12,
		Time:     "202501-01 12:00:00",
		Data:     "Testing block",
		PrevHash: "",
		Hash:     "",
		Nonce:    0,
	}

	_, _, err := MineBlock(block)

	if err == nil {
		t.Errorf(`MineBlock() didn't return wrong time error %v`, err)
	}
}

// Calls block.MineBlock with a wrong nonce, checking if there is error message
func TestMineBlockWrongNonce(t *testing.T) {
	block := Block{
		Index:    12,
		Time:     "2025-01-01 12:00:00",
		Data:     "Testing block",
		PrevHash: "",
		Hash:     "",
		Nonce:    -10,
	}

	_, _, err := MineBlock(block)

	if err == nil {
		t.Errorf(`MineBlock() didn't return wrong nonce error %v`, err)
	}
}
