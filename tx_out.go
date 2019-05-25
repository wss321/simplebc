package main

import (
	"bytes"
	"fmt"
)

//创建交易输出
func NewTxOut(receiver string, value float64) TxOut {
	txO := TxOut{Value: value, PubKeyH: Addr2PubKeyHash(receiver)}
	return txO
}

func (txOut TxOut) PossessBy(pubKeyH []byte) bool {
	return bytes.Compare(txOut.PubKeyH, pubKeyH) == 0
}
func (txOut TxOut) ShowInfo() {
	fmt.Printf("Value:\t%v\n", txOut.Value)
	fmt.Printf("Address:\t%s\n", PubKeyHash2Addr(txOut.PubKeyH))
}
