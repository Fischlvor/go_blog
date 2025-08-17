package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"server/global"
	"server/model/other"
	"server/utils"
)

// JwtService 提供与高德相关的服务
type GaodeService struct {
}

// GetLocationByIP 根据IP地址获取地理位置信息
func (gaodeService *GaodeService) GetLocationByIP(ip string) (other.IPResponse, error) {
	data := other.IPResponse{}

	if is, err := isPrivateIP(ip); is {
		data.Adcode = "440113"
		return data, err
	}
	key := global.Config.Gaode.Key
	urlStr := "https://restapi.amap.com/v3/ip"
	method := "GET"
	params := map[string]string{
		"ip":  ip,
		"key": key,
	}
	res, err := utils.HttpRequest(urlStr, method, nil, params, nil)
	if err != nil {
		return data, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return data, fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}

	byteData, err := io.ReadAll(res.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(byteData, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

// GetWeatherByAdcode 根据城市编码获取实时天气信息
func (gaodeService *GaodeService) GetWeatherByAdcode(adcode string) (other.Live, error) {
	data := other.WeatherResponse{}
	key := global.Config.Gaode.Key
	urlStr := "https://restapi.amap.com/v3/weather/weatherInfo"
	method := "GET"
	params := map[string]string{
		"city": adcode,
		"key":  key,
	}
	res, err := utils.HttpRequest(urlStr, method, nil, params, nil)
	if err != nil {
		return other.Live{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return other.Live{}, fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}

	byteData, err := io.ReadAll(res.Body)
	if err != nil {
		return other.Live{}, err
	}

	err = json.Unmarshal(byteData, &data)
	if err != nil {
		return other.Live{}, err
	}

	// 检查是否有返回的天气数据
	if len(data.Lives) == 0 {
		return other.Live{}, fmt.Errorf("no live weather data available") // 没有天气数据时返回错误
	}

	// 返回当天的天气数据
	return data.Lives[0], nil
}

// isPrivateIP 判断给定的IP地址是否为内网IP
// 支持IPv4和IPv6地址的判断
// 返回true表示是内网IP，false表示是公网IP，error表示解析失败
func isPrivateIP(ip string) (bool, error) {
	// 解析输入的IP地址
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false, errors.New("无效的IP地址格式")
	}

	// IPv4内网地址段判断
	if ip4 := parsedIP.To4(); ip4 != nil {
		return isPrivateIPv4(ip4), nil
	}

	// IPv6内网地址段判断
	return isPrivateIPv6(parsedIP), nil
}

// isPrivateIPv4 判断IPv4地址是否为内网地址
func isPrivateIPv4(ip net.IP) bool {
	// RFC 1918: 10.0.0.0/8
	if ip[0] == 10 {
		return true
	}

	// RFC 1918: 172.16.0.0/12
	if ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31 {
		return true
	}

	// RFC 1918: 192.168.0.0/16
	if ip[0] == 192 && ip[1] == 168 {
		return true
	}

	// RFC 6598: 100.64.0.0/10 (共享地址空间)
	if ip[0] == 100 && ip[1] >= 64 && ip[1] <= 127 {
		return true
	}

	// RFC 5735: 127.0.0.0/8 (回环地址)
	if ip[0] == 127 {
		return true
	}

	// RFC 5735: 169.254.0.0/16 (链路本地地址)
	if ip[0] == 169 && ip[1] == 254 {
		return true
	}

	// RFC 5737: 192.0.0.0/24, 192.0.2.0/24, 198.51.100.0/24, 203.0.113.0/24 (测试地址)
	if (ip[0] == 192 && ip[1] == 0 && (ip[2] == 0 || ip[2] == 2)) ||
		(ip[0] == 198 && ip[1] == 51 && ip[2] == 100) ||
		(ip[0] == 203 && ip[1] == 0 && ip[2] == 113) {
		return true
	}

	// RFC 3849: 2001:db8::/32 (文档地址)
	if ip[0] == 0x20 && ip[1] == 0x01 && ip[2] == 0xdb && ip[3] == 0x8 {
		return true
	}

	return false
}

// isPrivateIPv6 判断IPv6地址是否为内网地址
func isPrivateIPv6(ip net.IP) bool {
	// 链路本地单播地址 fe80::/10
	if ip[0] == 0xfe && (ip[1]&0xc0) == 0x80 {
		return true
	}

	// 唯一本地地址 fc00::/7
	if (ip[0] & 0xfe) == 0xfc {
		return true
	}

	// 回环地址 ::1
	if ip.Equal(net.ParseIP("::1")) {
		return true
	}

	// 文档地址 2001:db8::/32
	if ip[0] == 0x20 && ip[1] == 0x01 && ip[2] == 0xdb && ip[3] == 0x8 {
		return true
	}

	return false
}
