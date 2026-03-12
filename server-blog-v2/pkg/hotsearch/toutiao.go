package hotsearch

import (
	"io"
	"net/http"
	"strconv"

	"github.com/tidwall/gjson"
)

type Toutiao struct{}

func (*Toutiao) GetHotSearchData(maxNum int) (HotSearchData, error) {
	resp, err := http.Get("https://www.toutiao.com/hot-event/hot-board/?origin=toutiao_pc")
	if err != nil {
		return HotSearchData{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return HotSearchData{}, err
	}

	jsonStr := string(body)

	timeStr := gjson.Get(jsonStr, "impr_id").Str
	var updateTime string
	if len(timeStr) >= 14 {
		updateTime = timeStr[:4] + "-" + timeStr[4:6] + "-" + timeStr[6:8] + " " + timeStr[8:10] + ":" + timeStr[10:12] + ":" + timeStr[12:14]
	}

	var hotList []HotItem

	for i := 0; i < maxNum; i++ {
		if index := gjson.Get(jsonStr, "data."+strconv.Itoa(i)+".ClusterId"); !index.Exists() {
			break
		}
		hotList = append(hotList, HotItem{
			Index:       i + 1,
			Title:       gjson.Get(jsonStr, "data."+strconv.Itoa(i)+".Title").Str,
			Description: "",
			Image:       gjson.Get(jsonStr, "data."+strconv.Itoa(i)+".Image.url").Str,
			Popularity:  gjson.Get(jsonStr, "data."+strconv.Itoa(i)+".HotValue").Str,
			URL:         gjson.Get(jsonStr, "data."+strconv.Itoa(i)+".Url").Str,
		})
	}

	return HotSearchData{Source: "头条热榜", UpdateTime: updateTime, HotList: hotList}, nil
}
