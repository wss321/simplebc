package main

//命令行
import (
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"strconv"
)

func init() {
	gob.Register(elliptic.P256())
}
func Command() {
	//实例化cli
	app := cli.NewApp()
	//Name可以设定应用的名字
	app.Name = "block chain"
	// Version可以设定应用的版本号
	app.Version = "1.0.0"
	// Commands用于创建命令
	app.Commands = []cli.Command{
		{
			// 命令的名字
			Name: "createblockchain",
			// 命令的缩写，就是不输入language只输入lang也可以调用命令
			Aliases: []string{"createbc", "cbc"},
			// 命令的用法注释，这里会在输入 程序名 -help的时候显示命令的使用方法
			Usage: "Create a Block Chain:\tcreateblockchain <address>",
			// 命令的处理函数
			Action: func(c *cli.Context) error {
				address := c.Args().First()
				if IsValidAddr([]byte(address)) {
					CreateBlockChain(dbFile, string(address))
				} else {
					log.Panic("Invalid address")
				}
				return nil
			},
		},
		{
			Name:    "createwallet",
			Aliases: []string{"createwl", "cw"},
			Usage:   "Create a wallet:\tcreatewallet",
			Action: func(c *cli.Context) error {
				wallet := NewWallet()
				wallet.ShowInfo()
				return nil
			},
		},
		{
			Name:    "getbalance",
			Aliases: []string{"gb"},
			Usage:   "Get balance:\tgetbalance <address>",
			Action: func(c *cli.Context) error {
				address := c.Args().First()
				if IsValidAddr([]byte(address)) {
					bc := LoadBlockChain(dbFile)
					balance := bc.Balance(string(address))
					fmt.Printf("Balance of %s is %v\n", address, balance)
				} else {
					log.Panic("Invalid address")
				}
				return nil
			},
		},

		{
			Name:    "send",
			Aliases: []string{"sd"},
			Usage:   "Send coin to others:\t send <from> <to> <amout>",
			Action: func(c *cli.Context) error {
				from := c.Args().First()
				to := c.Args().Get(1)
				amout := c.Args().Get(2)
				//fmt.Printf("from :%s\nto: %s\n", from, to)
				bc := LoadBlockChain(dbFile)
				wallet := LoadWallet(from)
				amoutFloat, _ := strconv.ParseFloat(amout, 64)
				wallet.Send(to, amoutFloat, bc)
				return nil
			},
		},
		{
			Name:    "showchain",
			Aliases: []string{"show", "sc"},
			Usage:   "Show Block Chain Info:\tshowchain",
			Action: func(c *cli.Context) error {

				bc := LoadBlockChain(dbFile)
				bc.ShowChainInfo()
				return nil
			},
		},
		{
			Name:    "showwallets",
			Aliases: []string{"sw"},
			Usage:   "Show all wallets address:\tshowwallets",
			Action: func(c *cli.Context) error {
				ShowAllWallets()
				return nil
			},
		},
	}
	// 接受os.Args启动程序
	e := app.Run(os.Args)
	if e != nil {
		log.Panic(e)
	}
}
