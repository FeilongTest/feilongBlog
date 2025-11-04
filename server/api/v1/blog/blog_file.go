package blog

import (
	"feilongBlog/global"
	"feilongBlog/model/common/response"
	"feilongBlog/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FileApi struct{}

var fileUploadService = service.ServiceGroupApp.BlogService.FileService

func (a *FileApi) UploadFile(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.BLOG_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	filePath, fileName, err := fileUploadService.UploadFile(header) // 文件上传后拿到文件路径
	if err != nil {
		global.BLOG_LOG.Error("上传文件失败!", zap.Error(err))
		response.FailWithMessage("上传文件失败", c)
		return
	}
	response.OkWithDetailed(map[string]string{
		"url": filePath, // 图片 src ，必须
		"alt": fileName, // 图片描述文字，非必须
	}, "上传成功", c)
}
