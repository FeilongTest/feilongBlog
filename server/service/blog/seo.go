package blog

import (
	"feilongBlog/global"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// SeoService 负责服务端 SEO 渲染、站点地图与搜索引擎推送
type SeoService struct {
}

// PageSeo 描述一个页面的 SEO 信息
type PageSeo struct {
	Title       string           // 页面标题
	Description string           // 页面描述
	Keywords    string           // 页面关键词
	Canonical   string           // 规范地址
	Robots      string           // 抓取指令
	OgType      string           // Open Graph 类型
	Image       string           // 分享图片
	Published   int              // 发布时间戳
	Modified    int              // 更新时间戳
	Author      string           // 作者
	Section     string           // 所属分类
	JsonLd      []map[string]any // 结构化数据
	Body        string           // 首屏直出内容
}

// seoService SEO 服务实例
var seoService = SeoService{}

// 首页模板缓存，前端重新部署后自动重新读取
var (
	seoIndexMu   sync.RWMutex
	seoIndexHTML string
	seoIndexAt   time.Time
)

const (
	// 首页模板中的 SEO 占位标记
	seoHeadPlaceholder = "<!--seo:head-->"
	seoAppPlaceholder  = "<div id=\"app\">"
	// 服务端注入的标签统一带上标记，方便前端接管后清理
	seoTagFlag = `data-seo="1"`
	// 模板缓存有效期，避免前端重新部署后必须重启后端
	seoIndexCacheTTL = 60 * time.Second
)

// SiteURL 返回站点主域名，未配置时返回空字符串
func (s *SeoService) SiteURL() string {
	return strings.TrimRight(strings.TrimSpace(global.BLOG_CONFIG.Seo.SiteUrl), "/")
}

// absoluteURL 把站内相对路径拼成绝对地址
func (s *SeoService) absoluteURL(path string) string {
	base := s.SiteURL()
	if base == "" {
		return ""
	}
	path = strings.TrimSpace(path)
	if path == "" || path == "/" {
		return base + "/"
	}
	return base + "/" + strings.TrimLeft(path, "/")
}

// loadIndexTemplate 读取前端构建产物的首页模板，文件变化后自动重新加载
func (s *SeoService) loadIndexTemplate() (string, error) {
	root := strings.TrimSpace(global.BLOG_CONFIG.Seo.WebRoot)
	if root == "" {
		return "", fmt.Errorf("未配置 seo.web-root，无法进行服务端渲染")
	}
	file := filepath.Join(root, "index.html")

	seoIndexMu.RLock()
	cached, cachedAt := seoIndexHTML, seoIndexAt
	seoIndexMu.RUnlock()
	if cached != "" && time.Since(cachedAt) < seoIndexCacheTTL {
		return cached, nil
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	template := string(content)
	if !strings.Contains(template, seoHeadPlaceholder) || !strings.Contains(template, seoAppPlaceholder) {
		return "", fmt.Errorf("前端模板 %s 缺少 SEO 占位标记，请确认已更新 web/index.html", file)
	}
	seoIndexMu.Lock()
	seoIndexHTML, seoIndexAt = template, time.Now()
	seoIndexMu.Unlock()
	return template, nil
}

// RenderPage 为前端路由渲染带完整 SEO 信息的 HTML，status 为建议返回的 HTTP 状态码
func (s *SeoService) RenderPage(path, rawQuery string) (page string, status int, err error) {
	template, err := s.loadIndexTemplate()
	if err != nil {
		return "", 500, err
	}
	seo, status := s.buildPageSeo(path, rawQuery)
	return s.inject(template, seo), status, nil
}

// inject 把 SEO 信息写入首页模板
func (s *SeoService) inject(template string, seo *PageSeo) string {
	result := replaceTitle(template, seo.Title)
	result = strings.Replace(result, seoHeadPlaceholder, s.renderHead(seo), 1)
	if seo.Body != "" {
		result = strings.Replace(result, seoAppPlaceholder, seoAppPlaceholder+seo.Body, 1)
	}
	return result
}

// replaceTitle 替换模板中的 title 标签
func replaceTitle(template, title string) string {
	start := strings.Index(template, "<title>")
	if start < 0 {
		return template
	}
	end := strings.Index(template[start:], "</title>")
	if end < 0 {
		return template
	}
	end += start + len("</title>")
	return template[:start] + "<title>" + html.EscapeString(title) + "</title>" + template[end:]
}
