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
	"github.com/echosoar/news/utils"
)

type Result struct {
	Distances []*DitanceItem
	Items     []*ResultItem
}

type ResultItem struct {
	Title    string            `json:"title"`
	Links    []spider.NewsItem `json:"links"`
	Time     int64             `json:"time"`
	Keywords []string          `json:"keywords"`
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
	x := simHash.GetJieba()
	defer x.Free()

	nowTimeStamp := time.Now().Unix()
	var daySec int64 = 24 * 3600
	for _, item := range list {
		if nowTimeStamp-item.Time > daySec {
			continue
		}
		hash, keywords := simHash.Calc(x, item.Title)
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
					distanceItem.item.Keywords = keywords
				}
				if item.Time > distanceItem.item.Time {
					distanceItem.item.Time = item.Time
				}
				distanceItem.item.Links = append(distanceItem.item.Links, item)
				break
			}
		}
		if !isEqual {
			resultItem := ResultItem{
				Title:    item.Title,
				Time:     item.Time,
				Links:    []spider.NewsItem{item},
				Keywords: keywords,
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
		return result.Items[j].Time < result.Items[i].Time
	})

	size := len(result.Items)

	if size > 100 {
		size = 100
	}

	result.Items = result.Items[0:size]

	now := utils.FormatNow("2006-01-02 15:04:05")

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
				mdStr += "    +  " + spider.ItemToHtml(&link) + "\n"
			}
			mdStr += "\n"
		} else {
			mdStr += strconv.Itoa(index+1) + ". " + spider.ItemToHtml(&(item.Links[0])) + "\n"
		}
	}

	md, _ := os.Create("readme.md")
	defer md.Close()
	md.Write([]byte(mdStr))
}
