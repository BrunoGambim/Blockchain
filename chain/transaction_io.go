package blockchain

type TxOutput struct {
	Value     int
	PublicKey string
}
type TxInput struct {
	ID          []byte
	OutputIndex int
	Signature   string
}

func (input *TxInput) CanUnlock(data string) bool {
	return input.Signature == data
}

func (output *TxOutput) CanBeUnlocked(data string) bool {
	return output.PublicKey == data
}
