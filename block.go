package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/cbergoon/merkletree"
	"log"
	"time"
)

func NewBlock(txs []*Tx, preBH []byte) *Block {
	if len(txs)==0{
		log.Panic("Number of Transaction can't be 0, while creating Block.")
	}
	blockHeader := &BlockHeader{PreH: preBH, Timestamp: time.Now().Unix(), NBits: nBits}
	block := &Block{Header: blockHeader, TxCnt: uint(len(txs)), Txs: txs}
	block.Header.MerkleRootHash = block.MerkleRootHash()
	pow := &PoW{Block: block, NBits: nBits}

	nonce, blockHash := pow.FindNonce()
	block.Hash = blockHash
	block.Header.Nonce = uint(nonce)
	return block
}

//创世区块
func GenesisBlock(coinBase *Tx) *Block {
	return NewBlock([]*Tx{coinBase}, []byte{})
}

//计算全部交易的hash
//创建一个符合merkletree的类型
type TxMT struct {
	tx *Tx
}

func (txc TxMT) CalculateHash() ([]byte, error) {
	return txc.tx.Hash(), nil
}

func (txc TxMT) Equals(other merkletree.Content) (bool, error) {
	return txc.tx == other.(*TxMT).tx, nil
}

//树根hash
func (b *Block) MerkleRootHash() []byte {
	var txsHash []merkletree.Content
	for _, tx := range b.Txs {
		txsHash = append(txsHash, TxMT{tx})
	}

	MT, err := merkletree.NewTree(txsHash)
	if err != nil {
		log.Panic(err)
	}
	return MT.Root.Hash
}

//Block转byte数组，用于保存到数据库
func (b Block) Encode() []byte {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

//byte数组转Block，便于从数据库中读取
func ByteArr2Block(bt []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(bt))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

//输出Block的信息
func (b Block) ShowBlockInfo() {
	fmt.Printf("\n-------------  Block Infomation  ---------------- \n")
	fmt.Printf("\nNumber Of Transactions:\t%v \n", b.TxCnt)
	fmt.Printf("Timestamp:\t%s \n", Timestamp2Time(b.Header.Timestamp))
	cbtx := b.FindCoinBaseTx()
	fmt.Printf("Created By:\t%s \n", PubKeyHash2Addr(cbtx.TxOut[0].PubKeyH))
	fmt.Printf("Nonce:\t%v\n", b.Header.Nonce)
	fmt.Printf("Difficulty:\t%v bits\n", b.Header.NBits)
	fmt.Printf("Block Reward:\t%v\n", cbtx.TxOut[0].Value)
	fmt.Printf("Hash:\t%x\n", b.Hash)
	fmt.Printf("Previous Block:\t%x\n", b.Header.PreH)
	fmt.Printf("Merkle Root:\t%x\n", b.Header.MerkleRootHash)
	fmt.Printf("\n-------------------------------------------------- \n")

}

//找出CoinBase交易
func (b Block) FindCoinBaseTx() *Tx {
	for _, tx := range b.Txs {
		if IsCoinBaseTx(*tx) {
			return tx
		}
	}
	return nil
}
