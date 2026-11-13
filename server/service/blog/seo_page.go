package blog

import (
	"errors"
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"fmt"
	"gorm.io/gorm"
	"net/url"
	"strconv"
	"strings"
)

// 列表页直出文章数量
const seoListLimit = 10

// buildPageSeo 根据请求路径构建页面 SEO 信息，第二个返回值是建议的 HTTP 状态码
func (s *SeoService) buildPageSeo(path, rawQuery string) (*PageSeo, int, error) {
	setting, err := settingService.GetSetting()
	if err != nil {
		return nil, 500, err
	}
	cleanPath := strings.TrimRight(strings.TrimSpace(path), "/")
	if cleanPath == "" {
		cleanPath = "/"
	}

	switch {
	case cleanPath == "/" || cleanPath == "/index":
		seo := s.homeSeo(setting)
		if err := s.applyListingQuery(seo, setting, 0, rawQuery); err != nil {
			return nil, 500, err
		}
		return seo, 200, nil
	case strings.EqualFold(cleanPath, "/about"):
		return s.aboutSeo(setting), 200, nil
	case strings.EqualFold(cleanPath, "/archive"):
		return s.archiveSeo(setting), 200, nil
	case strings.HasPrefix(cleanPath, "/category/") && len(strings.Split(strings.Trim(cleanPath, "/"), "/")) == 2:
		// 路径形如 /category/:id，分类 ID 位于第 2 段
		if id := pathSegmentID(cleanPath, 2); id > 0 {
			seo, err := s.categorySeo(setting, id)
			if err == nil {
				if err := s.applyListingQuery(seo, setting, id, rawQuery); err != nil {
					return nil, 500, err
				}
				return seo, 200, nil
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, 500, err
			}
		}
	case strings.HasPrefix(cleanPath, "/content/") && len(strings.Split(strings.Trim(cleanPath, "/"), "/")) == 3:
		// 路径形如 /content/:fid/:id，文章 ID 位于第 3 段
		if id := pathSegmentID(cleanPath, 3); id > 0 {
			seo, err := s.articleSeo(setting, id)
			if err == nil {
				return seo, 200, nil
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, 500, err
			}
		}
	}

	// 登录页、后台与不存在的地址不参与收录
	return s.plainSeo(setting, cleanPath), statusForPath(cleanPath), nil
}

// statusForPath 判断路径应该返回的状态码，SPA 内部页面返回 200，其余返回 404
func statusForPath(path string) int {
	switch {
	case path == "/sign-in", path == "/404", path == "/500",
		path == "/dashboard", path == "/admin", strings.HasPrefix(path, "/admin/"):
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
	seo.Body = s.renderCrawlNav() + s.renderListBody(setting.SiteName, setting.Description, s.recentArticles(0, seoListLimit))
	return seo
}

// renderCrawlNav 输出不依赖 JavaScript 的站内入口，方便爬虫发现关于页、归档和分类
func (s *SeoService) renderCrawlNav() string {
	var b strings.Builder
	b.WriteString(`<nav class="seo-prerender">`)
	b.WriteString(`<a href="/about">关于</a>`)
	b.WriteString(` · <a href="/archive">文章归档</a>`)
	b.WriteString(` · <a href="/rss.xml">RSS</a>`)
	var categories []model.Category
	if err := global.BLOG_DB.Order("sort desc, id asc").Find(&categories).Error; err == nil && len(categories) > 0 {
		b.WriteString("<ul>")
		for _, category := range categories {
			if strings.TrimSpace(category.Name) == "" {
				continue
			}
			b.WriteString(fmt.Sprintf(`<li><a href="/category/%d">%s</a></li>`, category.ID, escapeHTML(category.Name)))
		}
		b.WriteString("</ul>")
	}
	b.WriteString("</nav>")
	return b.String()
}

// renderAboutBody 生成关于页首屏直出内容
func (s *SeoService) renderAboutBody(setting model.Setting, author, bio, email string) string {
	var b strings.Builder
	b.WriteString(`<article class="seo-prerender"><h1>关于 `)
	b.WriteString(escapeHTML(setting.SiteName))
	b.WriteString("</h1>")
	if setting.Description != "" {
		b.WriteString("<p>")
		b.WriteString(escapeHTML(setting.Description))
		b.WriteString("</p>")
	}
	if author != "" {
		b.WriteString("<p>作者：")
		b.WriteString(escapeHTML(author))
		b.WriteString("</p>")
	}
	if bio != "" {
		b.WriteString("<p>")
		b.WriteString(escapeHTML(bio))
		b.WriteString("</p>")
	}
	if email != "" {
		b.WriteString(`<p>联系邮箱：<a href="mailto:`)
		b.WriteString(escapeHTML(email))
		b.WriteString(`">`)
		b.WriteString(escapeHTML(email))
		b.WriteString("</a></p>")
	}
	b.WriteString(`<p><a href="/">返回首页</a> · <a href="/archive">文章归档</a></p></article>`)
	return b.String()
}

// renderArchiveBody 按年份输出全部文章链接
func (s *SeoService) renderArchiveBody(entries []model.ArchiveEntry) string {
	var b strings.Builder
	b.WriteString(`<div class="seo-prerender"><h1>文章归档</h1>`)
	if len(entries) == 0 {
		b.WriteString("<p>暂时还没有文章。</p></div>")
		return b.String()
	}
	year := ""
	opened := false
	for _, item := range entries {
		itemYear := "更早"
		if date := localDate(item.Ctime); len(date) >= 4 {
			itemYear = date[:4]
		}
		if itemYear != year {
			if opened {
				b.WriteString("</ul>")
			}
			year = itemYear
			opened = true
			b.WriteString("<h2>")
			b.WriteString(escapeHTML(year))
			b.WriteString("</h2><ul>")
		}
		b.WriteString(`<li><a href="`)
		b.WriteString(escapeHTML(fmt.Sprintf("/content/%d/%d", item.Fid, item.ID)))
		b.WriteString(`">`)
		b.WriteString(escapeHTML(item.Title))
		b.WriteString("</a>")
		if date := localDate(item.Ctime); date != "" {
			b.WriteString(fmt.Sprintf(` <time datetime="%s">%s</time>`, isoTime(item.Ctime), escapeHTML(date)))
		}
		b.WriteString("</li>")
	}
	if opened {
		b.WriteString("</ul>")
	}
	b.WriteString("</div>")
	return b.String()
}

// aboutJsonLd 关于页结构化数据，把作者实体指向固定地址
func (s *SeoService) aboutJsonLd(setting model.Setting, author, bio, email, image, canonical string) map[string]any {
	person := map[string]any{
		"@type": "Person",
		"name":  author,
		"url":   canonical,
	}
	if bio != "" {
		person["description"] = bio
	}
	if email != "" {
		person["email"] = email
	}
	if image != "" {
		person["image"] = image
	}
	data := map[string]any{
		"@context":   "https://schema.org",
		"@type":      "AboutPage",
		"name":       fmt.Sprintf("关于 %s", setting.SiteName),
		"url":        canonical,
		"mainEntity": person,
	}
	if setting.Description != "" {
		data["description"] = setting.Description
	}
	return data
}

// categorySeo 分类页 SEO 信息
func (s *SeoService) categorySeo(setting model.Setting, fid int) (*PageSeo, error) {
	var category model.Category
	if err := global.BLOG_DB.Where("id = ?", fid).First(&category).Error; err != nil {
		return nil, err
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
	return seo, nil
}

// plainSeo 不参与收录页面的 SEO 信息。
// 缺页不能复用首页标题、描述和规范地址，否则搜索引擎会把 404 当成站点主页。
func (s *SeoService) plainSeo(setting model.Setting, path string) *PageSeo {
	title := setting.SiteName
	description := ""
	if statusForPath(path) == 404 {
		title = "页面不存在"
		if setting.SiteName != "" {
			title += " - " + setting.SiteName
		}
		description = "你访问的页面不存在或已删除。"
	}
	return &PageSeo{
		Title:       title,
		Description: description,
		Robots:      "noindex, nofollow",
		OgType:      "website",
	}
}

// aboutSeo 关于页，承接搜索引擎根据页脚 About 生成的站内链接
func (s *SeoService) aboutSeo(setting model.Setting) *PageSeo {
	author := setting.Author
	bio := ""
	email := ""
	avatar := ""
	users := UserService{}
	if contact, err := users.GetPublicContact(); err == nil {
		if contact.TrueName != "" {
			author = contact.TrueName
		}
		bio = strings.TrimSpace(contact.Bio)
		email = strings.TrimSpace(contact.Email)
		avatar = strings.TrimSpace(contact.Pic)
	}
	if author == "" {
		author = setting.SiteName
	}
	description := bio
	if description == "" {
		description = setting.Description
	}
	title := fmt.Sprintf("关于 - %s", setting.SiteName)
	canonical := s.absoluteURL("/about")
	seo := &PageSeo{
		Title:       title,
		Description: description,
		Keywords:    joinNonEmpty(",", "关于", author, setting.SiteName, setting.Keywords),
		Canonical:   canonical,
		Robots:      "index, follow",
		OgType:      "profile",
		Image:       s.absoluteResource(avatar),
		Author:      author,
	}
	seo.JsonLd = append(seo.JsonLd, s.aboutJsonLd(setting, author, bio, email, seo.Image, canonical))
	seo.Body = s.renderAboutBody(setting, author, bio, email)
	return seo
}

// archiveSeo 归档页，把全部已发布文章以链接形式直出
func (s *SeoService) archiveSeo(setting model.Setting) *PageSeo {
	title := fmt.Sprintf("文章归档 - %s", setting.SiteName)
	description := fmt.Sprintf("%s 的全部文章归档。", setting.SiteName)
	canonical := s.absoluteURL("/archive")
	articles := ArticleService{}
	entries, err := articles.ListArchive(2000)
	if err != nil {
		entries = nil
	}
	seo := &PageSeo{
		Title:       title,
		Description: description,
		Keywords:    joinNonEmpty(",", "归档", setting.SiteName, setting.Keywords),
		Canonical:   canonical,
		Robots:      "index, follow",
		OgType:      "website",
		Author:      setting.Author,
	}
	seo.JsonLd = append(seo.JsonLd, s.collectionJsonLd(title, description, canonical))
	seo.Body = s.renderArchiveBody(entries)
	return seo
}

// absoluteResource 把可公开访问的图片地址转换成绝对地址。
func (s *SeoService) absoluteResource(url string) string {
	return publicMediaURL(url)
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
	if setting.Author != "" {
		publisher := map[string]any{"@type": "Person", "name": setting.Author}
		if authorURL := s.absoluteURL("/about"); authorURL != "" {
			publisher["url"] = authorURL
		}
		data["publisher"] = publisher
	}
	if siteURL != "" {
		data["potentialAction"] = map[string]any{
			"@type":       "SearchAction",
			"target":      siteURL + "/?keyword={search_term_string}",
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

// applyListingQuery 使搜索、分页的直出内容和收录策略与前端保持一致。
func (s *SeoService) applyListingQuery(seo *PageSeo, setting model.Setting, fid int, rawQuery string) error {
	query, _ := url.ParseQuery(rawQuery)
	keyword := strings.TrimSpace(query.Get("keyword"))
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	if page > 1000000 {
		page = 1000000
	}
	var articles []model.Article
	db := global.BLOG_DB.Where("status = ?", 0)
	if fid > 0 {
		db = db.Where("fid = ?", fid)
	}
	if keyword != "" {
		db = db.Where("title LIKE ?", "%"+keyword+"%")
	}
	if err := db.Order("istop desc, ctime desc").Offset((page - 1) * seoListLimit).Limit(seoListLimit).Find(&articles).Error; err != nil {
		return err
	}
	path := "/"
	if fid > 0 {
		path = fmt.Sprintf("/category/%d", fid)
	}
	params := url.Values{}
	if keyword != "" {
		params.Set("keyword", keyword)
		seo.Title = fmt.Sprintf("“%s”的搜索结果 - %s", keyword, setting.SiteName)
		seo.Description = fmt.Sprintf("在 %s 中搜索“%s”的相关文章。", setting.SiteName, keyword)
	}
	if page > 1 {
		params.Set("page", strconv.Itoa(page))
		seo.Title += fmt.Sprintf(" - 第 %d 页", page)
	}
	if keyword != "" || page > 1 {
		seo.Robots = "noindex, follow"
		seo.JsonLd = nil
	}
	seo.Canonical = s.absoluteURL(path)
	if len(params) > 0 {
		seo.Canonical += "?" + params.Encode()
	}
	seo.Body = s.renderCrawlNav() + s.renderListBody(seo.Title, seo.Description, articles)
	return nil
}
