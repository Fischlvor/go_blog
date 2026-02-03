package middleware

import (
	"server/internal/model/database"
	"server/internal/service"
	"server/pkg/global"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/ua-parser/uap-go/uaparser"
	"go.uber.org/zap"
)

// LoginRecord 是一个中间件，用于记录登录日志
func LoginRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 异步记录日志
		go func() {
			gaodeService := service.ServiceGroupApp.GaodeService
			var address string
			ip := c.ClientIP()
			loginMethod := c.DefaultQuery("flag", "email") // 若未传递flag参数，则默认为"email"
			userAgent := c.Request.UserAgent()

			// 获取用户IP的地理位置
			address = getAddressFromIP(ip, gaodeService)

			// 解析用户的浏览器、操作系统和设备信息
			os, deviceInfo, browserInfo := parseUserAgent(userAgent)

			// 创建登录记录
			// 注意：登录接口是公开路由，没有 SSO 中间件
			// TokenNext 设置的是 "user_uuid"，不是 "claims"
			var userUUID uuid.UUID
			if val, exists := c.Get("user_uuid"); exists {
				userUUID = val.(uuid.UUID)
			}
			login := database.Login{
				UserUUID:    userUUID, // 从 context 获取（登录成功时设置）
				LoginMethod: loginMethod,
				IP:          ip,
				Address:     address,
				OS:          os,
				DeviceInfo:  deviceInfo,
				BrowserInfo: browserInfo,
				Status:      c.Writer.Status(),
			}

			// 将登录记录存储到数据库
			if err := global.DB.Create(&login).Error; err != nil {
				global.Log.Error("Failed to record login", zap.Error(err))
			}
		}()
	}
}

// 获取IP地址对应的地理位置信息
func getAddressFromIP(ip string, gaodeService service.GaodeService) string {
	res, err := gaodeService.GetLocationByIP(ip)
	if err != nil || res.Province == "" {
		return "未知"
	}
	if res.City != "" && res.Province != res.City {
		return res.Province + "-" + res.City
	}
	return res.Province
}

// 解析用户代理（User-Agent）字符串，提取操作系统、设备信息和浏览器信息
func parseUserAgent(userAgent string) (os, deviceInfo, browserInfo string) {
	os = userAgent
	deviceInfo = userAgent
	browserInfo = userAgent

	parser := uaparser.NewFromSaved()
	cli := parser.Parse(userAgent)
	os = cli.Os.Family
	deviceInfo = cli.Device.Family
	browserInfo = cli.UserAgent.Family

	return
}
