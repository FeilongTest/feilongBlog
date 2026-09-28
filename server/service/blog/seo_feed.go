package blog

import (
	"bytes"
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"fmt"
	"strings"
	"sync"
	"time"
)

// 订阅源输出的文章数量
const seoFeedLimit = 20

// 订阅内容缓存，内容变化后自动刷新
var (
	seoFeedMu  sync.RWMutex
	seoFeedXML string
	seoFeedAt  time.Time
)

// Feed 生成 RSS 2.0 订阅内容，结果带缓存
func (s *SeoService) Feed() string {
	ttl := global.BLOG_CONFIG.Seo.SitemapTTL
	if ttl <= 0 {
		ttl = 600
	}
	seoFeedMu.RLock()
	cached, cachedAt := seoFeedXML, seoFeedAt
	seoFeedMu.RUnlock()
	if cached != "" && time.Since(cachedAt) < time.Duration(ttl)*time.Second {
		return cached
	}

	result := s.buildFeed()
	seoFeedMu.Lock()
	seoFeedXML, seoFeedAt = result, time.Now()
	seoFeedMu.Unlock()
	return result
}

// ResetFeed 内容变化后清空订阅缓存
func (s *SeoService) ResetFeed() {
	seoFeedMu.Lock()
	seoFeedXML, seoFeedAt = "", time.Time{}
	seoFeedMu.Unlock()
}

// buildFeed 从数据库读取文章生成订阅内容
func (s *SeoService) buildFeed() string {
	siteURL := s.SiteURL()
	if siteURL == "" {
		return ""
	}
	setting, err := settingService.GetSetting()
	if err != nil {
		setting = defaultSetting()
	}

	var articles []model.Article
	if err := global.BLOG_DB.Where("status = ?", 0).Order("ctime desc").Limit(seoFeedLimit).Find(&articles).Error; err != nil {
		return ""
	}

	// 分类名称用于输出 category 标签
	categoryNames := s.categoryNameMap()

	var b bytes.Buffer
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">` + "\n")
	b.WriteString("  <channel>\n")
	b.WriteString(fmt.Sprintf("    <title>%s</title>\n", escapeXML(setting.SiteName)))
	b.WriteString(fmt.Sprintf("    <link>%s/</link>\n", escapeXML(siteURL)))
	b.WriteString(fmt.Sprintf("    <description>%s</description>\n", escapeXML(setting.Description)))
	b.WriteString("    <language>zh-CN</language>\n")
	if len(articles) > 0 {
		b.WriteString(fmt.Sprintf("    <lastBuildDate>%s</lastBuildDate>\n", rssTime(articles[0].EditTime)))
	}
	b.WriteString(fmt.Sprintf("    <generator>%s</generator>\n", escapeXML(setting.SiteName)))
	b.WriteString(fmt.Sprintf("    <atom:link href=%q rel=\"self\" type=\"application/rss+xml\" />\n", siteURL+"/rss.xml"))

	for _, article := range articles {
		link := s.ArticleURL(article.Fid, article.ID)
		b.WriteString("    <item>\n")
		b.WriteString(fmt.Sprintf("      <title>%s</title>\n", escapeXML(article.Title)))
		b.WriteString(fmt.Sprintf("      <link>%s</link>\n", escapeXML(link)))
		b.WriteString(fmt.Sprintf("      <guid isPermaLink=\"true\">%s</guid>\n", escapeXML(link)))
		b.WriteString(fmt.Sprintf("      <pubDate>%s</pubDate>\n", rssTime(article.Ctime)))
		if name := categoryNames[article.Fid]; name != "" {
			b.WriteString(fmt.Sprintf("      <category>%s</category>\n", escapeXML(name)))
		}
		body := articleContentHTML(article)
		if body == "" {
			body = plainText(article.Content, article.ContentFormat, 200)
		}
		b.WriteString("      <description><![CDATA[")
		b.WriteString(escapeCDATA(body))
		b.WriteString("]]></description>\n")
		b.WriteString("    </item>\n")
	}

	b.WriteString("  </channel>\n</rss>\n")
	return b.String()
}

// categoryNameMap 返回分类 ID 到名称的映射
func (s *SeoService) categoryNameMap() map[int]string {
	result := make(map[int]string)
	var categories []model.Category
	if err := global.BLOG_DB.Find(&categories).Error; err != nil {
		return result
	}
	for _, category := range categories {
		result[int(category.ID)] = category.Name
	}
	return result
}

// rssTime 把时间戳转换成 RSS 使用的时间格式
func rssTime(timestamp int) string {
	if timestamp <= 0 {
		return time.Now().Format(time.RFC1123Z)
	}
	return time.Unix(int64(timestamp), 0).Format(time.RFC1123Z)
}

// escapeCDATA 处理 CDATA 段落中不能出现的结束标记
func escapeCDATA(text string) string {
	return strings.ReplaceAll(text, "]]>", "]]]]><![CDATA[>")
}
