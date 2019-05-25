package main

import (
	"bytes"
	"fmt"
)

//新的交易输入
func NewTxIn(PreTxH, Signature, PubKey []byte, PreTxOutIndex int) TxIn {
	return TxIn{PreTxH: PreTxH, PreTxOutIndex: PreTxOutIndex, PubKey: PubKey, Signature: Signature}
}

func (txIn TxIn) PossessBy(PubKeyH []byte) bool {
	KeyHash := Ripemd160Hash(txIn.PubKey)
	return bytes.Compare(KeyHash, PubKeyH) == 0
}

//验证与前一个交易的公钥哈希是否相同
func (txIn TxIn) VerifyPubKeyHash(preTx Tx) bool{
	PubKeyHash := Ripemd160Hash(txIn.PubKey)
	return bytes.Compare(PubKeyHash, preTx.TxOut[txIn.PreTxOutIndex].PubKeyH) == 0
}

//显示交易输入的信息
func (txIn TxIn) ShowInfo() {
	if txIn.PreTxOutIndex != -1 {
		fmt.Printf("Previous Transaction Hash:\t%x\n", txIn.PreTxH)
		fmt.Printf("Output Index in Previous Transaction:\t%x\n", txIn.PreTxOutIndex)
		fmt.Printf("Signatrue:\t%x\n", txIn.Signature)
		addr := PubKeyHash2Addr(Ripemd160Hash(txIn.PubKey))
		fmt.Printf("Address:\t%s\n", addr)
	} else {
		fmt.Printf("\tNull\n")
	}
}
