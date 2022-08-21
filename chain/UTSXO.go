package blockchain

import (
	"encoding/hex"
	"log"

	"go.etcd.io/bbolt"
)

var (
	utxoBucket = []byte("utxo")
)

type UTXOSet struct {
	Chain *Blockchain
}

func NewUTXOSet(chain *Blockchain) *UTXOSet {
	return &UTXOSet{
		Chain: chain,
	}
}

func (u *UTXOSet) Update(block *Block) {
	err := u.Chain.Database.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(utxoBucket)
		if bucket == nil {
			return bbolt.ErrBucketNotFound
		}
		for _, transaction := range block.Transactions {
			if transaction.IsCoinBase() == false {
				for _, input := range transaction.Inputs {
					updatedOuts := TxOutputs{}
					item := bucket.Get(input.ID)
					if item == nil {
						log.Panic("Output doesn't exist")
					}
					outputs := DeserializeOutputs(item)
					for outputIndex, output := range outputs.Outputs {
						if outputIndex != input.OutputIndex {
							updatedOuts.Outputs = append(updatedOuts.Outputs, output)
						}
					}
					err := bucket.Delete(input.ID)
					if err != nil {
						return err
					}
					if len(updatedOuts.Outputs) != 0 {
						if err = bucket.Put(input.ID, updatedOuts.Serialize()); err != nil {
							return err
						}
					}
				}
			}
			newOutputs := TxOutputs{}
			for _, output := range transaction.Outputs {
				newOutputs.Outputs = append(newOutputs.Outputs, output)
			}

			if err := bucket.Put(transaction.ID, newOutputs.Serialize()); err != nil {
				return err
			}
		}

		return nil
	})

	Handle(err)
}

func (u *UTXOSet) CountTransactions() int {
	counter := 0
	err := u.Chain.Database.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(utxoBucket)
		if bucket == nil {
			return bbolt.ErrBucketNotFound
		}

		bucket.ForEach(func(key []byte, value []byte) error {
			counter++
			return nil
		})
		return nil
	})
	Handle(err)

	return counter
}

func (u *UTXOSet) Reindex() {
	u.DeleteAll()
	unspentOutputs := u.Chain.FindUnspentTransactionOutputs()

	err := u.Chain.Database.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucket(utxoBucket)
		if err != nil {
			return err
		}
		for transactionId, outputs := range unspentOutputs {
			key, err := hex.DecodeString(transactionId)
			if err != nil {
				return err
			}
			bucket.Put(key, outputs.Serialize())
		}

		return nil
	})
	Handle(err)
}

func (u *UTXOSet) DeleteAll() {
	err := u.Chain.Database.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(utxoBucket)
		if bucket != nil {
			return tx.DeleteBucket(utxoBucket)
		} else {
			return nil
		}
	})
	Handle(err)
}

func (u *UTXOSet) FindUnspentTransactionOutputs(publicHashKey []byte) []TxOutput {
	var unspentTransactionOutputs []TxOutput

	err := u.Chain.Database.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(utxoBucket)
		if bucket == nil {
			return bbolt.ErrBucketNotFound
		}

		bucket.ForEach(func(key []byte, item []byte) error {
			txOutputs := DeserializeOutputs(item)
			for _, output := range txOutputs.Outputs {
				if output.IsLockedWithKey(publicHashKey) {
					unspentTransactionOutputs = append(unspentTransactionOutputs, output)
				}
			}
			return nil
		})
		return nil
	})
	Handle(err)

	return unspentTransactionOutputs
}

func (u *UTXOSet) FindSpendableOutputs(publicHashKey []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0

	err := u.Chain.Database.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(utxoBucket)
		if bucket == nil {
			return bbolt.ErrBucketNotFound
		}

		bucket.ForEach(func(key []byte, item []byte) error {
			txOutputs := DeserializeOutputs(item)
			for outputIndex, output := range txOutputs.Outputs {
				if output.IsLockedWithKey(publicHashKey) && accumulated < amount {
					accumulated += output.Value
					unspentOutputs[hex.EncodeToString(key)] = append(unspentOutputs[hex.EncodeToString(key)], outputIndex)

					if accumulated >= amount {
						return nil
					}
				}
			}
			return nil
		})
		return nil
	})
	Handle(err)

	return accumulated, unspentOutputs
}
