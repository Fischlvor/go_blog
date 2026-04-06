package input

// UpsertSiteSetting 站点配置写入参数。
type UpsertSiteSetting struct {
	SettingKey   string
	SettingValue string
	SettingType  string
	Description  string
	IsPublic     bool
}
