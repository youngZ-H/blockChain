package main

import (
)
func main() {
	//向区块连中添加新的区块
	bc := CreateBlockChain("张三")
	cli := CLI{bc}
	cli.Run()
	//bc.AddBlock("确认过眼神")
	//bc.AddBlock("你是对的人")
	////创建迭代器
	//it := bc.NewIterator()
	//for {
	//	block := it.Next()
	//	fmt.Printf("PreHash：%x\n", block.PreHash)
	//	fmt.Printf("hash :%x\n", block.Hash)
	//	fmt.Printf("%s\n", block.Data)
	//	//在主函数中调用校验函数
	//	pow := NewProofWork(&block)
	//	fmt.Printf("IsValid %v\n",pow.IsValid())
	//	if len(block.PreHash) == 0{
	//		fmt.Println("区块遍历结束")
	//		break
	//	}
	//
	//}
	////打印前区块的hash值
	//for i, block := range bc.blocks {
	//	fmt.Printf("区块高度%d\n", i)
	//	fmt.Printf("PreHash：%x\n", block.PreHash)
	//	fmt.Printf("hash :%x\n", block.Hash)
	//	fmt.Printf("%s\n", block.Data)
	//	//在主函数中调用校验函数
	//	pow := NewProofWork(block)
	//	fmt.Printf("IsValid %v\n",pow.IsValid())
	//}

}
