package main

import (
	"fmt"

	"github.com/echosoar/news/spider"
)

func main() {
	list := spider.Get()
	fmt.Println("test", list)
}
