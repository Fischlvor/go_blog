package output

// AdvertisementInfo 广告信息。
type AdvertisementInfo struct {
	List  []AdvertisementItem `json:"list"`
	Total int64               `json:"total"`
}

// AdvertisementItem 广告项。
type AdvertisementItem struct {
	ID       int64  `json:"id"`
	AdName   string `json:"ad_name"`
	AdImage  string `json:"ad_image"`
	AdLink   string `json:"ad_link"`
	AdType   int    `json:"ad_type"`
	IsActive bool   `json:"is_active"`
}
