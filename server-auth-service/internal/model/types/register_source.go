package types

import "encoding/json"

// RegisterSource 注册来源枚举
type RegisterSource int

const (
	Email  RegisterSource = iota // 0: 邮箱注册
	QQ                           // 1: QQ登录
	Wechat                       // 2: 微信登录
	Github                       // 3: Github登录
)

func (r RegisterSource) String() string {
	switch r {
	case Email:
		return "email"
	case QQ:
		return "qq"
	case Wechat:
		return "wechat"
	case Github:
		return "github"
	default:
		return "unknown"
	}
}

// MarshalJSON 实现 JSON 序列化
func (r RegisterSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

// UnmarshalJSON 实现 JSON 反序列化
func (r *RegisterSource) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "email":
		*r = Email
	case "qq":
		*r = QQ
	case "wechat":
		*r = Wechat
	case "github":
		*r = Github
	default:
		*r = Email
	}
	return nil
}
