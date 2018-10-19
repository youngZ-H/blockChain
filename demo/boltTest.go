package main

import (
	"blockExercise/bolt"
	"log"
	"fmt"
)

func main()  {
	db, err := bolt.Open("test.db",0600,nil)
	if err != nil{
		panic("打开数据库错误")
		return
	}
	defer db.Close()
	//操作数据库，向数据库中写数据
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil{
			//如果找不到抽屉，就新建一个
			bucket, err = tx.CreateBucket([]byte("b1"))
			if err !=nil{
				panic("新建bucket发生错误")
			}
		}
		//向数据库中写入数据
		bucket.Put([]byte("11111"),[]byte("hello"))
		bucket.Put([]byte("22222"),[]byte("world"))
		return nil
	})

	//从数据库中读取数据
	db.View(func(tx *bolt.Tx) error {
		//找到抽屉，进行数据的读取
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil{
			log.Panic("发生错误，数据不应该为空")
		}
		//进行数据的读取
		v1 := bucket.Get([]byte("11111"))
		v2 := bucket.Get([]byte("22222"))
		fmt.Printf("%s\n",v1)
		fmt.Printf("%s\n",v2)
		return nil
	})
}
