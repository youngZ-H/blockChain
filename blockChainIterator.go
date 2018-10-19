package main

import (
	"blockExercise/bolt"
	"log"
)

//定义区块迭代器
type BlockChainIterator struct {
	db *bolt.DB
	currentHashPointer []byte
}
//为区块链绑定迭代器的方法
func (bc *BlockChain)NewIterator()*BlockChainIterator  {

	return &BlockChainIterator{
		db:bc.db,
		currentHashPointer:bc.tail,
	}

}
//为迭代器制定指针前移的方法
func (it *BlockChainIterator)Next()Block  {
	var block Block
	//读取数据库中的数据
	it.db.View(func(tx *bolt.Tx) error {
		//查找抽屉
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil{
			log.Panic("不应该没有数据，（迭代器）")
		}

		//读取hash值
		hash := bucket.Get(it.currentHashPointer)
		//进行解码
		block = Deserialize(hash)
		//游标hash左移
		it.currentHashPointer = block.PreHash
		return nil
	})
	return block
}
