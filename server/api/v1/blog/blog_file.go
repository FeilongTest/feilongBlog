package blog

import (
	"feilongBlog/global"
	"feilongBlog/model/common/response"
	"feilongBlog/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
	"strings"
)

type FileApi struct{}

var fileUploadService = service.ServiceGroupApp.BlogService.FileService

func (a *FileApi) UploadFile(c *gin.Context) {
	maxMB := global.BLOG_CONFIG.R2.MaxSizeMB
	if maxMB <= 0 {
		maxMB = 10
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxMB<<20)
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.BLOG_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	file, err := header.Open()
	if err != nil {
		response.FailWithMessage("无法读取上传文件", c)
		return
	}
	defer file.Close()
	sniff := make([]byte, 512)
	n, _ := file.Read(sniff)
	contentType := http.DetectContentType(sniff[:n])
	if !strings.HasPrefix(contentType, "image/") || contentType == "image/svg+xml" {
		response.FailWithMessage("仅支持常见位图格式，禁止 SVG", c)
		return
	}
	extension := map[string]string{"image/jpeg": ".jpg", "image/png": ".png", "image/gif": ".gif", "image/webp": ".webp", "image/bmp": ".bmp", "image/x-ms-bmp": ".bmp", "image/x-icon": ".ico", "image/vnd.microsoft.icon": ".ico"}[contentType]
	if extension == "" {
		response.FailWithMessage("不支持该图片格式，请使用 JPG、PNG、GIF、WebP 或 BMP", c)
		return
	}
	header.Filename = strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename)) + extension
	header.Header.Set("Content-Type", contentType)
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
