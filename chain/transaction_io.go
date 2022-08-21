package blockchain

import (
	"bytes"

	"gambim.com/blockchain/wallet"
)

type TxOutput struct {
	Value         int
	PublicKeyHash []byte
}
type TxInput struct {
	ID          []byte
	OutputIndex int
	Signature   []byte
	PublicKey   []byte
}

func (input *TxInput) UsesKey(publicKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(input.PublicKey)

	return bytes.Compare(lockingHash, publicKeyHash) == 0
}

func (output *TxOutput) Lock(address []byte) {
	publicKeyFullHash := wallet.Base58Decode(address)
	publicKeyHash := publicKeyFullHash[1 : len(publicKeyFullHash)-wallet.GetChecksumLength()]
	output.PublicKeyHash = publicKeyHash
}

func (output *TxOutput) IsLockedWithKey(publicHashKey []byte) bool {
	return bytes.Compare(output.PublicKeyHash, publicHashKey) == 0
}

func NewTransactionOutput(value int, address string) *TxOutput {
	transactionOutput := &TxOutput{Value: value}
	transactionOutput.Lock([]byte(address))
	return transactionOutput
}
