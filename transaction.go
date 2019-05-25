package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"
)

//计算一个交易的hash(签名之前的)
func (tx Tx) Hash() []byte {
	return CalcHash(tx.ToByteArr())
}

//交易类型转byte数组，以便计算hash
func (tx Tx) ToByteArr() []byte {
	var buf bytes.Buffer
	txCp := tx
	if IsCoinBaseTx(txCp) == false {
		txCp.LockTime = 0
	}
	for _, txIn := range txCp.TxIn {
		txIn.Signature = nil
	}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(txCp)
	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

//对交易签名
func (tx *Tx) Sign(privKey ecdsa.PrivateKey, bc *Blockchain) {
	if IsCoinBaseTx(*tx) {
		return
	}

	txCp := *tx
	for i, txIn := range txCp.TxIn {
		prevTx := bc.FindTx(txIn.PreTxH)
		if prevTx == nil {
			break
		}
		txCp.TxIn[i].Signature = nil
		//签名
		hash := txCp.Hash()
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, hash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.TxIn[i].Signature = signature
	}
	tx.LockTime = time.Now().Unix()

}

//输出交易信息
func (tx Tx) ShowInfo() {
	fmt.Printf("\n-------------  Transaction Infomation  ---------------- \n")
	fmt.Printf("Hash:\t%x \n", tx.Hash())
	fmt.Printf("Lock Time:\t%s \n", Timestamp2Time(tx.LockTime))
	fmt.Printf("\nNumber Of Input:\t%v \n", tx.TxInCnt)
	for i, txIn := range tx.TxIn {
		fmt.Printf("\n\t----  Info of TxInput %v  --- \n", i)
		txIn.ShowInfo()
	}

	fmt.Printf("\nNumber Of Output:\t%v \n", tx.TxOutCnt)
	for i, txOut := range tx.TxOut {
		fmt.Printf("\n\t----Info of TxOutput %v  --- \n", i)
		txOut.ShowInfo()
	}
	fmt.Printf("\n-------------------------------------------------- \n")
}

//创建CoinBase交易
func CoinBaseTx(receiver, disc string) *Tx {
	if disc == "" {
		disc = fmt.Sprintf("Create By '%s'", receiver)
	}
	txIn := NewTxIn([]byte{}, nil, []byte(receiver), -1)
	txOut := NewTxOut(receiver, Reward)
	tx := NewTx([]TxIn{txIn}, []TxOut{txOut})
	return tx
}

//计算交易输入和输出数量
func (tx *Tx) setIOCnt() {
	var txInCnt = uint(0)
	var txOutCnt uint
	for _, txI := range tx.TxIn {
		if txI.PreTxOutIndex == -1 && txI.Signature == nil && len(txI.PreTxH) == 0 {
			continue
		}
		txInCnt++
	}
	txOutCnt = uint(len(tx.TxOut))
	tx.TxInCnt = txInCnt
	tx.TxOutCnt = txOutCnt
}

func IsCoinBaseTx(tx Tx) bool {
	return tx.TxInCnt == 0 && tx.TxOutCnt == 1
}

//创建交易
func NewTx(txIns []TxIn, txOuts []TxOut) *Tx {
	tx := &Tx{TxIn: txIns, TxOut: txOuts, LockTime: time.Now().Unix()}
	tx.setIOCnt()
	return tx
}

//sender发送数字币给receiver
func Send(receiver, sender string, amout float64, bc *Blockchain, w Wallet) *Tx {
	/*
		1.找到未花费的UTXO查看够不够钱

		2.把未花费的UTXO作为输入，接收方作为输出产生交易

		3.交易签名
	*/
	var spandingUtxos []TxOut //即将用于消费的UTXO
	var txHashs [][]byte
	var txOutIndexes []int
	var inputs []TxIn
	var outputs []TxOut
	values := 0.0
	unspendUTXOs, txHashs, txOutIndexes := bc.FindUnspentUTXOs(Addr2PubKeyHash(sender))
	for _, utxo := range unspendUTXOs {
		values += utxo.Value
		spandingUtxos = append(spandingUtxos, utxo)
		if values >= amout {
			break
		}
	}
	if values < amout {
		fmt.Printf("%s Not enough Balance, Remain %v", sender, values)
		os.Exit(1)
	}

	//创建输入
	for i, _ := range spandingUtxos {
		txi := NewTxIn(txHashs[i], nil, w.PublicKey, txOutIndexes[i])
		inputs = append(inputs, txi)
	}
	//创建输出
	outputs = append(outputs, NewTxOut(receiver, amout))
	if amout < values {
		outputs = append(outputs, NewTxOut(sender, values-amout))
	}
	//fmt.Printf("outputs %v", len(outputs))
	tx := NewTx(inputs, outputs)
	tx.Sign(w.PrivateKey, bc) //签名
	return tx
}
