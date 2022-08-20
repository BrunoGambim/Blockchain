package blockchain

import (
	"bytes"
	"crypto/sha256"
	"reflect"
	"testing"
)

func TestDeriveHash(t *testing.T) {
	prevHash := sha256.Sum256([]byte("teste"))
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
