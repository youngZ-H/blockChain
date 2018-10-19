package main

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

//定义工作量证明结构体
type ProofOfWork struct {
	//指定区块
	block *Block
	//制定目标难度值
	target *big.Int
}
//定义新建工作量证明的方法
func NewProofWork(block *Block)*ProofOfWork {
	//定义一个工作量证明的对象
	pow := ProofOfWork{
		block:block,
	}
	//指定一个难度值,字符串类型
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"
	//定义一个辅助变量，将上面的难度值转换为big.int类型
	tmpInt := big.Int{}
	tmpInt.SetString(targetStr,16)
	pow.target = &tmpInt
	return &pow
}
//为工作量证明绑定数据准备的方法
func (pow *ProofOfWork)PrepareData(nonce uint64)[]byte  {
	block := pow.block
			//进行区块数据的拼装
		tmp := [][]byte{
		Uint64ToByte(block.Version),
		block.PreHash,
		block.MerKleRoot,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(nonce),
		//block.Data,
		//只对区块头做hash运算，区块体通过梅克尔根影响
	}
		//进行数据的拼装
		blockInfo := bytes.Join(tmp,[]byte{})
		return blockInfo
}
//为工作量证明绑定不断计算哈希的函数，即进行挖矿的过程
func (pow *ProofOfWork)Run()([]byte,uint64) {
	//1. 拼装数据（区块的数据，还有不断变化的随机数）
	//2. 做哈希运算
	//3. 与pow中的target进行比较
	//a. 找到了，退出返回
	//b. 没找到，继续找，随机数加1


	//定义随机数
	var nonce uint64
	//定义hash
	var hash [32]byte

	//1.根据不同的随机数，进行循环计算哈希值的运算,直到条件满足退出循环，返回随机数和哈希值
	fmt.Println("正在挖矿......")
	for{
		//调用工作量证明的准备函数
		blockInfo := pow.PrepareData(nonce)
		//进行sha256的编码运算
		hash = sha256.Sum256(blockInfo)
		//定义一个中间的变量
		tmpInt := big.Int{}
		//将hash转为字节流和给定的难度值进行比较
		tmpInt.SetBytes(hash[:])
		//调用big.cmp方法，进行比较
		if tmpInt.Cmp(pow.target) == -1{
			//如果条件成立，说明找到目标的hash值,找到后，退出循环，并返计算出的hash值和随机数
			fmt.Printf("挖矿成功，hash : %x, nonce：%d\n",hash[:], nonce)
			//退出循环
			//break
			return hash[:], nonce
		}else {
			//挖矿不成功，随机数加1，继续进行挖矿
			nonce++
		}
	}
	//返回hash值和随机数
}
//为工作量证明绑定校验函数
func (pow *ProofOfWork)IsValid()bool  {
	//获取拼装好的数据，计算hash值
	nonce := pow.block.Nonce
	blockInfo := pow.PrepareData(nonce)
	//对拼装好的数据进行hash运算
	hash := sha256.Sum256(blockInfo)
	fmt.Printf("校验数据：hash:%x, nonce:%v\n",hash[:],nonce)
	//将计算好的hash值与目标难度值进行比较
	tmpInt := big.Int{}
	tmpInt.SetBytes(hash[:])
	if tmpInt.Cmp(pow.target) == -1{
		//如果比目标值小，返回true
		return true
	}else {
		return false
	}

}

