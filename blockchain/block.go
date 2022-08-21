package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nounce       int
}

func CreateBlock(transactions []*Transaction, prevHash []byte) *Block {
	newBlock := &Block{
		PrevHash:     prevHash,
		Transactions: transactions,
	}
	proofOfWork := NewProof(newBlock)
	nounce, hash := proofOfWork.Run()

	newBlock.Hash = hash[:]
	newBlock.Nounce = nounce

	return newBlock
}

func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, transaction := range block.Transactions {
		txHashes = append(txHashes, transaction.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
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
