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
	Title string            `json:"title"`
	Links []spider.NewsItem `json:"links"`
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
			if simHash.IsEqual(distanceItem.distance, distance, 3) {
				isEqual = true
				isExists := false
				for _, link := range distanceItem.item.Links {
					if link.Origin == item.Origin && link.Title == item.Title {
						isExists = true
						break
					}
				}
				if isExists {
					break
				}
				if len(item.Title) > len(distanceItem.item.Title) {
					distanceItem.item.Title = item.Title
				}
				distanceItem.item.Links = append(distanceItem.item.Links, item)
				break
			}
		}
		if !isEqual {
			resultItem := ResultItem{
				Title: item.Title,
				Links: []spider.NewsItem{item},
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
		jLinksLen := len(result.Items[j].Links)
		iLinksLen := len(result.Items[i].Links)
		if jLinksLen != iLinksLen {
			return jLinksLen < iLinksLen
		}
		return len(result.Items[j].Title) < len(result.Items[i].Title)
	})

	size := len(result.Items)

	if size > 100 {
		size = 100
	}

	result.Items = result.Items[0:size]

	now := time.Now().Format("2006-01-02 15:04:05")

	jsonStr, _ := json.Marshal(result.Items)

	json, _ := os.Create("./result/news.json")
	defer json.Close()
	json.Write(jsonStr)

	jsonp, _ := os.Create("./result/news.jsonp")
	defer jsonp.Close()
	jsonp.Write([]byte("/* */window.newsJsonp && window.newsJsonp(\"" + now + "\", " + string(jsonStr) + ");"))

	mdStr := "## News Update\n---\n" + now + "\n---\n"

	for index, item := range result.Items {
		if len(item.Links) > 1 {
			mdStr += strconv.Itoa(index+1) + ". " + item.Title + " (" + strconv.Itoa(len(item.Links)) + ")\n"
			for _, link := range item.Links {
				mdStr += "    +  <a target=\"_blank\" href=\"" + link.Link + "\">[" + link.Origin + "] " + link.Title + "</a>\n"
			}
			mdStr += "\n"
		} else {
			mdStr += strconv.Itoa(index+1) + ". <a target=\"_blank\" href=\"" + item.Links[0].Link + "\">[" + item.Links[0].Origin + "] " + item.Title + "</a>\n"
		}

	}

	md, _ := os.Create("readme.md")
	defer md.Close()
	md.Write([]byte(mdStr))
}
