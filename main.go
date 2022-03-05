package main

import (
	"fmt"

	"github.com/echosoar/news/simHash"
	"github.com/echosoar/news/spider"
)

func main() {
	list := spider.Get()

	fmt.Println(len(list))

	a := simHash.Calc("热烈祝贺十三届全国人大五次会议开幕")
	b := simHash.Calc("十三届全国人大五次会议今日上午9时开幕")

	ad := simHash.Distance(a)
	bd := simHash.Distance(b)
	fmt.Println(simHash.IsEqual(ad, bd, 10))
}
