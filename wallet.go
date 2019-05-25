package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/boltdb/bolt"
	"golang.org/x/crypto/ripemd160"
	"log"
	"os"
)

//判断地址是否有效
func IsValidAddr(addr []byte) bool {
	//将地址进行base58反编码，生成的其实是version+Pub Key hash+ checksum这25个字节
	versionPublicChecksum := Base58Decode(string(addr))

	//[25-4:],就是21个字节往后的数（22,23,24,25一共4个字节）
	checkSumBytes := versionPublicChecksum[len(versionPublicChecksum)-addressChecksumLen:]
	//[:25-4],就是前21个字节（1～21,一共21个字节）
	versionRipemd160 := versionPublicChecksum[:len(versionPublicChecksum)-addressChecksumLen]
	//取version+public+checksum的字节数组的前21个字节进行两次256哈希运算，取结果值的前4个字节
	checkBytes := CheckSum(versionRipemd160)
	//将checksum比较，如果一致则说明地址有效，返回true
	if bytes.Compare(checkSumBytes, checkBytes) == 0 {
		return true
	}

	return false
}

func (w Wallet) GetAddress() []byte {

	//调用Ripemd160Hash返回160位的Pub Key hash
	ripemd160Hash := Ripemd160Hash(w.PublicKey)

	//将version+Pub Key hash
	versionRipemd160hash := append([]byte{version}, ripemd160Hash...)

	//调用CheckSum方法返回前四个字节的checksum
	checkSumBytes := CheckSum(versionRipemd160hash)

	//将version+Pub Key hash+ checksum生成25个字节
	bts := append(versionRipemd160hash, checkSumBytes...)

	//将这25个字节进行base58编码并返回
	return []byte(Base58Encode(bts))
}

//向某人发送货币
func (w Wallet) Send(receiver string, amount float64, bc *Blockchain) {
	if IsValidAddr([]byte(receiver)) == false {
		log.Panic(fmt.Sprintf("Invalid addr %s", receiver))
	}
	tx := Send(receiver, string(w.GetAddress()), amount, bc, w)
	//fmt.Printf("%s\n",Timestamp2Time(tx.LockTime))
	cbtx := CoinBaseTx(string(w.GetAddress()), "") //发送者挖矿
	fmt.Printf("\r%s mining Block\n", string(w.GetAddress()))
	MineBlock([]*Tx{tx, cbtx}, bc)	//挖矿
	fmt.Printf("%s Sended %v Coins To %s.\n", w.GetAddress(), amount, receiver)
}

//取前4个字节
func CheckSum(payload []byte) []byte {
	//这里传入的payload其实是version+Pub Key hash，对其进行两次256运算
	hash1 := sha256.Sum256(payload)

	hash2 := sha256.Sum256(hash1[:])

	return hash2[:addressChecksumLen] //返回前四个字节，为CheckSum值
}

func Ripemd160Hash(publicKey []byte) []byte {

	//将传入的公钥进行256运算，返回256位hash值
	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)

	//将上面的256位hash值进行160运算，返回160位的hash值
	ripemd160hash := ripemd160.New()
	ripemd160hash.Write(hash)

	return ripemd160hash.Sum(nil) //返回Pub Key hash
}

// 创建钱包
func NewWallet() *Wallet {

	privateKey, publicKey := newKeyPair()
	wallet := &Wallet{privateKey, publicKey}
	wallet.Save()
	return wallet
}

// 通过私钥产生公钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	//这是一个曲线对象
	curve := elliptic.P256()
	//通过椭圆曲线加密算法生成私钥
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	//由私钥生成公钥
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

func Addr2PubKeyHash(addr string) []byte {
	pubKeyHash := Base58Decode(addr)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	return pubKeyHash
}

func PubKeyHash2Addr(b []byte) string {
	b = append([]byte{version}, b...)
	checkSumBytes := CheckSum(b)

	//将version+Pub Key hash+ checksum生成25个字节
	bts := append(b, checkSumBytes...)
	return Base58Encode(bts)
}

//Wallet转byte数组，用于保存到数据库
func (w Wallet) Encode() []byte {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(w)
	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

//byte数组转Wallet，便于从数据库中读取
func ByteArr2Wallet(bt []byte) *Wallet {
	var wallet Wallet

	decoder := gob.NewDecoder(bytes.NewReader(bt))
	err := decoder.Decode(&wallet)
	if err != nil {
		log.Panic(err)
	}

	return &wallet
}

//保存钱包
func (w Wallet) Save() {
	if isExists(walletFile) == false {
		dB := CreateDB(walletFile, Walletbucket) //创建数据库
		dB.Put(w.GetAddress(), w.Encode())
		dB.Close()
	}
	db := &DataBase{File: walletFile, Bucket: Walletbucket}
	db.LoadData()
	db.Put(w.GetAddress(), w.Encode())
	db.Close()
}

//取出钱包
func LoadWallet(address string) *Wallet {
	db := &DataBase{File: walletFile, Bucket: Walletbucket}
	db.LoadData()
	return ByteArr2Wallet(db.Get(address))
}

//保存钱包
func (w Wallet) ShowInfo() {
	fmt.Printf("Address:\t%s\n", string(w.GetAddress()))
}

//显示所有钱包
func ShowAllWallets() {
	if isExists(walletFile) == false {
		fmt.Printf("Wallet File Doesn't Exist.")
		os.Exit(1)
	}
	dB := &DataBase{File: walletFile, Bucket: Walletbucket}
	db := dB.LoadData()
	defer db.Close()

	err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(Walletbucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			fmt.Printf("Wallet:\t%s\n", k)
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}
