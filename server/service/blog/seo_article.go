package blog

import (
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"fmt"
	"strings"
)

// 正文直出长度上限，避免单页过大
const seoBodyLimit = 60000

// articleSeo 文章详情页 SEO 信息
func (s *SeoService) articleSeo(setting model.Setting, id int) (*PageSeo, bool) {
	var article model.Article
	if err := global.BLOG_DB.Where("id = ? AND status = ?", id, 0).First(&article).Error; err != nil {
		return nil, false
	}

	categoryName := ""
	if article.Fid > 0 {
		var category model.Category
		if err := global.BLOG_DB.Where("id = ?", article.Fid).First(&category).Error; err == nil {
			categoryName = category.Name
		}
	}
	author := setting.Author
	if author == "" {
		author = "飞龙"
	}

	canonical := s.absoluteURL(fmt.Sprintf("/content/%d/%d", article.Fid, article.ID))
	description := plainText(article.Content, article.ContentFormat, 150)
	if description == "" {
		description = fmt.Sprintf("%s 发布的最新文章：%s", setting.SiteName, article.Title)
	}

	seo := &PageSeo{
		Title:       fmt.Sprintf("%s - %s", article.Title, setting.SiteName),
		Description: description,
		Keywords:    joinNonEmpty(",", article.Title, categoryName, setting.SiteName),
		Canonical:   canonical,
		Robots:      "index, follow",
		OgType:      "article",
		Image:       article.Pic,
		Published:   article.Ctime,
		Modified:    article.EditTime,
		Author:      author,
		Section:     categoryName,
	}
	seo.JsonLd = append(seo.JsonLd, s.articleJsonLd(article, setting, categoryName, author, canonical, description))
	if categoryName != "" {
		seo.JsonLd = append(seo.JsonLd, s.breadcrumbJsonLd(categoryName, s.absoluteURL(fmt.Sprintf("/category/%d", article.Fid)), article.Title, canonical))
	}
	seo.Body = s.renderArticleBody(article, categoryName, author)
	return seo, true
}

// renderArticleBody 生成文章正文的首屏直出内容，供不执行 JS 的爬虫读取
func (s *SeoService) renderArticleBody(article model.Article, categoryName, author string) string {
	content := article.Content
	if strings.EqualFold(strings.TrimSpace(article.ContentFormat), "markdown") {
		content = markdownToHTML(content)
	} else {
		content = sanitizeArticleHTML(content)
	}
	if content == "" {
		content = fmt.Sprintf("<p>%s</p>", escapeHTML(plainText(article.Content, article.ContentFormat, 150)))
	}

	var b strings.Builder
	b.WriteString(`<div class="seo-prerender"><article>`)
	b.WriteString("<h1>")
	b.WriteString(escapeHTML(article.Title))
	b.WriteString("</h1><p>")
	if article.Ctime > 0 {
		b.WriteString(fmt.Sprintf("<time datetime=\"%s\">%s</time>", isoTime(article.Ctime), escapeHTML(localDate(article.Ctime))))
	}
	if categoryName != "" {
		b.WriteString(fmt.Sprintf(" · <a href=\"%s\">%s</a>", escapeHTML(fmt.Sprintf("/category/%d", article.Fid)), escapeHTML(categoryName)))
	}
	if author != "" {
		b.WriteString(" · ")
		b.WriteString(escapeHTML(author))
	}
	b.WriteString("</p><div class=\"seo-prerender-content\">")
	b.WriteString(truncateRunes(content, seoBodyLimit))
	b.WriteString("</div></article></div>")
	return b.String()
}

// articleJsonLd 文章结构化数据
func (s *SeoService) articleJsonLd(article model.Article, setting model.Setting, categoryName, author, canonical, description string) map[string]any {
	data := map[string]any{
		"@context":         "https://schema.org",
		"@type":            "BlogPosting",
		"headline":         article.Title,
		"description":      description,
		"mainEntityOfPage": map[string]any{"@type": "WebPage", "@id": canonical},
		"url":              canonical,
		"author":           map[string]any{"@type": "Person", "name": author},
		"publisher": map[string]any{
			"@type": "Organization",
			"name":  setting.SiteName,
		},
	}
	if article.Ctime > 0 {
		data["datePublished"] = isoTime(article.Ctime)
	}
	if article.EditTime > 0 {
		data["dateModified"] = isoTime(article.EditTime)
	}
	if categoryName != "" {
		data["articleSection"] = categoryName
	}
	if article.Pic != "" {
		data["image"] = []string{article.Pic}
	}
	return data
}

// breadcrumbJsonLd 面包屑结构化数据
func (s *SeoService) breadcrumbJsonLd(categoryName, categoryURL, articleTitle, articleURL string) map[string]any {
	return map[string]any{
		"@context": "https://schema.org",
		"@type":    "BreadcrumbList",
		"itemListElement": []map[string]any{
			{"@type": "ListItem", "position": 1, "name": "首页", "item": s.absoluteURL("/")},
			{"@type": "ListItem", "position": 2, "name": categoryName, "item": categoryURL},
			{"@type": "ListItem", "position": 3, "name": articleTitle, "item": articleURL},
		},
	}
}

// joinNonEmpty 使用分隔符连接非空字符串
func joinNonEmpty(sep string, values ...string) string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			result = append(result, value)
		}
	}
	return strings.Join(result, sep)
}
