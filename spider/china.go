package spider

import (
	"os"
	"regexp"

	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "china" {
		spiderManager.list = append(spiderManager.list, chinaSpider)
	}
}

func chinaSpider() []NewsItem {
	newsItems := make([]NewsItem, 0)
	urls := []string{
		"http://news.china.com.cn/node_7247300.htm",
	}

	for _, url := range urls {
		status, resp, err := fasthttp.Get(nil, url)
		if err != nil || status != fasthttp.StatusOK {
			continue
		}
		r, _ := regexp.Compile("[\n\r]")
		text := string(r.ReplaceAll(resp, []byte{}))
		reg, _ := regexp.Compile("<li><a\\s*href=\"(.*?)\"\\s*target=\"_blank\">(.*?)</a><span>(.*?)</span>")
		res := reg.FindAllStringSubmatch(text, -1)
		newsItems = make([]NewsItem, 0, len(res))
		for _, matchedItem := range res {
			newsItems = append(newsItems, NewsItem{
				Title:  matchedItem[2],
				Link:   matchedItem[1],
				Origin: "中国网",
			})
		}
	}

	return newsItems
}
