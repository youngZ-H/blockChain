package main

import (
	"blockExercise/bolt"
	"log"
	"bytes"
	"fmt"
	"time"
)

//引入区块连
type BlockChain struct {
	//blocks []*Block
	db *bolt.DB	//数据库用来存储区块的数据
	tail []byte //用来存储最后一个区块的hash
}

const blockChainDb  = "blockChain.db"
const blockBucket   = "blockBucket"

//创建区块连，利用数据库进行区块链的改写
func CreateBlockChain(address string ) *BlockChain {
/*	genesisBlock := GenesisBlock()
	return &BlockChain{
		blocks: []*Block{genesisBlock},
	}*/
	//打开数据库
	var LastHash []byte
	db, err := bolt.Open(blockChainDb,0600, nil)
	if err != nil{
		log.Panic("打开数据库时发生错误")
	}

	//对数据库进行操作
	db.Update(func(tx *bolt.Tx) error {
		//查找blockBucket抽屉，没有进行抽屉的创建
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil{
			//log.Panic("===========================所要查找的抽屉不存在")//当抽屉不存在时，进行创建
		bucket, err  = tx.CreateBucket([]byte(blockBucket))
		if err != nil{
			log.Panic("创建blockBucket发生错误")
		}
		//创建创世块
		genesisBlock := GenesisBlock(address)
		//向bucket中写入数据,把创世块的信息写入区块链中，hash值作为key，区块数据作为value
		bucket.Put(genesisBlock.Hash,genesisBlock.Serialize())
		bucket.Put([]byte("LastBlockHash"),genesisBlock.Hash)
		LastHash = genesisBlock.Hash
		}else {
			LastHash = bucket.Get([]byte("LastBlockHash"))
		}

		return nil
	})
	return &BlockChain{db,LastHash}
}

//生成创世块
func GenesisBlock(address string ) *Block {
	coinbaseTx := NewCoinbaseTx(address,"Go一期创世块，老牛逼啦")
	return CreateBlock([]*Transaction{coinbaseTx}, []byte{})
}

//为区块连结构体榜定添加区块的方法
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	//添加区块需要区块的数据和前区块的hash
	//获取前去块的hash
	/*lastBlock := bc.blocks[len(bc.blocks)-1]
	//获取前一区块的hash
	preHash := lastBlock.Hash
	//创建新的区块
	block := CreateBlock(data, preHash)
	//将新创建的区块添加的区块连中
	bc.blocks = append(bc.blocks, block)*/
	//利用数据库存储改写添加区块的方法
	// 获取最后一个区块的hash
	lastBlockHash := bc.tail
	//创建一个新的区块
	newBlock := CreateBlock(txs,lastBlockHash)
	//将数据存储到数据库
	bc.db.Update(func(tx *bolt.Tx) error {
		//查找抽屉是否存在
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil{
			log.Panic("数据不应该不存在，请检查错误")
		}else {
			//如果抽屉存在，向数据库中写入数据
			bucket.Put(newBlock.Hash,newBlock.Serialize())
		//更新最后一个区块的hash
			bucket.Put([]byte("LastBlockHash"),newBlock.Hash)
		}
		bc.tail = newBlock.Hash
		return nil
	})
}
func (bc *BlockChain) PrintBlockChain() {

	blockHeight := 0
	bc.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("blockBucket"))

		//从第一个key-> value 进行遍历，到最后一个固定的key时直接返回
		b.ForEach(func(k, v []byte) error {
			if bytes.Equal(k, []byte("LastBlockHash")) {
				return nil
			}

			block := Deserialize(v)
			//fmt.Printf("key=%x, value=%s\n", k, v)
			fmt.Printf("=============== 区块高度: %d ==============\n", blockHeight)
			blockHeight++
			fmt.Printf("版本号: %d\n", block.Version)
			fmt.Printf("前区块哈希值: %x\n", block.PreHash)
			fmt.Printf("梅克尔根: %x\n", block.MerKleRoot)
			timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
			fmt.Printf("时间戳: %s\n", timeFormat)
			fmt.Printf("难度值(随便写的）: %d\n", block.Difficulty)
			fmt.Printf("随机数 : %d\n", block.Nonce)
			fmt.Printf("当前区块哈希值: %x\n", block.Hash)
			fmt.Printf("区块数据 :%s\n", block.Transaction[0].TXInputs[0].Sig)
			return nil
		})
		return nil
	})
}
//获取指定地址的所有的UTXO
func (bc *BlockChain)FindUTXOs(address string )[]TXOutput  {
	var UTXO = []TXOutput{}

	txs := bc.FindNeedTransactions(address)
	for _,tx := range txs{
		//遍历每笔交易的outputs
		for _,output := range tx.TXOutputs{
			//判断每笔的output的签名地址是否时from的
			if output.PubKeyHash == address {
			UTXO = append(UTXO,output)

			}
		}
	}
	//定义一个map spendOutputs,用来存储消耗过的UTXO,k值是交易的id，value时输入的索引值
	/*spendOutputs := make(map[string][]int64)
	//遍历区块
	it := bc.NewIterator()
	for {
		block := it.Next()
		//遍历交易
		for _, tx := range block.Transaction{
			//遍历outputs,将所有的UTXO存储到切片中（在存储之前，要先进行过滤，看是否消耗过）
			OUTPUT:
			for i, output := range tx.TXOutputs{
				//fmt.Println("当前的索引值",i)
				//在进行UTXO存储之前，先进行过滤，将已经消耗过的UTXO去除
				//遍历spendOutputs,进行索引值的比对
				//先进行判断spendOutputs是否为空，如果为空，就不需要比对啦，因为不存在消耗的output
				//fmt.Println("222222222222222222222222")
				//fmt.Println("***********",spendOutputs[string(tx.TXID)])
				if spendOutputs[string(tx.TXID)] != nil{
					//fmt.Println("333333333333333333333333333")
					for _,index := range spendOutputs[string(tx.TXID)]{
					//判断当前的索引值是否存在消耗过的output中
					//fmt.Println("+++++++++++",index)
					if int64(i) == index{
						continue OUTPUT
					}
				}
			}
				//将所有的UTXO存储到UTXO切片中,进行判断，判断交易地址是否和我们要查询的地址相同
				if output.PubKeyHash == address{
					UTXO = append(UTXO,output)
					//fmt.Println("++++++++++++++++++",UTXO)
				}
			}
			//如果时挖矿交易，不做遍历
			if !tx.IsCoinbaseTx(){
				//fmt.Println("111111111111111111111111111")
					//遍历input，将自己所有消耗的output找出来，并存储在消耗spendOutputs切片中
					for _, input := range tx.TXInputs{
					//判断签名的地址是否要和我们查找的地址相同，相同则代表消耗过，存储到消耗的map中，共过滤进行使用
					if input.Sig == address{
					spendOutputs[string(input.TXid)] = append(spendOutputs[string(input.TXid)],input.Index)
					}

						//fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&",spendOutputs[string(tx.TXID)])
				}
			}
		}

		if len(block.PreHash) == 0{
			//fmt.Println("区块遍历结束")
			break
		}

	}*/
	 return UTXO

}
//为blockChain绑定查找最优的UTXO的方法实现
func(bc *BlockChain)FindNeedUTXOs(from string,amount float64)(map[string][]int64,float64){
	var utxos = map[string][]int64{}
	var cal float64
	//调用返回去未消费utxo的交易结构的方法
	txs := bc.FindNeedTransactions(from)
	//遍历交易切片，找到最合适的utxo
	for _,tx := range txs{
		//遍历每笔交易的outputs
		for i,output := range tx.TXOutputs{
			//判断每笔的output的签名地址是否时from的
			if output.PubKeyHash == from {


				//如果是的话，再判断目前找到的金额cal是否满足转账的要求amount
				if cal < amount{
					//将所有满足要求的utxo取出，存储起来
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)],int64(i) )
					//只有不满足时在将未消费的output的value取出
					cal += output.Value
					//在进行判断，看是否满足消费需求
					if cal >= amount{
						//如果满足需求就直接返回搜需要的值
						fmt.Printf("找到合适金额为:%f\n",cal)
						return utxos,cal
						}
					}else {
					fmt.Printf("对不起，您的余额不足，您的余额为%f, 所消费的金额为:%f\n",cal, amount)
				}
			}
		}
	}
	/*var utxos = map[string][]int64{}
	var cal float64
	//定义一个map spendOutputs,用来存储消耗过的UTXO,k值是交易的id，value时输入的索引值
	spendOutputs := make(map[string][]int64)
	//遍历区块
	it := bc.NewIterator()
	for {
		block := it.Next()
		//遍历交易
		for _, tx := range block.Transaction{
			//遍历outputs,将所有的UTXO存储到切片中（在存储之前，要先进行过滤，看是否消耗过）
			OUTPUT:
			for i, output := range tx.TXOutputs{
				//fmt.Println("当前的索引值",i)
				//在进行UTXO存储之前，先进行过滤，将已经消耗过的UTXO去除
				//遍历spendOutputs,进行索引值的比对
				//先进行判断spendOutputs是否为空，如果为空，就不需要比对啦，因为不存在消耗的output
				if spendOutputs[string(tx.TXID)] != nil{
					for _,index := range spendOutputs[string(tx.TXID)]{
					//判断当前的索引值是否存在消耗过的output中
					if int64(i) == index{
						continue OUTPUT
					}
				}
			}
				//将所有的UTXO存储到UTXO切片中,进行判断，判断交易地址是否和我们要查询的地址相同
				//找到需要的utxo，根据传进来的地址进行条件的判断
				if output.PubKeyHash == from {
					//判断要消费的金额和找到的金额的大小的对比，判断是否找到，
					if cal < amount{
						//将找到的未消费UTXO存储到map中
						utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)],int64(i))
						//将找到的未消费的金额加到cal中
						cal += output.Value
						//将加完的金额再判断一下，是否满足条件，如果满足条件，就跳出循环，返回utxos,cal
						if cal > amount{
							fmt.Printf("找到合适的消费金额%v\n",cal)
							return utxos,cal
						}
					}else {
						fmt.Printf("对不起，您的余额不足，您的余额为%f, 所消费的金额为:%f\n",cal, amount)
					}
				}

			}
			//如果时挖矿交易，不做遍历
			if !tx.IsCoinbaseTx(){
					//遍历input，将自己所有消耗的output找出来，并存储在消耗spendOutputs切片中
					for _, input := range tx.TXInputs{
					//判断签名的地址是否要和我们查找的地址相同，相同则代表消耗过，存储到消耗的map中，共过滤进行使用
					if input.Sig == from {
					spendOutputs[string(tx.TXID)] = append(spendOutputs[string(tx.TXID)],input.Index)
					}

				}
			}
		}

		if len(block.PreHash) == 0{
			//fmt.Println("区块遍历结束")
			break
		}

	}*/
	return utxos,cal

}
//为区块绑定获取未消费的utxo的交易，并返回
func(bc *BlockChain)FindNeedTransactions(address string)[]*Transaction{
	//var utxos = map[string][]int64{}
	var txs  []*Transaction
	//var cal float64
	//定义一个map spendOutputs,用来存储消耗过的UTXO,k值是交易的id，value时输入的索引值
	spendOutputs := make(map[string][]int64)
	//遍历区块
	it := bc.NewIterator()
	for {
		block := it.Next()
		//遍历交易
		for _, tx := range block.Transaction{
			//遍历outputs,将所有的UTXO存储到切片中（在存储之前，要先进行过滤，看是否消耗过）
			OUTPUT:
			for i, output := range tx.TXOutputs{
				//fmt.Println("当前的索引值",i)
				//在进行UTXO存储之前，先进行过滤，将已经消耗过的UTXO去除
				//遍历spendOutputs,进行索引值的比对
				//先进行判断spendOutputs是否为空，如果为空，就不需要比对啦，因为不存在消耗的output
				if spendOutputs[string(tx.TXID)] != nil{
					for _,index := range spendOutputs[string(tx.TXID)]{
					//判断当前的索引值是否存在消耗过的output中
					if int64(i) == index{
						continue OUTPUT
					}
				}
			}
				//将所有的UTXO存储到UTXO切片中,进行判断，判断交易地址是否和我们要查询的地址相同
				//找到需要的utxo，根据传进来的地址进行条件的判断
				if output.PubKeyHash == address {
					//判断要消费的金额和找到的金额的大小的对比，判断是否找到，
					/*if cal < amount{
						//将找到的未消费UTXO存储到map中
						utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)],int64(i))
						//将找到的未消费的金额加到cal中
						cal += output.Value
						//将加完的金额再判断一下，是否满足条件，如果满足条件，就跳出循环，返回utxos,cal
						if cal > amount{
							fmt.Printf("找到合适的消费金额%v\n",cal)
							return utxos,cal
						}
					}else {
						fmt.Printf("对不起，您的余额不足，您的余额为%f, 所消费的金额为:%f\n",cal, amount)
					}*/
					//通过地址表，找出属于这个地址的utxo的交易，并存储起来，然后返回交易结构
					txs = append(txs, tx)
				}
			}
			//如果时挖矿交易，不做遍历
			if !tx.IsCoinbaseTx(){
					//遍历input，将自己所有消耗的output找出来，并存储在消耗spendOutputs切片中
					for _, input := range tx.TXInputs{
					//判断签名的地址是否要和我们查找的地址相同，相同则代表消耗过，存储到消耗的map中，共过滤进行使用
					if input.Sig == address{
					spendOutputs[string(input.TXid)] = append(spendOutputs[string(input.TXid)],input.Index)
					}

				}
			}
		}

		if len(block.PreHash) == 0{
			//fmt.Println("区块遍历结束")
			break
		}

	}
	return txs

}