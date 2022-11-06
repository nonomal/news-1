package spider

import (
	"strings"
	"sync"

	"github.com/echosoar/news/utils"
)

var filterList [][]string = [][]string{
	{"新增", "确诊"},
	{"新增", "无症状"},
	{"新增", "感染"},
	// 非新闻
	{"学习", "精神"},
	{"坚持", "思想"},
	{"伟大复兴"},
	{"牢牢把握"},
	{"扎根"},
	{"基层", "乡村"},
	{"在基层"},
	{"奋斗"},
	{"创造", "明天"},
	{"绿水青山"},
	{"伟大", "时代"},
	{"为人民"},
	{"讲述"},
	{"微记录"},
	{"前程"},
	{"一起", "成长"},
	{"沿着", "足迹"},
	{"见证", "变化"},
	{"共同", "谱写"},
	{"出席", "致辞"},
	{"红火"},

	{"面对面"},
	{"摇篮"},
}

type NewsItem struct {
	Title  string
	Link   string
	Origin string
	Time   int64
}

type SpiderManager struct {
	list []Spider
}

var spiderManager = &SpiderManager{
	list: make([]Spider, 0),
}

type Spider = func() []NewsItem

func Get() []NewsItem {
	spiderCount := len(spiderManager.list)
	var wg sync.WaitGroup
	wg.Add(spiderCount)
	channel := make(chan *[]NewsItem, spiderCount)
	for i := 0; i < spiderCount; i++ {
		go func(ch chan *[]NewsItem, index int, wg *sync.WaitGroup) {
			list := spiderManager.list[index]()
			ch <- &list
			wg.Done()
		}(channel, i, &wg)
	}
	go func() {
		wg.Wait()
		close(channel)
	}()
	newsItems := make([]NewsItem, 0)
	for ch := range channel {
		for _, item := range *ch {
			if !isNeedFilter(item.Title) {
				newsItems = append(newsItems, item)
			} else {
			}

		}

	}
	return newsItems
}

func isNeedFilter(title string) bool {
	for _, filterWords := range filterList {
		isNeedFilter := true
		for _, word := range filterWords {
			if !strings.Contains(title, word) {
				isNeedFilter = false
				break
			}
		}
		if isNeedFilter {
			return true
		}
	}
	return false
}

func ItemToHtml(item *NewsItem) string {
	pubTime := ""
	if item.Time > 0 {
		pubTime = " - " + utils.FormatTime(item.Time, "01/02 15:04")
	}
	return "<a target=\"_blank\" href=\"" + item.Link + "\">" + item.Title + "</a> [" + item.Origin + pubTime + "]"
}
