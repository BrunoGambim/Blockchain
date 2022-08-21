package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"gambim.com/blockchain/wallet"
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

	txin := TxInput{ID: []byte{}, OutputIndex: -1, PublicKey: []byte(data)}
	txout := NewTransactionOutput(100, to)

	transaction := &Transaction{ID: nil, Inputs: []TxInput{txin}, Outputs: []TxOutput{*txout}}
	transaction.SetID()

	return transaction
}

func NewTransaction(from string, to string, amount int, chain *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	wallets, err := wallet.CreateWallets()
	Handle(err)
	wlt := wallets.GetWallet(from)
	publicKeyHash := wallet.PublicKeyHash(wlt.PublicKey)

	accumulated, validOutputs := chain.FindSpendableOutputs(publicKeyHash, amount)

	if accumulated < amount {
		log.Panic("Error: not enough funds")
	}

	for txIndex, outs := range validOutputs {
		txID, err := hex.DecodeString(txIndex)
		Handle(err)

		for _, out := range outs {
			input := TxInput{ID: txID, OutputIndex: out, PublicKey: wlt.PublicKey}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, *NewTransactionOutput(amount, to))

	if accumulated > amount {
		outputs = append(outputs, *NewTransactionOutput(accumulated-amount, from))
	}

	transaction := &Transaction{ID: nil, Inputs: inputs, Outputs: outputs}
	transaction.ID = transaction.Hash()
	chain.SignTransaction(transaction, wlt.PrivateKey)
	//transaction.SetID()

	return transaction
}

func (transaction Transaction) IsCoinBase() bool {
	return len(transaction.Inputs) == 1 && len(transaction.Inputs[0].ID) == 0 && transaction.Inputs[0].OutputIndex == -1
}

func (transaction *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(transaction)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func (transaction *Transaction) Hash() []byte {
	var hash [32]byte

	transactionCopy := *transaction
	transactionCopy.ID = []byte{}

	hash = sha256.Sum256(transaction.Serialize())

	return hash[:]
}

func (transaction *Transaction) Sign(privKey ecdsa.PrivateKey, prevTransactions map[string]Transaction) {
	if transaction.IsCoinBase() {
		return
	}

	for _, input := range transaction.Inputs {
		if prevTransactions[hex.EncodeToString(input.ID)].ID == nil {
			log.Panic("Error: Previous transaction is not exist")
		}
	}

	transactionCopy := transaction.TrimmedCopy()

	for inputIndex, input := range transactionCopy.Inputs {
		prevTransaction := prevTransactions[hex.EncodeToString(input.ID)]
		transactionCopy.Inputs[inputIndex].Signature = nil
		transactionCopy.Inputs[inputIndex].PublicKey = prevTransaction.Outputs[input.OutputIndex].PublicKeyHash
		transactionCopy.ID = transactionCopy.Hash()
		transactionCopy.Inputs[inputIndex].PublicKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, transactionCopy.ID)
		Handle(err)

		signature := append(r.Bytes(), s.Bytes()...)

		transaction.Inputs[inputIndex].Signature = signature
	}
}

func (transaction *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, input := range transaction.Inputs {
		inputs = append(inputs, TxInput{ID: input.ID, OutputIndex: input.OutputIndex})
	}

	for _, output := range transaction.Outputs {
		outputs = append(outputs, output)
	}

	transactionCopy := Transaction{ID: transaction.ID, Inputs: inputs, Outputs: outputs}

	return transactionCopy
}

func (transaction Transaction) Verify(prevTransactions map[string]Transaction) bool {
	if transaction.IsCoinBase() {
		return true
	}

	for _, input := range transaction.Inputs {
		if prevTransactions[hex.EncodeToString(input.ID)].ID == nil {
			log.Panic("Error: Previous transaction is not exist")
		}
	}

	transactionCopy := transaction.TrimmedCopy()
	curve := elliptic.P256()

	for inputIndex, input := range transaction.Inputs {
		prevTransaction := prevTransactions[hex.EncodeToString(input.ID)]
		transactionCopy.Inputs[inputIndex].Signature = nil
		transactionCopy.Inputs[inputIndex].PublicKey = prevTransaction.Outputs[input.OutputIndex].PublicKeyHash
		transactionCopy.ID = transactionCopy.Hash()
		transactionCopy.Inputs[inputIndex].PublicKey = nil

		r := big.Int{}
		s := big.Int{}
		signatureLength := len(input.Signature)
		r.SetBytes(input.Signature[:(signatureLength / 2)])
		s.SetBytes(input.Signature[(signatureLength / 2):])

		x := big.Int{}
		y := big.Int{}
		publicKeyLength := len(input.PublicKey)
		x.SetBytes(input.PublicKey[:(publicKeyLength / 2)])
		y.SetBytes(input.PublicKey[(publicKeyLength / 2):])

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, transactionCopy.ID, &r, &s) == false {
			return false
		}
	}
	return true
}

func (transaction *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("---Transaction %x:", transaction.ID))
	for i, input := range transaction.Inputs {
		lines = append(lines, fmt.Sprintf("     Input index %d:", i))
		lines = append(lines, fmt.Sprintf("        Input ID: %x", input.ID))
		lines = append(lines, fmt.Sprintf("        Output index: %d", input.OutputIndex))
		lines = append(lines, fmt.Sprintf("        Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("        Public key: %x", input.PublicKey))
	}

	for i, output := range transaction.Outputs {
		lines = append(lines, fmt.Sprintf("     Output index %d:", i))
		lines = append(lines, fmt.Sprintf("        Script: %x", output.PublicKeyHash))
		lines = append(lines, fmt.Sprintf("        Value: %d", output.Value))
	}

	return strings.Join(lines, "\n")
}
