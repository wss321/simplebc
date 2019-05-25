package main

import "crypto/ecdsa"

//区块链
type Blockchain struct {
	LastBlockH []byte //最后区块的hash
	Db         *DataBase
}
type BlockHeader struct {
	//version int32  //版本号
	PreH           []byte //前一块hash
	MerkleRootHash []byte //Merkle 树根的hash值
	Timestamp      int64  //时间戳
	NBits          uint   //挖矿难度，多少bit
	Nonce          uint
}

// 区块结构
type Block struct {
	//size   uint         //后面数据到块结束的字节数
	Hash   []byte
	Header *BlockHeader //区块头
	TxCnt  uint         //交易数量
	Txs    []*Tx        //交易
}

//交易
type Tx struct {
	TxInCnt  uint    //输入数量
	TxIn     []TxIn  //交易输入
	TxOutCnt uint    //输出数量
	TxOut    []TxOut //交易输出
	LockTime int64   //锁定时间
}

//交易输入
type TxIn struct {
	PreTxH        []byte //前置交易hash
	PreTxOutIndex int    //处在前置交易输出中的index
	Signature     []byte //签名
	PubKey        []byte //所有者公钥
}

//交易输出
type TxOut struct {
	Value   float64 //花费的数量
	PubKeyH []byte  //对方的公钥哈希
}

//工作量证明
type PoW struct {
	NBits uint //目标bits
	Block *Block
}

// 钱包
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}
