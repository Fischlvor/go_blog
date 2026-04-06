package output

// SettingValueWithKey 配置值及其 setting_key。
type SettingValueWithKey struct {
	Value      string `json:"value"`
	SettingKey string `json:"setting_key"`
}

// WebsiteInfo 网站信息。
type WebsiteInfo struct {
	Avatar          SettingValueWithKey `json:"avatar"`
	Title           SettingValueWithKey `json:"title"`
	Description     SettingValueWithKey `json:"description"`
	ProfileIntro    SettingValueWithKey `json:"profile_intro"`
	TechStack       SettingValueWithKey `json:"tech_stack"`
	WorkExperiences SettingValueWithKey `json:"work_experiences"`
	Version         SettingValueWithKey `json:"version"`
	CreatedAt       SettingValueWithKey `json:"created_at"`
	ICPFiling       SettingValueWithKey `json:"icp_filing"`
	BilibiliURL     SettingValueWithKey `json:"bilibili_url"`
	GithubURL       SettingValueWithKey `json:"github_url"`
	SteamURL        SettingValueWithKey `json:"steam_url"`
	Name            SettingValueWithKey `json:"name"`
	Job             SettingValueWithKey `json:"job"`
	Address         SettingValueWithKey `json:"address"`
	Email           SettingValueWithKey `json:"email"`
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

// SiteSettingDetail 站点配置详情。
type SiteSettingDetail struct {
	ID           int64  `json:"id"`
	SettingKey   string `json:"setting_key"`
	SettingValue string `json:"setting_value"`
	SettingType  string `json:"setting_type"`
	Description  string `json:"description"`
	IsPublic     bool   `json:"is_public"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
