package main

import "math"

const bucket = "blockchain"         //表单名
const nBits = 16                    //难度
const MaxNonce = math.MaxInt64      //最大数
const Reward = 12.5                 //挖矿奖励
const dbFile = "./db/blockchain.db" //数据库文件
const walletFile = "./db/wallet.db" //钱包文件
const Walletbucket = "Wallet"       //钱包表单名
const version = byte(0x00)          //定义版本号，一个字节
const addressChecksumLen = 4        //定义checksum长度为四个字节
