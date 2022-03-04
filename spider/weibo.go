package spider

func init() {
	spiderManager.list = append(spiderManager.list, weiboSpider)
}

func weiboSpider() []NewsItem {
	newsItems := make([]NewsItem, 0)
	newsItems = append(newsItems, NewsItem{
		title: "测试",
		link:  "test",
	})
	return newsItems
}
