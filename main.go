package main

import (
	"fmt"

	"gambim.com/blockchain/blockchain"
)

func main() {
	chain := blockchain.InitBlockchain()
	chain.AddBlock("First Block")
	chain.AddBlock("Second Block")
	chain.AddBlock("Third Block")
	/*for _, block := range chain.Blocks {
		fmt.Printf("data:%s  hash:%x  prev_hash:%x\n", block.Data, block.Hash, block.PrevHash)
	}*/

	iterator := chain.Iterator()
	for len(iterator.IteratorHash) != 0 {
		block := iterator.Next()
		fmt.Printf("data:%s  hash:%x  prev_hash:%x\n", block.Data, block.Hash, block.PrevHash)
	}

	//print("%x", iterator.IteratorHash)
	chain.Database.Close()
}
