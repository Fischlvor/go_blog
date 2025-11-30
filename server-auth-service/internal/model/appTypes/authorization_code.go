package appTypes

import "time"

// AuthorizationCode 授权码存储结构
type AuthorizationCode struct {
	Code         string
	UserUUID     string
	AppID        string
	RedirectURI  string
	ExpiresAt    time.Time
	AccessToken  string
	RefreshToken string
}
