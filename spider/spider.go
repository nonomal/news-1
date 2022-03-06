package spider

import (
	"sync"

	"github.com/echosoar/news/utils"
)

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
		newsItems = append(newsItems, *ch...)
	}
	return newsItems
}

func ItemToHtml(item *NewsItem) string {
	pubTime := ""
	if item.Time > 0 {
		pubTime = " - " + utils.FormatTime(item.Time, "01/02 15:04")
	}
	return "<a target=\"_blank\" href=\"" + item.Link + "\">" + item.Title + "</a> [" + item.Origin + pubTime + "]"
}
