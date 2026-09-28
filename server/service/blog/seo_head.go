package blog

import (
	"encoding/json"
	"fmt"
	"html"
	"strings"
	"time"
)

// settingService 站点设置服务
var settingService = SettingService{}

// renderHead 生成 head 中的 SEO 标签
func (s *SeoService) renderHead(seo *PageSeo) string {
	var b strings.Builder
	line := func(format string, args ...any) {
		b.WriteString("  ")
		b.WriteString(fmt.Sprintf(format, args...))
		b.WriteString("\n")
	}

	line(`<meta name="description" content="%s">`, html.EscapeString(seo.Description))
	if seo.Keywords != "" {
		line(`<meta name="keywords" content="%s">`, html.EscapeString(seo.Keywords))
	}
	if seo.Author != "" {
		line(`<meta name="author" content="%s">`, html.EscapeString(seo.Author))
	}
	line(`<meta name="robots" content="%s">`, html.EscapeString(seo.Robots))
	if seo.Canonical != "" {
		line(`<link rel="canonical" href="%s">`, html.EscapeString(seo.Canonical))
	}

	setting, err := settingService.GetSetting()
	siteName := "飞龙博客"
	if err == nil {
		siteName = setting.SiteName
	}
	line(`<meta property="og:site_name" content="%s">`, html.EscapeString(siteName))
	line(`<meta property="og:locale" content="zh_CN">`)
	line(`<meta property="og:type" content="%s">`, html.EscapeString(seo.OgType))
	line(`<meta property="og:title" content="%s">`, html.EscapeString(seo.Title))
	line(`<meta property="og:description" content="%s">`, html.EscapeString(seo.Description))
	if seo.Canonical != "" {
		line(`<meta property="og:url" content="%s">`, html.EscapeString(seo.Canonical))
	}
	if seo.Image != "" {
		line(`<meta property="og:image" content="%s">`, html.EscapeString(seo.Image))
		line(`<meta name="twitter:card" content="summary_large_image">`)
		line(`<meta name="twitter:image" content="%s">`, html.EscapeString(seo.Image))
	} else {
		line(`<meta name="twitter:card" content="summary">`)
	}
	line(`<meta name="twitter:title" content="%s">`, html.EscapeString(seo.Title))
	line(`<meta name="twitter:description" content="%s">`, html.EscapeString(seo.Description))
	if seo.Published > 0 {
		line(`<meta property="article:published_time" content="%s">`, isoTime(seo.Published))
	}
	if seo.Modified > 0 {
		line(`<meta property="article:modified_time" content="%s">`, isoTime(seo.Modified))
	}
	if seo.Section != "" {
		line(`<meta property="article:section" content="%s">`, html.EscapeString(seo.Section))
	}

	// 搜索引擎站长平台验证，便于在站长后台完成站点归属校验
	if err == nil {
		if setting.BaiduVerify != "" {
			line(`<meta name="baidu-site-verification" content="%s">`, html.EscapeString(setting.BaiduVerify))
		}
		if setting.GoogleVerify != "" {
			line(`<meta name="google-site-verification" content="%s">`, html.EscapeString(setting.GoogleVerify))
		}
		if setting.BingVerify != "" {
			line(`<meta name="msvalidate.01" content="%s">`, html.EscapeString(setting.BingVerify))
		}
	}

	for _, item := range seo.JsonLd {
		payload, err := json.Marshal(item)
		if err != nil {
			continue
		}
		line(`<script type="application/ld+json">%s</script>`, payload)
	}

	// 统一添加标记，前端接管后可以整体清理
	result := b.String()
	result = strings.ReplaceAll(result, "<meta ", "<meta "+seoTagFlag+" ")
	result = strings.ReplaceAll(result, `<link rel="canonical"`, `<link `+seoTagFlag+` rel="canonical"`)
	result = strings.ReplaceAll(result, `<script type="application/ld+json">`, `<script `+seoTagFlag+` type="application/ld+json">`)
	return result
}

// isoTime 把时间戳格式化为 ISO 8601 字符串
func isoTime(timestamp int) string {
	if timestamp <= 0 {
		return ""
	}
	return time.Unix(int64(timestamp), 0).Local().Format("2006-01-02T15:04:05-07:00")
}
