package blockchain

import (
	"crypto/sha256"
	"reflect"
	"testing"
)

func TestCreateBlock(t *testing.T) {
	prevHash := sha256.Sum256([]byte("test prev hash"))
	blockData := "test block"

	block := &Block{
		Data:     []byte(blockData),
		PrevHash: prevHash[:],
	}

	blockCreatedByFunc := CreateBlock(blockData, prevHash[:])

	if !reflect.DeepEqual(blockCreatedByFunc, block) {
		t.Errorf("Error on block creation")
	}
}

func TestGenesis(t *testing.T) {
	block := CreateBlock("Genesis", []byte{})

	blockCreatedByFunc := Genesis()

	if !reflect.DeepEqual(blockCreatedByFunc, block) {
		t.Errorf("Error on block genesis")
	}
}
