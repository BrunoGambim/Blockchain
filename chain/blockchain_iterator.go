package blockchain

import (
	"go.etcd.io/bbolt"
)

type BlockchainIterator struct {
	IteratorHash []byte
	Database     *bbolt.DB
}

func (iterator *BlockchainIterator) Next() *Block {
	var block *Block

	err := iterator.Database.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("blockchain bucket"))
		if bucket == nil {
			Handle(bbolt.ErrBucketNotFound)
		}

		blockBlob := bucket.Get(iterator.IteratorHash)
		block = Deserialize(blockBlob)
		return nil
	})
	Handle(err)

	iterator.IteratorHash = block.PrevHash

	return block
}
