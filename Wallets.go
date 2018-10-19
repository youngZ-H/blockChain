package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"io/ioutil"
	"crypto/elliptic"
	"os"
)
//定义常量
const wallet  = "wallet.dat"
//创建一个wallets用来存储wallet和address
//定义一个钱包的结构体
type Wallets struct {
	WalletsMap map[string]*Wallet
}
//新建一个wallets的方法
func NewWallets()*Wallets  {
	var ws Wallets
	ws.WalletsMap = make(map[string]*Wallet)
	ws.LoadFile()
	return &ws
}
//创建钱包的方法
func (ws *Wallets)CreateWallet()string{
	//新创建一个钱包
	wt := NewWallet()
	//新生成一个地址
	address := wt.NewAddress()
	//定义一个存储钱包的结构体
	//var ws Wallets
	//开辟一个新的map，然后给map赋值，返回钱包
	//ws.WalletsMap = make(map[string]*Wallet)
	ws.WalletsMap[address] = wt
	ws.SaveToFile()
	return address
}
//为wallets绑定保存到文件中的方法
func (ws *Wallets)SaveToFile()  {
	//定义一个缓存区用阿里存储编码后的数据
	var buffer bytes.Buffer
	//gob编码前需要先进行注册
	gob.Register(elliptic.P256())
	//定义一个编码器
	enconder := gob.NewEncoder(&buffer)
	//进行编码
	err := enconder.Encode(ws)
	if err != nil{
		log.Panic(err)
	}
	//将编码后的收取写入wallet.dat文件中
	ioutil.WriteFile(wallet,buffer.Bytes(),0600)
}

func (ws *Wallets)LoadFile()  {
	//打开文件之前，首先判断一下文件是否存在
	_, err := os.Stat(wallet)
	//fmt.Println("22222222222222222222222222")
	//fmt.Println(err)
	if os.IsNotExist(err){
		//fmt.Println("33333333333333333333333")
		return
	}
	//打开文件,进行文件的读取
	content, err := ioutil.ReadFile(wallet)
	if err != nil{
		log.Panic(err)
	}
	//将读取到的内容进行解码
	//进行gob的注册，新建解码器，进行解码
	gob.Register(elliptic.P256())
	//创建解码器
	deconder := gob.NewDecoder(bytes.NewReader(content))
	//定义一个本地的钱包结构，进行本地读取的信息进行内存的存储
	var wsLocal Wallets
	err = deconder.Decode(&wsLocal)
	if err != nil{
		log.Panic(err)
	}
	//进行钱包结构的赋值，注意：结构体中含有map的，不能直接从外部对结构体进行赋值，要单独提出直接对map进行赋值
	ws.WalletsMap = wsLocal.WalletsMap
}
//列举钱包中存放的所有的地址
func (ws *Wallets)ListAllAddress()[]string  {
	//定义一个字符创切片用来存储取出的所有的地址
	var addresses []string
	//遍历map，取出所有的地址
	for adrs,_ := range ws.WalletsMap{
		//将所有的地址存储起来
		addresses = append(addresses,adrs)
	}
	return addresses
}




























