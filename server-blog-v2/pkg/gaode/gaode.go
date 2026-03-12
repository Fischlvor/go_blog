package gaode

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Client 高德地图 API 客户端。
type Client struct {
	key        string
	httpClient *http.Client
}

// NewClient 创建高德客户端。
func NewClient(key string) *Client {
	return &Client{
		key: key,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// IPResponse IP 定位响应。
type IPResponse struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Province string `json:"province"`
	City     string `json:"city"`
	Adcode   string `json:"adcode"`
}

// WeatherResponse 天气响应。
type WeatherResponse struct {
	Status string `json:"status"`
	Count  string `json:"count"`
	Info   string `json:"info"`
	Lives  []Live `json:"lives"`
}

// Live 实时天气。
type Live struct {
	Province      string `json:"province"`
	City          string `json:"city"`
	Adcode        string `json:"adcode"`
	Weather       string `json:"weather"`
	Temperature   string `json:"temperature"`
	WindDirection string `json:"winddirection"`
	WindPower     string `json:"windpower"`
	Humidity      string `json:"humidity"`
	ReportTime    string `json:"reporttime"`
}

// GetLocationByIP 根据 IP 获取位置信息。
func (c *Client) GetLocationByIP(ip string) (*IPResponse, error) {
	// 检查是否为内网 IP，如果是则使用默认城市
	if isPrivateIP(ip) {
		return &IPResponse{
			Status: "1",
			Adcode: "440113", // 番禺
		}, nil
	}

	params := url.Values{}
	params.Set("key", c.key)
	params.Set("ip", ip)

	resp, err := c.httpClient.Get("https://restapi.amap.com/v3/ip?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result IPResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("gaode api error: %s", result.Info)
	}

	return &result, nil
}

// isPrivateIP 判断是否为内网 IP。
func isPrivateIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return true // 无效 IP 当作内网处理
	}

	// IPv4 内网地址
	if ip4 := parsedIP.To4(); ip4 != nil {
		// 10.0.0.0/8
		if ip4[0] == 10 {
			return true
		}
		// 172.16.0.0/12
		if ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31 {
			return true
		}
		// 192.168.0.0/16
		if ip4[0] == 192 && ip4[1] == 168 {
			return true
		}
		// 127.0.0.0/8 (回环)
		if ip4[0] == 127 {
			return true
		}
	}

	// IPv6 回环 ::1
	if parsedIP.Equal(net.ParseIP("::1")) {
		return true
	}

	return false
}

// GetWeatherByAdcode 根据城市编码获取天气。
func (c *Client) GetWeatherByAdcode(adcode string) (*Live, error) {
	params := url.Values{}
	params.Set("key", c.key)
	params.Set("city", adcode)
	params.Set("extensions", "base")

	resp, err := c.httpClient.Get("https://restapi.amap.com/v3/weather/weatherInfo?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result WeatherResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("gaode api error: %s", result.Info)
	}

	if len(result.Lives) == 0 {
		return nil, fmt.Errorf("no weather data available")
	}

	return &result.Lives[0], nil
}

// GetWeatherString 获取格式化的天气字符串。
func (c *Client) GetWeatherString(ip string) (string, error) {
	location, err := c.GetLocationByIP(ip)
	if err != nil {
		return "", err
	}

	live, err := c.GetWeatherByAdcode(location.Adcode)
	if err != nil {
		return "", err
	}

	return c.FormatWeather(live), nil
}

// FormatWeather 格式化天气信息。
func (c *Client) FormatWeather(live *Live) string {
	return fmt.Sprintf("地区：%s-%s 天气：%s 温度：%s°C 风向：%s 风级：%s 湿度：%s%%",
		live.Province, live.City, live.Weather, live.Temperature,
		live.WindDirection, live.WindPower, live.Humidity)
}
