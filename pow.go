package main

import (
	"bytes"
	"fmt"
)

// 将Block的数据打包以方便进一步做hash运算
func (pow PoW) BlockData2ByteArr(nonce int64) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.Header.PreH,
			pow.Block.MerkleRootHash(),
			Int2ByteArr(int64(pow.Block.TxCnt)),
			Int2ByteArr(int64(pow.Block.Header.Timestamp)),
			Int2ByteArr(int64(pow.Block.Header.NBits)),
			Int2ByteArr(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

//寻找合适的nonce
func (pow PoW) FindNonce() (int64, []byte) {
	var hash []byte
	nonce := int64(0)
	for nonce < MaxNonce {
		hash = CalcHash(pow.BlockData2ByteArr(nonce))
		hashInt := ByteArr2Int(hash)

		hashStr := fmt.Sprintf("%0b", hashInt) //转二进制的字符数组
		bits0 := 8*len(hash)-len(hashStr)
		//fmt.Printf("\rMining Block:\tnonce=%v\thash:%x\tCnt:%v", nonce, hash, bits0)
		//判断是否符合目标
		if bits0 >= int(pow.NBits) {
			break
		} else {
			nonce++
		}
	}
	return nonce, hash

}
