package blog

import (
	"feilongBlog/global"
	"feilongBlog/model/common/response"
	"feilongBlog/service/blog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SeoApi struct {
}

var seoService = blog.SeoService{}

// 后端接口前缀，这些地址不能返回 HTML
// 注意：后台页面使用 /admin 开头的地址，因此这里不能包含 /admin，
// 真正的后台接口已经注册为 Gin 路由，不会走到 NoRoute
var apiPathPrefixes = []string{"/blog", "/base", "/health"}

// RenderPage 为前台路由注入 SEO 信息，未命中接口的请求都会走到这里
func (a *SeoApi) RenderPage(c *gin.Context) {
	path := c.Request.URL.Path

	// IndexNow 密钥文件，便于搜索引擎校验推送来源
	if key := seoService.IndexNowKey(); key != "" && path == "/"+key+".txt" {
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(key))
		return
	}

	if isAPIPath(path) {
		c.JSON(http.StatusNotFound, response.Response{Code: response.ERROR, Data: map[string]interface{}{}, Msg: "接口不存在"})
		return
	}

	html, status, err := seoService.RenderPage(path, c.Request.URL.RawQuery)
	if err != nil {
		global.BLOG_LOG.Error("服务端渲染页面失败", zap.String("path", path), zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.Response{Code: response.ERROR, Data: map[string]interface{}{}, Msg: "页面渲染失败"})
		return
	}
	c.Header("Cache-Control", "no-cache")
	c.Data(status, "text/html; charset=utf-8", []byte(html))
}

// Sitemap 输出站点地图
func (a *SeoApi) Sitemap(c *gin.Context) {
	content := seoService.Sitemap()
	if content == "" {
		c.Data(http.StatusNotFound, "text/plain; charset=utf-8", []byte("未配置 seo.site-url，站点地图不可用"))
		return
	}
	c.Header("Cache-Control", "public, max-age=600")
	c.Data(http.StatusOK, "application/xml; charset=utf-8", []byte(content))
}

// Robots 输出 robots.txt
func (a *SeoApi) Robots(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=3600")
	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(seoService.Robots()))
}

// isAPIPath 判断请求是否属于后端接口
func isAPIPath(path string) bool {
	for _, prefix := range apiPathPrefixes {
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}
	return false
}
