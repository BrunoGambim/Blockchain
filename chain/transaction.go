package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func (transaction *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(transaction)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	transaction.ID = hash[:]
}

func CoinBaseTx(to string, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxInput{ID: []byte{}, OutputIndex: -1, Signature: data}
	txout := TxOutput{Value: 100, PublicKey: to}

	transaction := &Transaction{ID: nil, Inputs: []TxInput{txin}, Outputs: []TxOutput{txout}}
	transaction.SetID()

	return transaction
}

func NewTransaction(from string, to string, amount int, chain *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	accumulated, validOutputs := chain.FindSpendableOutputs(from, amount)

	if accumulated < amount {
		log.Panic("Error: not enough funds")
	}

	for txIndex, outs := range validOutputs {
		txID, err := hex.DecodeString(txIndex)
		Handle(err)

		for _, out := range outs {
			input := TxInput{ID: txID, OutputIndex: out, Signature: from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{Value: amount, PublicKey: to})

	if accumulated > amount {
		outputs = append(outputs, TxOutput{Value: accumulated - amount, PublicKey: from})
	}

	transaction := &Transaction{ID: nil, Inputs: inputs, Outputs: outputs}
	transaction.SetID()

	return transaction
}

func (transaction Transaction) IsCoinBase() bool {
	return len(transaction.Inputs) == 1 && len(transaction.Inputs[0].ID) == 0 && transaction.Inputs[0].OutputIndex == -1
}
