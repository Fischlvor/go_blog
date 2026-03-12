package hotsearch

// HotItem 热搜项。
type HotItem struct {
	Index       int    `json:"index"`       // 排名
	Title       string `json:"title"`       // 标题
	Description string `json:"description"` // 描述
	Image       string `json:"image"`       // 图片
	Popularity  string `json:"popularity"`  // 热度
	URL         string `json:"url"`         // 链接
}

// HotSearchData 热搜数据。
type HotSearchData struct {
	Source     string    `json:"source"`      // 数据源
	UpdateTime string    `json:"update_time"` // 更新时间
	HotList    []HotItem `json:"hot_list"`    // 热搜列表
}

// Source 热搜数据源接口。
type Source interface {
	GetHotSearchData(maxNum int) (HotSearchData, error)
}

// NewSource 创建热搜数据源。
func NewSource(sourceStr string) Source {
	switch sourceStr {
	case "baidu":
		return &Baidu{}
	case "kuaishou":
		return &Kuaishou{}
	case "toutiao":
		return &Toutiao{}
	case "zhihu":
		return &Zhihu{}
	default:
		return nil
	}
}
