package blockchain

import (
	"go.etcd.io/bbolt"
)

type Blockchain struct {
	LastHash []byte
	Database *bbolt.DB
}

func InitBlockchain() *Blockchain {
	var lastHash []byte = nil
	db, err := bbolt.Open("tmp/blocks/my.db", 0600, nil)

	Handle(err)
	/*return &Blockchain{
		Blocks: []*Block{
			Genesis(),
		},
	}*/

	err = db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("blockchain bucket"))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte("blockchain bucket"))
			Handle(err)
		}
		lastHash = bucket.Get([]byte("last hash"))
		if lastHash == nil {
			genesis := Genesis()
			err = bucket.Put(genesis.Hash, genesis.Serialize())
			Handle(err)

			err = bucket.Put([]byte("last hash"), genesis.Hash)
			Handle(err)

			lastHash = genesis.Hash
		}
		return nil
	})

	Handle(err)

	return &Blockchain{
		LastHash: lastHash,
		Database: db,
	}
}

func (chain *Blockchain) AddBlock(data string) {
	var err error
	newBlock := CreateBlock(data, chain.LastHash)

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
}

func (chain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		IteratorHash: chain.LastHash,
		Database:     chain.Database,
	}
}
