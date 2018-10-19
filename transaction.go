package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
)

//定义交易结构体
type Transaction struct {
	TXID []byte 	//	交易ID
	TXInputs []TXInput //交易输入数组
	TXOutputs []TXOutput //交易输出数组
}
//定义交易输入
type TXInput struct {
	TXid []byte //交易id
	Index int64 //索引值
	Sig string //解锁脚本，我们用地址模拟
}
//定义交易输出
type TXOutput struct {
	Value float64 //交易金额
	PubKeyHash string //锁定脚本，我们用地址模拟
}

const reward  = 50

//绑定提供hash值的方法
func (tx *Transaction)SetHash()  {
	//定义一个缓存区用来存储编码
	var buffer bytes.Buffer
	//定义一个编码器
	encoder := gob.NewEncoder(&buffer)
	//进行编码
	err := encoder.Encode(tx)
	if err != nil{
		log.Panic(err)
	}
	//对编码后的数据进行hash运算
	//将编码后的数据赋值给一个变量
	data := buffer.Bytes()
	//进行sha256编码
	hash := sha256.Sum256(data)
	//将交易编码后的值赋值给交易id
	tx.TXID = hash[:]

}
//提供创建交易的方法，挖矿交易
func NewCoinbaseTx(address, data string )*Transaction  {
	//挖矿交易，没有输入来源，所以没有交易id，没有索引值，解锁脚本为挖矿塞入的信息，转账金额为挖矿的所得值
	//交易输入
	input := TXInput{[]byte{},-1,data}
	//交易输出
	output := TXOutput{reward,address}
	tx := Transaction{[]byte{},[]TXInput{input},[]TXOutput{output}}
	//对挖矿交易进行hash运算
	tx.SetHash()
	return &tx

}
//实现一个函数，判断当前交易是否为挖矿交易
func (tx *Transaction)IsCoinbaseTx()bool {
	//如果时挖矿交易，交易id为空，值存在一个input，无索引值
	//if len(tx.TXInputs) == 1{
	//	input := tx.TXInputs[0]
	//	if bytes.Equal(tx.TXID,[]byte{}) && input.Index == -1{
	//		return true
	//	}
	//}
	//return false
	if len(tx.TXInputs) == 1 && len(tx.TXInputs[0].TXid) == 0 && tx.TXInputs[0].Index == -1 {
		return true
	}
	return false
}
//创建普通交易
func NewTransaction(from, to string, amount float64, bc *BlockChain)*Transaction  {
	//普通交易需要找到最优的UTXO
	utxos, resValue := bc.FindNeedUTXOs(from,amount)
	//判断余额是否充足
	if resValue < amount{
		fmt.Println("余额不足")
		return nil
	}
	//创建交易输入和交易输出
	var inputs []TXInput
	var outputs []TXOutput
	//创建交易输入，遍历utxos
	for id, indexArray := range utxos{
		//遍历utxos获得交易id和交易索引值的数组，在遍历交易索引值数组
		for _,index := range indexArray{
			input := TXInput{[]byte(id),index,from}
			inputs = append(inputs,input)
		}
	}

	//创建交易输出
	output := TXOutput{amount,to}
	outputs = append(outputs,output)
	//判断一下是否要找零
	if resValue > amount{
		outputs = append(outputs,TXOutput{resValue-amount,from })
		//fmt.Println("=================",resValue - amount)
		//fmt.Println("=======================",outputs)
	}
	//将数据加入到交易的结构体中
	tx := Transaction{[]byte{},inputs,outputs}
	//将交易的信息进行hash运算
	tx.SetHash()
	return &tx
}