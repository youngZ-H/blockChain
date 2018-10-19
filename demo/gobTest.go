package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

//定义一个结构体
type Person struct {
	Name string
	Age int
}
//定义一个函数进行编码和解码
func main(){
	//声明一个对象
	var tom Person
	tom.Name = "Tom"
	tom.Age = 20
	//定义一个buffer用来存储编码的字节流
	var buffer bytes.Buffer
	//定义一个编码器
	encoder := gob.NewEncoder(&buffer)
	//进行编码
	err := encoder.Encode(&tom)
	if err != nil{
		panic("编码发生错误")
	}
	fmt.Printf("编码后的结果为%v\n",buffer.Bytes())
	//定义一个解码器进行解码
	decoder := gob.NewDecoder(bytes.NewReader(buffer.Bytes()))
	//进行解码
	//定义一个对象，用来存储解码后的结果
	var Tom Person
	err = decoder.Decode(&Tom)
	if err != nil{
		panic("解码时发生错误")
	}
	fmt.Printf("解码后的结果为：%v\n", Tom)


}

