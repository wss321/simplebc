/*
1.验证交易
2.工作量证明验证
3.
*/

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//验证交易
func (bc *Blockchain) VerifyTx(tx Tx) bool {
	if IsCoinBaseTx(tx) {
		return true
	}

	txCp := tx
	curve := elliptic.P256()

	for i, txIn := range txCp.TxIn {
		prevTx := bc.FindTx(txIn.PreTxH)
		if prevTx == nil {
			break
		}
		if txCp.TxIn[i].VerifyPubKeyHash(*prevTx) == false {
			return false
		}
		txCp.TxIn[i].Signature = nil
		//签名
		hash := txCp.Hash()
		//txCp.TxIn[i].PubKey = nil
		r := big.Int{}
		s := big.Int{}
		sigLen := len(txIn.Signature)

		r.SetBytes(txIn.Signature[:(sigLen / 2)])
		s.SetBytes(txIn.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(txIn.PubKey)
		x.SetBytes(txIn.PubKey[:(keyLen / 2)])
		y.SetBytes(txIn.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}

		if ecdsa.Verify(&rawPubKey, hash, &r, &s) == false {
			return false
		}

	}

	return true
}

//验证工作量
func (pow PoW) VerifyPoW() bool {
	data := pow.BlockData2ByteArr(int64(pow.Block.Header.Nonce))
	var hashInt big.Int
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	hashStr := fmt.Sprintf("%0b", hashInt) //转二进制的字符数组
	bits0 := 8*len(hash) - len(hashStr)
	//判断是否符合目标
	return bits0 >= int(pow.NBits)

}

//验证区块
//func (bc Blockchain) VarifyBlock(b Block) bool {
//
//	return false
//}
