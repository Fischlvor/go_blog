package api

import (
	"server/internal/model/database"
	"server/internal/model/request"
	"server/internal/model/response"
	"server/pkg/global"
	"server/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ImageApi struct {
}

// ImageUpload 上传图片
func (imageApi *ImageApi) ImageUpload(c *gin.Context) {
	_, header, err := c.Request.FormFile("image")
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	url, err := imageService.ImageUpload(header)
	if err != nil {
		global.Log.Error("Failed to upload image:", zap.Error(err))
		response.FailWithMessage("Failed to upload image", c)
		return
	}
	// local: 直接返回相对路径
	// qiniu: 数据库存 key，这里返回拼接域名后的完整 URL
	if global.Config.System.OssType == "qiniu" {
		url = global.Config.Qiniu.ImgPath + url
	}
	response.OkWithDetailed(response.ImageUpload{
		Url:     url,
		OssType: global.Config.System.OssType,
	}, "Successfully uploaded image", c)
}

// ImageDelete 删除图片
func (imageApi *ImageApi) ImageDelete(c *gin.Context) {
	var req request.ImageDelete
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = imageService.ImageDelete(req)
	if err != nil {
		global.Log.Error("Failed to delete image:", zap.Error(err))
		response.FailWithMessage("Failed to delete image", c)
		return
	}
	response.OkWithMessage("Successfully deleted image", c)
}

// ImageList 获取图片列表
func (imageApi *ImageApi) ImageList(c *gin.Context) {
	var pageInfo request.ImageList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	imageList, total, err := imageService.ImageList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get image list:", zap.Error(err))
		response.FailWithMessage("Failed to get image list", c)
		return
	}
	// 拼接 Image.URL 的对外 URL
	switch items := imageList.(type) {
	case []database.Image:
		for i := range items {
			items[i].URL = utils.PublicURLFromDB(items[i].URL)
		}
		imageList = items
	}
	response.OkWithData(response.PageResult{
		List:  imageList,
		Total: total,
	}, c)
}
