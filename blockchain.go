package main

import (
	"bytes"
	"fmt"
	"os"
)

//创建新的区块链
func CreateBlockChain(file, receiver string) *Blockchain {
	if isExists(file) == true {

		fmt.Println("BlockChain already exists.")
		os.Exit(1)
	}

	var lastHash []byte
	fmt.Printf("%s Creating GenesisBlock\n", receiver)
	cbtx := CoinBaseTx(receiver, "")
	gb := GenesisBlock(cbtx)
	lastHash = gb.Hash

	dB := CreateDB(file, bucket) //创建数据库
	dB.Put(lastHash, gb.Encode())

	dB.Put([]byte("lastHash"), lastHash)
	dB.CurrentHash = lastHash
	bc := &Blockchain{LastBlockH: lastHash, Db: dB}
	return bc
}

//加载已经创建的数据库中的区块链
func LoadBlockChain(file string) *Blockchain {
	db := &DataBase{File: file, Bucket: bucket}
	db.LoadData()
	lastHash := db.Get("lastHash")
	db.CurrentHash = lastHash
	//fmt.Printf("last hash :%x", lastHash)
	bc := &Blockchain{LastBlockH: lastHash, Db: db}
	return bc
}

//遍历区块链的函数
func (bc *Blockchain) BlockIter() *Block {
	return bc.Db.BlockIter()
}

//重置当前hash至最后block
func (bc *Blockchain) ResetCurrent() {
	bc.Db.ResetCurrent()
}

//关闭区块链文件
func (bc *Blockchain) Close() {
	bc.Db.Close()
}

//添加区块
func (bc *Blockchain) AddBlock(bl *Block) {
	hash := bl.Hash
	bc.Db.Put([]byte("lastHash"), hash)
	bc.Db.Put(hash, bl.Encode())
	bc.Db.CurrentHash = hash
	bc.LastBlockH = hash
}

//找到某人未花费的UTXO
func (bc Blockchain) FindUnspentUTXOs(PubKeyHash []byte) ([]TxOut, [][]byte, []int) {
	var unspentUTXOs []TxOut
	var txHashes [][]byte
	var txOutIndex []int

	spentUTXOIndexes := make(map[string][]int)

	//找出已花费的UTXO
	currentHash := bc.Db.CurrentHash
	bc.ResetCurrent()
	for {
		block := bc.BlockIter() //遍历区块
		txs := block.Txs
		//println(len(txs))
		for _, tx := range txs { //遍历交易
			for _, txIn := range tx.TxIn { //交易输入都认为是已花费的UTXO
				if IsCoinBaseTx(*tx) == false {
					//println(false)
					spentUTXOIndexes[string(txIn.PreTxH)] = append(spentUTXOIndexes[string(txIn.PreTxH)], txIn.PreTxOutIndex)
				}
				//println(true)
			}

		}
		if len(block.Header.PreH) == 0 {
			break
		}
	}
	//for i, j := range spentUTXOIndexes {
	//	fmt.Printf("%x\t%v\n", i, j)
	//}
	//再次遍历，找出未花费的UTXO
	bc.ResetCurrent()
	for {
		block := bc.BlockIter() //区块
		txs := block.Txs
		for _, tx := range txs { //交易
			txhash := tx.Hash()
			for index, txOut := range tx.TxOut { //UTXO
				if txOut.PossessBy(PubKeyHash) {
					spendIndexes := spentUTXOIndexes[string(txhash)]
					unspend := true
					if spendIndexes != nil {
						for _, spendIndex := range spendIndexes {
							if index == spendIndex {
								unspend = false
								break
							}
						}
					}
					if unspend {
						unspentUTXOs = append(unspentUTXOs, txOut)
						txHashes = append(txHashes, txhash)
						txOutIndex = append(txOutIndex, index)
					}
				}
			}

		}
		if len(block.Header.PreH) == 0 {
			break
		}
	}
	bc.Db.CurrentHash = currentHash //还原当前hash
	return unspentUTXOs, txHashes, txOutIndex

}

//计算某人的余额
func (bc Blockchain) Balance(addr string) float64 {
	amount := 0.0
	pubKeyHash := Addr2PubKeyHash(addr)
	UTXOs, _, _ := bc.FindUnspentUTXOs(pubKeyHash)
	for _, utxo := range UTXOs {
		amount += utxo.Value
	}
	return amount
}

//找出某个交易
func (bc Blockchain) FindTx(txHash []byte) *Tx {
	currentHash := bc.Db.CurrentHash
	bc.ResetCurrent()
	for {
		block := bc.BlockIter() //区块
		for _, tx := range block.Txs {
			//fmt.Printf("hash1:%x\nhash2:%x\n", tx.Hash(), txHash)
			if bytes.Compare(tx.Hash(), txHash) == 0 {
				bc.Db.CurrentHash = currentHash //还原当前hash
				return tx
			}
		}
		if len(block.Header.PreH) == 0 {
			break
		}
	}
	bc.Db.CurrentHash = currentHash //还原当前hash
	return nil

}

//显示bc信息
func (bc Blockchain) ShowChainInfo() {
	currentHash := bc.Db.CurrentHash
	bc.ResetCurrent()
	for {
		lastBlock := bc.BlockIter()
		lastBlock.ShowBlockInfo()
		for _, tx := range lastBlock.Txs {
			tx.ShowInfo()
		}
		if len(lastBlock.Header.PreH) == 0 {
			break
		}
	}
	bc.Db.CurrentHash = currentHash //还原当前hash
}
