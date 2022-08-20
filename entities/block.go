package entities

import (
	"bytes"
	"crypto/sha256"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

func (block *Block) DeriveHash() {
	info := bytes.Join([][]byte{block.Data, block.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	block.Hash = hash[:]
}

func (block *Block) CreateBlock(data string, prevHash []byte) *Block {
	newBlock := &Block{
		PrevHash: prevHash,
		Data:     []byte(data),
	}
	newBlock.DeriveHash()
	return block
}
