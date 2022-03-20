package simHash

import (
	"path"
	"runtime"

	"github.com/yanyiwu/gojieba"
)

// 词性对照 https://gist.github.com/hscspring/c985355e0814f01437eaf8fd55fd7998
func GetJieba() *gojieba.Jieba {
	_, filePath, _, _ := runtime.Caller(0)
	newDict := path.Join(path.Dir(filePath), "dict.utf8")
	x := gojieba.NewJieba(
		gojieba.DICT_PATH,
		gojieba.HMM_PATH,
		newDict,
		gojieba.IDF_PATH,
		gojieba.STOP_WORDS_PATH,
	)
	return x
}
