package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

type DataBase struct {
	File        string   //文件名
	Bucket      string   //表单名
	DB          *bolt.DB //读到的数据
	CurrentHash []byte   //当前区块的Hash
}

func CreateDB(file, bucket string) *DataBase {
	if isExists(file) {
		fmt.Printf("%s already exists.", file)
		os.Exit(-1)
	}
	bdb, err := bolt.Open(file, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	db := &DataBase{File: file, DB: bdb}
	db.CreateBucket(bucket)
	db.Bucket = bucket
	return db

}
func (db *DataBase) LoadData() *bolt.DB {
	dbd, err := bolt.Open(db.File, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	db.DB = dbd
	return dbd
}

//获取表单中的数据
func (db *DataBase) Get(key string) []byte {
	var value []byte
	err := db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.Bucket)) //获取表单
		value = b.Get([]byte(key))        //获取表中key对应的数据
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return value
}

//创建表单
func (db *DataBase) CreateBucket(bucket string) {
	err := db.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

//写入表单
func (db *DataBase) Put(key []byte, value []byte) []byte {
	err := db.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.Bucket)) //获取表单
		err := b.Put(key, value)          //写入表单
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return value
}

//上一个Block
func (db *DataBase) BlockIter() *Block {
	blockByteArr := db.Get(string(db.CurrentHash))
	block := ByteArr2Block(blockByteArr)
	db.CurrentHash = block.Header.PreH
	return block
}

//重置当前hash到最后一个区块
func (db *DataBase) ResetCurrent() {
	db.CurrentHash = db.Get("lastHash")
}

func (db *DataBase) Close() {
	db.DB.Close()

}
