package blockchain

import (
	"bytes"
	"crypto/sha256"
	"reflect"
	"testing"
)

func TestDeriveHash(t *testing.T) {
	prevHash := sha256.Sum256([]byte("test prev hash"))
	blockData := "test block"

	hash := sha256.Sum256(bytes.Join([][]byte{[]byte(blockData), prevHash[:]}, []byte{}))

	block := Block{
		Data:     []byte(blockData),
		PrevHash: prevHash[:],
	}
	block.DeriveHash()

	if !reflect.DeepEqual(hash[:], block.Hash) {
		t.Errorf("Error on hash calc")
	}
}

func TestCreateBlock(t *testing.T) {
	prevHash := sha256.Sum256([]byte("test prev hash"))
	blockData := "test block"

	block := &Block{
		Data:     []byte(blockData),
		PrevHash: prevHash[:],
	}
	block.DeriveHash()

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
