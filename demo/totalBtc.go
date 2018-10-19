package main

import (
	"fmt"
)

func main()  {
	total := 0.0
	reduceWard:= 50.0
	reduceCount := 0
	blockInterval := 21 //单位为万
	for reduceWard > 0 {
		reduceCount++
		if reduceCount > 31{
			break
		}
		sum := float64(blockInterval) * reduceWard
		total += sum
		reduceWard *= 0.5
	}
	fmt.Printf("产生的比特比总数：%f万, 衰减次数%v",total,reduceCount)
}