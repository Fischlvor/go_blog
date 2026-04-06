package request

// WebsiteConfig 网站配置更新请求。
type WebsiteConfig struct {
	Avatar          SettingField `json:"avatar"`
	Title           SettingField `json:"title"`
	Description     SettingField `json:"description"`
	ProfileIntro    SettingField `json:"profile_intro"`
	TechStack       SettingField `json:"tech_stack"`
	WorkExperiences SettingField `json:"work_experiences"`
	Version         SettingField `json:"version"`
	CreatedAt       SettingField `json:"created_at"`
	ICPFiling       SettingField `json:"icp_filing"`
	BilibiliURL     SettingField `json:"bilibili_url"`
	GithubURL       SettingField `json:"github_url"`
	SteamURL        SettingField `json:"steam_url"`
	Name            SettingField `json:"name"`
	Job             SettingField `json:"job"`
	Address         SettingField `json:"address"`
	Email           SettingField `json:"email"`
}

// SettingField 配置字段。
type SettingField struct {
	Value      string `json:"value"`
	SettingKey string `json:"setting_key"`
}
