package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math"
	"math/big"
)

const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	proofOfWork := &ProofOfWork{Block: block, Target: target}

	return proofOfWork
}

func (proofOfWork *ProofOfWork) InitData(nounce int) []byte {
	data := bytes.Join(
		[][]byte{
			proofOfWork.Block.PrevHash,
			proofOfWork.Block.Data,
			ToHex(int64(nounce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

func (proofOfWork *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nounce := 0

	for nounce < math.MaxInt64 {
		data := proofOfWork.InitData(nounce)
		hash = sha256.Sum256(data)

		//fmt.Printf("\r%x\n", hash)

		intHash.SetBytes(hash[:])

		if intHash.Cmp(proofOfWork.Target) == -1 {
			break
		} else {
			nounce++
		}
	}

	return nounce, hash[:]
}

func (proofOfWork *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := proofOfWork.InitData(proofOfWork.Block.Nounce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(proofOfWork.Target) == -1
}

func ToHex(num int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}
