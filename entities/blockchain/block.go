package blockchain

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
