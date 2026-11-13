package blog

import (
	"bytes"
	"encoding/xml"
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"fmt"
	"strings"
	"sync"
	"time"
)

// 站点地图缓存，避免爬虫每次抓取都查询数据库
var (
	seoSitemapMu  sync.RWMutex
	seoSitemapXML string
	seoSitemapAt  time.Time
)

// sitemapItem 站点地图中的单个地址
type sitemapItem struct {
	Loc        string
	LastMod    string
	ChangeFreq string
	Priority   string
}

// Sitemap 生成站点地图，结果带缓存
func (s *SeoService) Sitemap() string {
	ttl := global.BLOG_CONFIG.Seo.SitemapTTL
	if ttl <= 0 {
		ttl = 600
	}
	seoSitemapMu.RLock()
	cached, cachedAt := seoSitemapXML, seoSitemapAt
	seoSitemapMu.RUnlock()
	if cached != "" && time.Since(cachedAt) < time.Duration(ttl)*time.Second {
		return cached
	}

	result := s.buildSitemap()
	seoSitemapMu.Lock()
	seoSitemapXML, seoSitemapAt = result, time.Now()
	seoSitemapMu.Unlock()
	return result
}

// ResetSitemap 内容变化后清空站点地图缓存
func (s *SeoService) ResetSitemap() {
	seoSitemapMu.Lock()
	seoSitemapXML, seoSitemapAt = "", time.Time{}
	seoSitemapMu.Unlock()
}

// buildSitemap 从数据库读取内容生成站点地图
func (s *SeoService) buildSitemap() string {
	siteURL := s.SiteURL()
	if siteURL == "" {
		return ""
	}

	items := []sitemapItem{
		{Loc: siteURL + "/", ChangeFreq: "daily", Priority: "1.0"},
		{Loc: siteURL + "/archive", ChangeFreq: "daily", Priority: "0.7"},
		{Loc: siteURL + "/about", ChangeFreq: "monthly", Priority: "0.4"},
		{Loc: siteURL + "/tags", ChangeFreq: "weekly", Priority: "0.5"},
	}

	var categories []model.Category
	if err := global.BLOG_DB.Order("sort desc").Find(&categories).Error; err == nil {
		for _, category := range categories {
			items = append(items, sitemapItem{
				Loc:        fmt.Sprintf("%s/category/%d", siteURL, category.ID),
				ChangeFreq: "weekly",
				Priority:   "0.6",
			})
		}
	}

	if tags, err := (&TagService{}).List(true); err == nil {
		for _, tag := range tags {
			items = append(items, sitemapItem{Loc: fmt.Sprintf("%s/tag/%d", siteURL, tag.ID), ChangeFreq: "weekly", Priority: "0.5"})
		}
	}
	var articles []model.Article
	if err := global.BLOG_DB.Where("status = ?", 0).Order("ctime desc").Find(&articles).Error; err == nil {
		for _, article := range articles {
			lastMod := localDate(article.EditTime)
			if lastMod == "" {
				lastMod = localDate(article.Ctime)
			}
			items = append(items, sitemapItem{
				Loc:        fmt.Sprintf("%s/content/%d/%d", siteURL, article.Fid, article.ID),
				LastMod:    lastMod,
				ChangeFreq: "monthly",
				Priority:   "0.8",
			})
		}
	}

	return renderSitemap(items)
}

// renderSitemap 把地址列表转换成 sitemap.xml
func renderSitemap(items []sitemapItem) string {
	var b bytes.Buffer
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")
	for _, item := range items {
		b.WriteString("  <url>\n    <loc>")
		b.WriteString(escapeXML(item.Loc))
		b.WriteString("</loc>\n")
		if item.LastMod != "" {
			b.WriteString("    <lastmod>")
			b.WriteString(escapeXML(item.LastMod))
			b.WriteString("</lastmod>\n")
		}
		if item.ChangeFreq != "" {
			b.WriteString("    <changefreq>")
			b.WriteString(item.ChangeFreq)
			b.WriteString("</changefreq>\n")
		}
		if item.Priority != "" {
			b.WriteString("    <priority>")
			b.WriteString(item.Priority)
			b.WriteString("</priority>\n")
		}
		b.WriteString("  </url>\n")
	}
	b.WriteString("</urlset>\n")
	return b.String()
}

// escapeXML 转义 XML 文本
func escapeXML(text string) string {
	var b bytes.Buffer
	if err := xml.EscapeText(&b, []byte(text)); err != nil {
		return text
	}
	return b.String()
}

// Robots 生成 robots.txt 内容。
// 具体爬虫要写在通配规则前面。Cloudflare 如果在响应前插入同名 Disallow，需要在控制台关闭托管 robots.txt。
func (s *SeoService) Robots() string {
	var b strings.Builder
	for _, agent := range []string{"Baiduspider", "Baiduspider-image", "Googlebot", "bingbot", "Sogou web spider", "360Spider", "YisouSpider"} {
		b.WriteString("User-agent: ")
		b.WriteString(agent)
		b.WriteString("\nAllow: /\n")
		writeRobotExclusions(&b)
		b.WriteString("\n")
	}
	b.WriteString("User-agent: *\n")
	b.WriteString("Allow: /\n")
	writeRobotExclusions(&b)
	b.WriteString("\n")
	if siteURL := s.SiteURL(); siteURL != "" {
		b.WriteString("Sitemap: ")
		b.WriteString(siteURL)
		b.WriteString("/sitemap.xml\n")
	}
	return b.String()
}

// CanonicalRedirect 把关于页和归档页的大小写、尾部斜杠统一到规范地址
func (s *SeoService) CanonicalRedirect(path string) string {
	clean := strings.TrimRight(strings.TrimSpace(path), "/")
	if clean == "" {
		return ""
	}
	switch strings.ToLower(clean) {
	case "/about":
		if path != "/about" {
			return "/about"
		}
	case "/archive":
		if path != "/archive" {
			return "/archive"
		}
	}
	return ""
}

func writeRobotExclusions(b *strings.Builder) {
	for _, path := range []string{"/admin", "/reading-list", "/dashboard", "/sign-in", "/blog/", "/article/", "/*?keyword=", "/*&keyword=", "/*?page=", "/*&page="} {
		b.WriteString("Disallow: " + path + "\n")
	}
}
