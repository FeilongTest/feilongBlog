package blog

import (
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"fmt"
	"strconv"
	"strings"
)

// 列表页直出文章数量
const seoListLimit = 10

// buildPageSeo 根据请求路径构建页面 SEO 信息，第二个返回值是建议的 HTTP 状态码
func (s *SeoService) buildPageSeo(path, rawQuery string) (*PageSeo, int) {
	setting, err := settingService.GetSetting()
	if err != nil {
		setting = defaultSetting()
	}
	cleanPath := strings.TrimRight(strings.TrimSpace(path), "/")
	if cleanPath == "" {
		cleanPath = "/"
	}

	switch {
	case cleanPath == "/" || cleanPath == "/index":
		return s.homeSeo(setting), 200
	case strings.HasPrefix(cleanPath, "/category/"):
		// 路径形如 /category/:id，分类 ID 位于第 2 段
		if id := pathSegmentID(cleanPath, 2); id > 0 {
			if seo, ok := s.categorySeo(setting, id); ok {
				return seo, 200
			}
		}
	case strings.HasPrefix(cleanPath, "/content/"):
		// 路径形如 /content/:fid/:id，文章 ID 位于第 3 段
		if id := pathSegmentID(cleanPath, 3); id > 0 {
			if seo, ok := s.articleSeo(setting, id); ok {
				return seo, 200
			}
		}
	}

	// 登录页、后台与不存在的地址不参与收录
	return s.plainSeo(setting, cleanPath), statusForPath(cleanPath)
}

// statusForPath 判断路径应该返回的状态码，SPA 内部页面返回 200，其余返回 404
func statusForPath(path string) int {
	switch {
	case path == "/sign-in", path == "/404", path == "/500",
		path == "/dashboard", strings.HasPrefix(path, "/admin"):
		return 200
	default:
		return 404
	}
}

// pathSegmentID 取出路径中第 index 段（从 1 开始）的数字
func pathSegmentID(path string, index int) int {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) < index {
		return 0
	}
	id, err := strconv.Atoi(segments[index-1])
	if err != nil || id <= 0 {
		return 0
	}
	return id
}

// homeSeo 首页 SEO 信息
func (s *SeoService) homeSeo(setting model.Setting) *PageSeo {
	seo := &PageSeo{
		Title:       fmt.Sprintf("%s - %s", setting.SiteName, setting.SiteSlogan),
		Description: setting.Description,
		Keywords:    setting.Keywords,
		Canonical:   s.absoluteURL("/"),
		Robots:      "index, follow",
		OgType:      "website",
		Author:      setting.Author,
	}
	seo.JsonLd = append(seo.JsonLd, s.websiteJsonLd(setting))
	seo.Body = s.renderListBody(setting.SiteName, setting.Description, s.recentArticles(0, seoListLimit))
	return seo
}

// categorySeo 分类页 SEO 信息
func (s *SeoService) categorySeo(setting model.Setting, fid int) (*PageSeo, bool) {
	var category model.Category
	if err := global.BLOG_DB.Where("id = ?", fid).First(&category).Error; err != nil {
		return nil, false
	}
	articles := s.recentArticles(fid, seoListLimit)
	title := fmt.Sprintf("%s 分类下的文章 - %s", category.Name, setting.SiteName)
	description := fmt.Sprintf("「%s」分类下的全部文章，共 %d 篇，由 %s 整理发布。", category.Name, len(articles), setting.SiteName)
	seo := &PageSeo{
		Title:       title,
		Description: description,
		Keywords:    joinNonEmpty(",", category.Name, setting.SiteName, setting.Keywords),
		Canonical:   s.absoluteURL(fmt.Sprintf("/category/%d", fid)),
		Robots:      "index, follow",
		OgType:      "website",
		Author:      setting.Author,
		Section:     category.Name,
	}
	seo.JsonLd = append(seo.JsonLd, s.collectionJsonLd(title, description, seo.Canonical))
	seo.Body = s.renderListBody(category.Name, description, articles)
	return seo, true
}

// plainSeo 不参与收录页面的 SEO 信息
func (s *SeoService) plainSeo(setting model.Setting, path string) *PageSeo {
	return &PageSeo{
		Title:       setting.SiteName,
		Description: setting.Description,
		Keywords:    setting.Keywords,
		Canonical:   s.absoluteURL(path),
		Robots:      "noindex, follow",
		OgType:      "website",
		Author:      setting.Author,
	}
}

// recentArticles 读取指定分类下的最新文章
func (s *SeoService) recentArticles(fid int, limit int) []model.Article {
	db := global.BLOG_DB.Model(&model.Article{}).Where("status = ?", 0)
	if fid > 0 {
		db = db.Where("fid = ?", fid)
	}
	var list []model.Article
	if err := db.Order("istop desc, ctime desc").Limit(limit).Find(&list).Error; err != nil {
		return nil
	}
	return list
}

// renderListBody 生成列表页的首屏直出内容
func (s *SeoService) renderListBody(title, description string, list []model.Article) string {
	var b strings.Builder
	b.WriteString(`<div class="seo-prerender">`)
	b.WriteString(fmt.Sprintf("<h1>%s</h1>", escapeHTML(title)))
	if description != "" {
		b.WriteString(fmt.Sprintf("<p>%s</p>", escapeHTML(description)))
	}
	if len(list) == 0 {
		b.WriteString("</div>")
		return b.String()
	}
	b.WriteString("<ul>")
	for _, article := range list {
		link := fmt.Sprintf("/content/%d/%d", article.Fid, article.ID)
		b.WriteString("<li><h2><a href=\"")
		b.WriteString(escapeHTML(link))
		b.WriteString("\">")
		b.WriteString(escapeHTML(article.Title))
		b.WriteString("</a></h2>")
		if article.Ctime > 0 {
			b.WriteString(fmt.Sprintf("<time datetime=\"%s\">%s</time>", isoTime(article.Ctime), escapeHTML(localDate(article.Ctime))))
		}
		if summary := plainText(article.Content, article.ContentFormat, 120); summary != "" {
			b.WriteString(fmt.Sprintf("<p>%s</p>", escapeHTML(summary)))
		}
		b.WriteString("</li>")
	}
	b.WriteString("</ul></div>")
	return b.String()
}

// websiteJsonLd 站点结构化数据
func (s *SeoService) websiteJsonLd(setting model.Setting) map[string]any {
	siteURL := s.SiteURL()
	data := map[string]any{
		"@context": "https://schema.org",
		"@type":    "WebSite",
		"name":     setting.SiteName,
		"url":      siteURL + "/",
	}
	if setting.Description != "" {
		data["description"] = setting.Description
	}
	if siteURL != "" {
		data["potentialAction"] = map[string]any{
			"@type":       "SearchAction",
			"target":      siteURL + "/index?keyword={search_term_string}",
			"query-input": "required name=search_term_string",
		}
	}
	return data
}

// collectionJsonLd 列表页结构化数据
func (s *SeoService) collectionJsonLd(name, description, url string) map[string]any {
	return map[string]any{
		"@context":    "https://schema.org",
		"@type":       "CollectionPage",
		"name":        name,
		"description": description,
		"url":         url,
	}
}
