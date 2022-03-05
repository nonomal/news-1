package spider

import (
	"encoding/json"
	"regexp"

	"github.com/valyala/fasthttp"
)

func init() {
	spiderManager.list = append(spiderManager.list, jiemianSpider)
}

type jiemianJson struct {
	Rst string `json:"rst"`
}

func jiemianSpider() []NewsItem {
	url := "https://a.jiemian.com/index.php?m=lists&a=ajaxNews&cid=4&page=1"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	jiemianJsonStruct := jiemianJson{}
	json.Unmarshal(resp[1:len(resp)-1], &jiemianJsonStruct)

	reg, _ := regexp.Compile("\"item-date\">(.*?)</div><div class=\"item-main\"><p>\n\\s*.*?href=\"(.*?)\".*?\"_blank\">(.*?)</a>.*?\n\\s*(.*?)</p>")
	res := reg.FindAllStringSubmatch(jiemianJsonStruct.Rst, -1)
	newsItems = make([]NewsItem, 0, len(res))
	for _, matchedItem := range res {
		newsItems = append(newsItems, NewsItem{
			title: matchedItem[3],
			link:  matchedItem[2],
		})
	}
	return newsItems
}
