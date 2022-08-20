package blockchain

import (
	"reflect"
	"testing"
)

func TestInitBlockchain(t *testing.T) {
	blockchainInitiatedByFunc := InitBlockchain()

	chain := &Blockchain{
		blocks: []*Block{Genesis()},
	}

	if !reflect.DeepEqual(blockchainInitiatedByFunc, chain) || len(blockchainInitiatedByFunc.blocks) != 1 {
		t.Errorf("Error on blockchain initiation")
	}
}

func TestAddBlock(t *testing.T) {
	chain := InitBlockchain()
	chain.AddBlock("First Block")

	if len(chain.blocks) != 2 {
		t.Errorf("Error on adding block to the blockchain")
	}

	chain.AddBlock("Second Block")
	chain.AddBlock("Third Block")

	if len(chain.blocks) != 4 {
		t.Errorf("Error on adding block to the blockchain")
	}

	if !reflect.DeepEqual(chain.blocks[1].Data, []byte("First Block")) || !reflect.DeepEqual(chain.blocks[3].Data, []byte("Third Block")) {
		t.Errorf("Error on block added data")
	}

	if !reflect.DeepEqual(chain.blocks[0].Hash, chain.blocks[1].PrevHash) || !reflect.DeepEqual(chain.blocks[2].Hash, chain.blocks[3].PrevHash) {
		t.Errorf("Error on block prev hash")
	}
}
