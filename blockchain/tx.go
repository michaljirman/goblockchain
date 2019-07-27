package blockchain

import (
	"bytes"

	"github.com/michaljirman/goblockchain/wallet"
)

type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	PubKey    []byte
}

// Creates a new tx output from value and address
func NewTxOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}

// Text if tx input uses pubKeyHash
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

// Locks the tx output
func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := wallet.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-wallet.ChecksumLength]
	out.PubKeyHash = pubKeyHash
}

// Checks if tx output is locked with a key
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}
