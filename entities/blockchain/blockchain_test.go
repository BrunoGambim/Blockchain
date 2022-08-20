package blockchain

import (
	"reflect"
	"testing"
)

func TestInitBlockchain(t *testing.T) {
	blockchainInitiatedByFunc := InitBlockchain()

	chain := &Blockchain{
		Blocks: []*Block{Genesis()},
	}

	if !reflect.DeepEqual(blockchainInitiatedByFunc, chain) || len(blockchainInitiatedByFunc.Blocks) != 1 {
		t.Errorf("Error on blockchain initiation")
	}
}

func TestAddBlock(t *testing.T) {
	chain := InitBlockchain()
	chain.AddBlock("First Block")

	if len(chain.Blocks) != 2 {
		t.Errorf("Error on adding block to the blockchain")
	}

	chain.AddBlock("Second Block")
	chain.AddBlock("Third Block")

	if len(chain.Blocks) != 4 {
		t.Errorf("Error on adding block to the blockchain")
	}

	if !reflect.DeepEqual(chain.Blocks[1].Data, []byte("First Block")) || !reflect.DeepEqual(chain.Blocks[3].Data, []byte("Third Block")) {
		t.Errorf("Error on block added data")
	}

	if !reflect.DeepEqual(chain.Blocks[0].Hash, chain.Blocks[1].PrevHash) || !reflect.DeepEqual(chain.Blocks[2].Hash, chain.Blocks[3].PrevHash) {
		t.Errorf("Error on block prev hash")
	}
}
