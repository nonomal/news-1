package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/echosoar/news/simHash"
	"github.com/echosoar/news/spider"
	"github.com/yanyiwu/gojieba"
)

type Result struct {
	Distances []*DitanceItem
	Items     []*ResultItem
}

type ResultItem struct {
	Title string   `json:"title"`
	Links []string `json:"links"`
}

type DitanceItem struct {
	distance uint64
	item     *ResultItem
}

func main() {
	list := spider.Get()

	result := Result{
		Distances: make([]*DitanceItem, 0),
		Items:     make([]*ResultItem, 0),
	}
	fmt.Println(len(list))
	x := gojieba.NewJieba()
	defer x.Free()
	for _, item := range list {
		hash := simHash.Calc(x, item.Title)
		distance := simHash.Distance(hash)

		isEqual := false
		for _, distanceItem := range result.Distances {
			if simHash.IsEqual(distanceItem.distance, distance, 5) {
				isEqual = true
				distanceItem.item.Links = append(distanceItem.item.Links, item.Link)
			}
		}
		if !isEqual {
			resultItem := ResultItem{
				Title: item.Title,
				Links: []string{item.Link},
			}
			distanceItem := DitanceItem{
				distance: distance,
				item:     &resultItem,
			}
			result.Items = append(result.Items, &resultItem)
			result.Distances = append(result.Distances, &distanceItem)
		}
	}
	sort.Slice(result.Items, func(i, j int) bool {
		return len(result.Items[j].Links) < len(result.Items[i].Links)
	})

	result.Items = result.Items[0:50]

	nowTime := time.Now()
	timeStr := strconv.Itoa(nowTime.Year()) + "-" + strconv.Itoa(int(nowTime.Month())) + "-" + strconv.Itoa(nowTime.Day())
	jsonStr, _ := json.Marshal(result.Items)

	json, _ := os.Create("./result/news.json")
	defer json.Close()
	json.Write(jsonStr)

	jsonp, _ := os.Create("./result/news.jsonp")
	defer jsonp.Close()
	jsonp.Write([]byte("/* */window.newsJsonp && window.newsJsonp(" + timeStr + ", " + string(jsonStr) + ");"))

	mdStr := "## News Update\n---\n" + timeStr + "\n---\n"

	for index, item := range result.Items {
		mdStr += strconv.Itoa(index+1) + ". <a target=\"_blank\" href=\"" + item.Links[0] + "\">" + item.Title + " (" + strconv.Itoa(len(item.Links)) + ")</a>\n"
	}

	md, _ := os.Create("readme.md")
	defer md.Close()
	md.Write([]byte(mdStr))
}
