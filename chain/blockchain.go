package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"runtime"

	"go.etcd.io/bbolt"
)

const (
	genesisData = "First Transaction from Genesis"
	dbPath      = "tmp/blocks"
	dbName      = "my.db"
)

type Blockchain struct {
	LastHash []byte
	Database *bbolt.DB
}

func DBExists() bool {
	_, err := os.Stat(fmt.Sprintf("%s/%s", dbPath, dbName))
	return !os.IsNotExist(err)
}

func InitBlockchain(address string) *Blockchain {
	var lastHash []byte = nil

	if DBExists() {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	db, err := bbolt.Open(fmt.Sprintf("%s/%s", dbPath, dbName), 0600, nil)
	Handle(err)

	err = db.Update(func(tx *bbolt.Tx) error {

		bucket, err := tx.CreateBucket([]byte("blockchain bucket"))
		Handle(err)

		transaction := CoinBaseTx(address, genesisData)
		genesis := Genesis(transaction)
		err = bucket.Put(genesis.Hash, genesis.Serialize())
		Handle(err)

		err = bucket.Put([]byte("last hash"), genesis.Hash)
		lastHash = genesis.Hash

		return err
	})

	Handle(err)

	return &Blockchain{
		LastHash: lastHash,
		Database: db,
	}
}

func ContinueBlockchain(address string) *Blockchain {
	var lastHash []byte = nil

	if !DBExists() {
		fmt.Println("Blockchain not exists")
		runtime.Goexit()
	}

	db, err := bbolt.Open(fmt.Sprintf("%s/%s", dbPath, dbName), 0600, nil)
	Handle(err)

	err = db.View(func(tx *bbolt.Tx) error {

		bucket := tx.Bucket([]byte("blockchain bucket"))
		lastHash = bucket.Get([]byte("last hash"))

		return err
	})
	Handle(err)

	return &Blockchain{
		LastHash: lastHash,
		Database: db,
	}
}

func (chain *Blockchain) FindUnspentTransactionOutputs() map[string]TxOutputs {
	unspentTransactionOutputs := make(map[string]TxOutputs)
	spentTransactions := make(map[string][]int)

	iterator := chain.Iterator()
	for len(iterator.IteratorHash) != 0 {
		block := iterator.Next()

		for _, transaction := range block.Transactions {
			transactionId := hex.EncodeToString(transaction.ID)

		Outputs:
			for outputIndex, output := range transaction.Outputs {
				if spentTransactions[transactionId] != nil {
					for _, spentOut := range spentTransactions[transactionId] {
						if spentOut == outputIndex {
							continue Outputs
						}
					}
				}
				outputs := unspentTransactionOutputs[transactionId]
				outputs.Outputs = append(outputs.Outputs, output)
				unspentTransactionOutputs[transactionId] = outputs

			}
			if transaction.IsCoinBase() == false {
				for _, input := range transaction.Inputs {
					inputTransactionId := hex.EncodeToString(input.ID)
					spentTransactions[inputTransactionId] = append(spentTransactions[inputTransactionId], input.OutputIndex)
				}
			}
		}
	}
	return unspentTransactionOutputs
}

func (chain *Blockchain) AddBlock(transactions []*Transaction) *Block {
	var err error
	newBlock := CreateBlock(transactions, chain.LastHash)

	err = chain.Database.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("blockchain bucket"))
		if bucket == nil {
			Handle(bbolt.ErrBucketNotFound)
		}

		if chain.LastHash != nil {
			err = bucket.Put(newBlock.Hash, newBlock.Serialize())
			Handle(err)

			err = bucket.Delete([]byte("last hash"))
			Handle(err)

			err = bucket.Put([]byte("last hash"), newBlock.Hash)
			Handle(err)
			chain.LastHash = newBlock.Hash
		}
		return nil
	})

	return newBlock
}

func (chain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		IteratorHash: chain.LastHash,
		Database:     chain.Database,
	}
}

func (chain *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	iterator := chain.Iterator()
	for len(iterator.IteratorHash) != 0 {
		block := iterator.Next()

		for _, transaction := range block.Transactions {
			if bytes.Compare(transaction.ID, ID) == 0 {
				return *transaction, nil
			}
		}
	}
	return Transaction{}, errors.New("Transaction does not exist")
}

func (chain *Blockchain) SignTransaction(transaction *Transaction, privateKey ecdsa.PrivateKey) {
	prevTransactions := make(map[string]Transaction)
	for _, input := range transaction.Inputs {
		prevTransaction, err := chain.FindTransaction(input.ID)
		Handle(err)
		prevTransactions[hex.EncodeToString(prevTransaction.ID)] = prevTransaction
	}

	transaction.Sign(privateKey, prevTransactions)
}

func (chain *Blockchain) VerifyTransaction(transaction *Transaction, privateKey ecdsa.PrivateKey) {
	prevTransactions := make(map[string]Transaction)
	for _, input := range transaction.Inputs {
		prevTransaction, err := chain.FindTransaction(input.ID)
		Handle(err)
		prevTransactions[hex.EncodeToString(prevTransaction.ID)] = prevTransaction
	}

	transaction.Verify(prevTransactions)
}
