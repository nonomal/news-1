package utils

import (
	"math"
	"strconv"
	"strings"
	"time"
)

var daySec = 86400

func GetTimezone() *time.Location {
	locate, lerr := time.LoadLocation("Asia/Beijing")
	if lerr != nil {
		locate = time.FixedZone("CST", 8*3600)
	}
	return locate
}

func FormatTimeHMToUnix(hm string) int64 {
	now := time.Now()
	nowUnix := now.Unix()
	locate := GetTimezone()
	pubTime, _ := time.ParseInLocation("2006-01-02 15:04", now.Format("2006-01-02")+" "+hm, locate)
	pubTimeUnix := pubTime.Unix()
	if pubTimeUnix < nowUnix {
		return pubTimeUnix
	}

	diff := math.Ceil(float64(pubTimeUnix-nowUnix) / float64(daySec))
	return pubTimeUnix - int64(diff)*int64(daySec)
}

func FormatTimeYMDToUnix(tStr string) int64 {
	locate := GetTimezone()
	pubTime, _ := time.ParseInLocation("2006-01-02", tStr, locate)
	return pubTime.Unix()
}

func FormatTimemdToUnix(tStr string) int64 {
	locate := GetTimezone()
	now := time.Now()
	pubTime, _ := time.ParseInLocation("2006-1-2 15:04", strconv.Itoa(now.Year())+"-"+tStr, locate)
	return pubTime.Unix()
}

func FormatTimeYMDHMSToUnix(tStr string) int64 {
	locate := GetTimezone()
	pubTime, _ := time.ParseInLocation("2006-01-02 15:04:05", tStr, locate)
	return pubTime.Unix()
}

func FormatTime(unix int64, layout string) string {
	now := time.Unix(unix, 0)
	locate := GetTimezone()
	return now.In(locate).Format(layout)
}

func FormatNow(layout string) string {
	return FormatTime(time.Now().Unix(), layout)
}

func FormatTimeAgo(ago string) int64 {
	now := time.Now().Unix()
	if strings.HasSuffix(ago, "分钟前") {
		minNum, err := strconv.Atoi(ago[0 : len(ago)-len("分钟前")])
		if err != nil {
			return 0
		}
		now -= int64(minNum * 60)
	} else if strings.HasSuffix(ago, "小时前") {
		minNum, err := strconv.Atoi(ago[0 : len(ago)-len("小时前")])
		if err != nil {
			return 0
		}
		now -= int64(minNum * 60 * 60)
	}
	if now > 0 {
		return now
	}
	return 0
}
