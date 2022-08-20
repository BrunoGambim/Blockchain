package blockchain

type Blockchain struct {
	blocks []*Block
}

func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, newBlock)
}

func InitBlockchain() *Blockchain {
	return &Blockchain{
		blocks: []*Block{
			Genesis(),
		},
	}
}
