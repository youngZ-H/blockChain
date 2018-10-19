package main

import (
	"os"
	"fmt"
	"strconv"
)

//定义命令行结构体
type CLI struct {
	bc *BlockChain
}

//定义一个全局常量，用来存储输入命令的用法
const Usage = `
	printChain               			"正向打印区块"
	printChainR               			"反向打印区块"
	getBalance --address ADDRESS        "获取指定地址的余额"
	send FROM TO AMOUNT MINER DATA		"FROM 发送地址；TO 接收地址；AMMOUNT 转账的金额 MINER 矿工 DATA 矿工添加的信息"
	newWallet							"创建新的钱包，生成一对新的公钥私钥对"
	listAllAddress						"列举出钱包中所有的地址"
`
//为命令行结构体绑定一个函数
func (cli *CLI) Run() {
	args := os.Args
	//对输入的命令进行判断
	if len(args) < 2 {
		fmt.Println("请输入正确的命令")
		fmt.Println(Usage)
		return
	}
	//利用switch进行输入的命令的判断
	cmd := args[1]
	switch cmd {
	/*case "addBlock":
		fmt.Println("添加区块")
		//判断添加数据的条件
		if len(args) == 4 && args[2] == "--data" {
			//进行数据的添加
			data := args[3]
			cli.AddBlock(data)
		}*/
	case "printChain":
		fmt.Println("打印区块")
		//cli.PrintBlockChain()
		cli.bc.PrintBlockChain()
	case "printChainR":
		fmt.Println("反向打印区块")
		cli.PrintBlockChainReverse()
	case "getBalance":
		fmt.Println("获取余额")
		//获取打印余额的命令行条件
		if len(args) == 4 && args[2] == "--address" {
			//进行数据的添加
			address := args[3]
			cli.GetBalance(address)
		}
	case "send":
		fmt.Println("进行转账")
		if len(args) != 7{
			fmt.Println("输入命令有误，请重新输入")
			fmt.Println(Usage)
			return
		}
		//获取输入的各个参数
		//send FROM TO AMOUNT MINER DATA		"FROM 发送地址；TO 接收地址；AMMOUNT 转账的金额 MINER 矿工 DATA 矿工添加的信息"
		from := args[2]
		to := args[3]
		amount,_ := strconv.ParseFloat(args[4],64)
		miner := args[5]
		data := args[6]
		cli.Send(from, to, amount, miner, data)
	case "newWallet":
		fmt.Println("生成新的私钥公钥对")
		cli.NewWallet()
	case "listAllAddress":
		fmt.Println("打印所有的地址......")
		cli.ListAllAddress()


	default:
		fmt.Println("输入的命令错误")
		fmt.Println(Usage)
	}
}
