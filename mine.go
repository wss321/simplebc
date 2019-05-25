package main

//挖区块
func MineBlock(txs []*Tx, bc *Blockchain) *Block {
	var txVerified []*Tx
	for _, tx := range txs {
		if bc.VerifyTx(*tx) == true {
			txVerified = append(txVerified, tx)
		}
	}
	block := NewBlock(txVerified, bc.LastBlockH)
	bc.AddBlock(block) //添加区块到区块链上
	return block
}


