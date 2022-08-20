package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nounce   int
}

func CreateBlock(data string, prevHash []byte) *Block {
	newBlock := &Block{
		PrevHash: prevHash,
		Data:     []byte(data),
	}
	proofOfWork := NewProof(newBlock)
	nounce, hash := proofOfWork.Run()

	newBlock.Hash = hash[:]
	newBlock.Nounce = nounce

	return newBlock
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func (block *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(block)

	Handle(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
