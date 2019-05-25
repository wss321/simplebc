package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

//计算hash
func CalcHash(b []byte) []byte {
	h := sha256.Sum256(b)
	return h[:]
}

func Int2ByteArr(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func ByteArr2Int(buf []byte) *big.Int {
	var hashInt big.Int

	return hashInt.SetBytes(buf)
}

//文件是否存在
func isExists(f string) bool {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return false
	}

	return true
}

//时间戳转时间
func Timestamp2Time(ts int64) string {
	format := time.Unix(int64(ts), 0).Format("2006-01-02 15:04:05")
	return format
}

//base58编码
// alphabet is the modified base58 alphabet used by Bitcoin.
const BTCAlphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

//const FlickrAlphabet = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

var bigRadix = big.NewInt(58)
var bigZero = big.NewInt(0)

// Decode decodes a modified base58 string to a byte slice, using BTCAlphabet
func Base58Decode(b string) []byte {
	return DecodeAlphabet(b, BTCAlphabet)
}

// Encode encodes a byte slice to a modified base58 string, using BTCAlphabet
func Base58Encode(b []byte) string {
	return EncodeAlphabet(b, BTCAlphabet)
}

// DecodeAlphabet decodes a modified base58 string to a byte slice, using alphabet.
func DecodeAlphabet(b, alphabet string) []byte {
	answer := big.NewInt(0)
	j := big.NewInt(1)

	for i := len(b) - 1; i >= 0; i-- {
		tmp := strings.IndexAny(alphabet, string(b[i]))
		if tmp == -1 {
			return []byte("")
		}
		idx := big.NewInt(int64(tmp))
		tmp1 := big.NewInt(0)
		tmp1.Mul(j, idx)

		answer.Add(answer, tmp1)
		j.Mul(j, bigRadix)
	}

	tmpval := answer.Bytes()

	var numZeros int
	for numZeros = 0; numZeros < len(b); numZeros++ {
		if b[numZeros] != alphabet[0] {
			break
		}
	}
	flen := numZeros + len(tmpval)
	val := make([]byte, flen, flen)
	copy(val[numZeros:], tmpval)

	return val
}

// Encode encodes a byte slice to a modified base58 string, using alphabet
func EncodeAlphabet(b []byte, alphabet string) string {
	x := new(big.Int)
	x.SetBytes(b)

	answer := make([]byte, 0, len(b)*136/100)
	for x.Cmp(bigZero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, bigRadix, mod)
		answer = append(answer, alphabet[mod.Int64()])
	}

	// leading zero bytes
	for _, i := range b {
		if i != 0 {
			break
		}
		answer = append(answer, alphabet[0])
	}

	// reverse
	alen := len(answer)
	for i := 0; i < alen/2; i++ {
		answer[i], answer[alen-1-i] = answer[alen-1-i], answer[i]
	}

	return string(answer)
}
//
////复制结构体
//func CopyStruct(src, dst interface{}) {
//	sval := reflect.ValueOf(src).Elem()
//	dval := reflect.ValueOf(dst).Elem()
//
//	for i := 0; i < sval.NumField(); i++ {
//		value := sval.Field(i)
//		name := sval.Type().Field(i).Name
//
//		dvalue := dval.FieldByName(name)
//		if dvalue.IsValid() == false {
//			continue
//		}
//		dvalue.Set(value) //这里默认共同成员的类型一样，否则这个地方可能导致 panic，需要简单修改一下。
//	}
//}
