package main

import (
	"time"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

//定义区块结构体
type Block struct {
	Version    uint64 //版本号
	PreHash    []byte //前区块hash
	MerKleRoot []byte //梅克尔根
	TimeStamp  uint64 //时间戳
	Difficulty uint64 //难度值
	Nonce      uint64 //随机数
	Hash       []byte //当前区块hash
	//Data       []byte //区块数据
	//真实的交易数据
	Transaction []*Transaction
}

//为区块绑定方法，创建新的区块
func CreateBlock(txs []*Transaction, preBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		PreHash:    preBlockHash,
		MerKleRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 100,
		Nonce:      100,
		Hash:       []byte{},
		//Data:       []byte(data),
		Transaction: txs,
	}
	//block.SetHash()
	block.MerKleRoot = block.CreateMeklerRoot()
	//进行工作量证明的调用
	pow := NewProofWork(&block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return &block

}
//uint64 向 【】byte转换函数的实现
func Uint64ToByte(num uint64)[]byte  {
	var buffer bytes.Buffer
	err := binary.Write(&buffer,binary.BigEndian,num)
	if err != nil{
		panic(err)
	}
	return buffer.Bytes()
}
//生成hash，为区块绑定生成hash的方法
/*func (block *Block) SetHash() {
	//拼装区块数据
	//var blockInfo []byte
	//blockInfo = append(blockInfo,uint64ToByte(block.Version)...)
	//blockInfo = append(blockInfo,block.PreHash...)
	//blockInfo = append(blockInfo,block.MerKleRoot...)
	//blockInfo = append(blockInfo,uint64ToByte(block.TimeStamp)...)
	//blockInfo = append(blockInfo,uint64ToByte(block.Difficulty)...)
	//blockInfo = append(blockInfo,uint64ToByte(block.Nonce)...)
	//blockInfo = append(blockInfo,block.Hash...)
	//blockInfo = append(blockInfo,block.Data...)
	tmp := [][]byte{
		Uint64ToByte(block.Version),
		block.PreHash,
		block.MerKleRoot,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(block.Nonce),
		block.Data,
	}
	blockInfo := bytes.Join(tmp,[]byte{})

	//生成hash
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]

}
//为block绑定转为字节流的方法*/
func (block *Block)Serialize()[]byte  {
	//定义一个buffer用来存储字节流
	var buffer bytes.Buffer
	//定义一个编码器
	encoder := gob.NewEncoder(&buffer)
	//进行编码
	err := encoder.Encode(&block)
	if err != nil{
		log.Panic("进行编码时发生错误")
	}
	return buffer.Bytes()
}
//定义一个解码器的函数
func Deserialize(data []byte) Block  {
	//定义一个解码器
	decoder := gob.NewDecoder(bytes.NewReader(data))
	//进行解码
	var block Block
	err := decoder.Decode(&block)
	if err != nil{
		log.Panic("解码时发生错误")
	}
	return block
}
//为block绑定生成梅克尔根的方法
func (block *Block)CreateMeklerRoot()[]byte  {
	//梅克尔根时将区块体的交易信息取hash运算，即将交易结构的TXID再次进行hash运算
	//将一个区块中的所有交易信心的hash值拼接起来再次进行hash运算
	var transactionInfo []byte
	txs := block.Transaction
	for _, tx := range txs{
		transactionInfo = append(transactionInfo, tx.TXID...)
	}
	//将所有的拼接信息再次进行hash运算
	hash := sha256.Sum256(transactionInfo)
	return hash[:]

}