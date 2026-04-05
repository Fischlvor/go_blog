package output

// WebsiteInfo 网站信息。
type WebsiteInfo struct {
	Avatar               string `json:"avatar"`
	Logo                 string `json:"logo"`
	FullLogo             string `json:"full_logo"`
	Title                string `json:"title"`
	Slogan               string `json:"slogan"`
	SloganEn             string `json:"slogan_en"`
	Description          string `json:"description"`
	Version              string `json:"version"`
	CreatedAt            string `json:"created_at"`
	ICPFiling            string `json:"icp_filing"`
	PublicSecurityFiling string `json:"public_security_filing"`
	BilibiliURL          string `json:"bilibili_url"`
	GithubURL            string `json:"github_url"`
	SteamURL             string `json:"steam_url"`
	Name                 string `json:"name"`
	Job                  string `json:"job"`
	Address              string `json:"address"`
	Email                string `json:"email"`
	QQImage              string `json:"qq_image"`
	WechatImage          string `json:"wechat_image"`
}

// HotItem 热搜项。
type HotItem struct {
	Index       int    `json:"index"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Popularity  string `json:"popularity"`
	URL         string `json:"url"`
}

// HotSearchData 热搜数据。
type HotSearchData struct {
	Source     string    `json:"source"`
	UpdateTime string    `json:"update_time"`
	HotList    []HotItem `json:"hot_list"`
}

// CalendarData 日历数据。
type CalendarData struct {
	Date         string `json:"date"`
	LunarDate    string `json:"lunar_date"`
	Ganzhi       string `json:"ganzhi"`
	Zodiac       string `json:"zodiac"`
	DayOfYear    string `json:"day_of_year"`
	SolarTerm    string `json:"solar_term"`
	Auspicious   string `json:"auspicious"`
	Inauspicious string `json:"inauspicious"`
}

// FooterLink 页脚链接。
type FooterLink struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}
