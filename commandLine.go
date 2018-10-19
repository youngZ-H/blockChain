package main

import (
	"fmt"

	"time"
)

//添加区块
/*func (cli *CLI) AddBlock(data string) {
	//cli.bc.AddBlock(txs)
	//cli.bc.AddBlock(data)
	fmt.Printf("添加区块成功！\n")
}*/
//反向打印区块
func (cli *CLI) PrintBlockChainReverse() {
	bc := cli.bc
	//创建迭代器
	it := bc.NewIterator()

	//调用迭代器，返回我们的每一个区块数据
	for {
		//返回区块，左移
		block := it.Next()

		fmt.Printf("===========================\n\n")
		fmt.Printf("版本号: %d\n", block.Version)
		fmt.Printf("前区块哈希值: %x\n", block.PreHash)
		fmt.Printf("梅克尔根: %x\n", block.MerKleRoot)
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("时间戳: %s\n", timeFormat)
		fmt.Printf("难度值(随便写的）: %d\n", block.Difficulty)
		fmt.Printf("随机数 : %d\n", block.Nonce)
		fmt.Printf("当前区块哈希值: %x\n", block.Hash)
		fmt.Printf("区块数据 :%s\n", block.Transaction[0].TXInputs[0].Sig)

		if len(block.PreHash) == 0 {
			fmt.Printf("区块链遍历结束！\n")
			break
		}
	}
}
//正向打印区块
func (cli *CLI) PrinBlockChain() {
	cli.bc.PrintBlockChain()
	fmt.Printf("打印区块链完成\n")
}
//获取指定地址的余额
func (cli *CLI)GetBalance(address string)  {
	utxos := cli.bc.FindUTXOs(address)
	total := 0.0
	for _, utxo := range utxos{
		total += utxo.Value
	}
	fmt.Printf("\"%s\"的余额为：%f\n",address,total)

}
//进行交易的实现
func (cli *CLI)Send(from, to string, amount float64, miner, data string)  {
	//创建挖矿交易
	coinbase := NewCoinbaseTx(miner,data)
	//创建普通交易
	tx := NewTransaction(from, to,amount,cli.bc)
	if tx == nil{
		//无效的交易，直接返回，结束函数
		return
	}
	//将交易信息记账到区块链中
	cli.bc.AddBlock([]*Transaction{coinbase,tx})
	fmt.Println("转账成功")
}
func (cli *CLI)NewWallet()  {
	//fmt.Println("11111111111111111111111111111")
	ws := NewWallets()
	address := ws.CreateWallet()
	//fmt.Printf("私钥%v\n",wallet.PrivateKey)
	//fmt.Printf("公钥%v\n",wallet.PublicKey)
	fmt.Printf("地址%s\n",address)

}
func (cli *CLI)ListAllAddress()  {
	ws := NewWallets()
	addresses := ws.ListAllAddress()
	//遍历所有的地址，并将其打印出来
	for _,adrs := range addresses{
		fmt.Printf("地址：%s\n",adrs)
	}
}

//正向打印区块
/*func (cli *CLI) PrintBlockChain() {
	bc := cli.bc
	db := bc.db
	blockHeight := 0
	//读取数据库
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		bucket.ForEach(func(k, v []byte) error {
			//判断是否读取到最后一个区块
			if bytes.Equal(k,[]byte("LastBlockHash")){
				fmt.Println("遍历区块链结束")
				return nil
			}
			block := Deserialize(v)
			fmt.Printf("===========================\n\n")
			fmt.Printf("当前区块高度%d\n",blockHeight)
			blockHeight++
			fmt.Printf("版本号: %d\n", block.Version)
			fmt.Printf("前区块哈希值: %x\n", block.PreHash)
			fmt.Printf("梅克尔根: %x\n", block.MerKleRoot)
			fmt.Printf("时间戳: %d\n", block.TimeStamp)
			fmt.Printf("难度值(随便写的）: %d\n", block.Difficulty)
			fmt.Printf("随机数 : %d\n", block.Nonce)
			fmt.Printf("当前区块哈希值: %x\n", block.Hash)
			fmt.Printf("区块数据 :%s\n", block.Data)

			return nil
		})
		return nil
	})
}*/

