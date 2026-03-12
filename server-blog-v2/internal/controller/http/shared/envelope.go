package shared

// Envelope 统一响应信封。
type Envelope struct {
	Code    string      `json:"code" example:"0000"`
	Message string      `json:"message" example:"ok"`
	Data    interface{} `json:"data,omitempty"`
}
