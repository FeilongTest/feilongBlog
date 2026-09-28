package blog

import (
	"bytes"
	"html"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
)

// 正文清洗与纯文本提取使用的正则
// 注意：Go 的正则不支持反向引用，这里对标签名做宽松匹配
var (
	seoBlockTagPattern = regexp.MustCompile(`(?is)<(?:script|style|iframe|noscript|form)\b[^>]*>.*?</(?:script|style|iframe|noscript|form)>`)
	seoTagPattern          = regexp.MustCompile(`(?s)<[^>]*>`)
	seoEventAttrPattern    = regexp.MustCompile(`(?i)\son[a-z]+\s*=\s*("[^"]*"|'[^']*'|[^\s>]+)`)
	seoSchemeHrefPattern   = regexp.MustCompile(`(?i)(href|src)\s*=\s*("[^"]*(?:javascript|vbscript):[^"]*"|'[^']*(?:javascript|vbscript):[^']*')`)
	seoSpacePattern        = regexp.MustCompile(`\s+`)
	seoMarkdownCodeFence   = regexp.MustCompile("(?s)```.*?```")
	seoMarkdownImage       = regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)`)
	seoMarkdownLink        = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	seoMarkdownHeading     = regexp.MustCompile(`(?m)^\s{0,3}#{1,6}\s*`)
	seoMarkdownQuoteAndBullet = regexp.MustCompile(`(?m)^\s{0,3}(>|[-*+]\s+|\d+\.\s+)`)
	seoMarkdownTable       = regexp.MustCompile(`(?m)^\s{0,3}\|`)
	seoMarkdownEmphasis    = regexp.MustCompile("[*_`~]+")
)

// seoMarkdown 用于把 Markdown 正文渲染成 HTML，供爬虫直接读取
var seoMarkdown = goldmark.New()

// markdownToHTML 把 Markdown 渲染为 HTML
func markdownToHTML(raw string) string {
	var buf bytes.Buffer
	if err := seoMarkdown.Convert([]byte(raw), &buf); err != nil {
		return ""
	}
	return buf.String()
}

// markdownToText 去掉 Markdown 语法，得到纯文本
func markdownToText(raw string) string {
	text := seoMarkdownCodeFence.ReplaceAllString(raw, " ")
	text = seoMarkdownImage.ReplaceAllString(text, " ")
	text = seoMarkdownLink.ReplaceAllString(text, "$1")
	text = seoMarkdownHeading.ReplaceAllString(text, "")
	text = seoMarkdownQuoteAndBullet.ReplaceAllString(text, "")
	text = seoMarkdownTable.ReplaceAllString(text, "")
	text = seoMarkdownEmphasis.ReplaceAllString(text, "")
	return text
}

// htmlToText 去掉 HTML 标签，得到纯文本
func htmlToText(raw string) string {
	text := seoBlockTagPattern.ReplaceAllString(raw, " ")
	text = seoTagPattern.ReplaceAllString(text, " ")
	return html.UnescapeString(text)
}

// plainText 把正文转换为用于描述与摘要的纯文本
func plainText(content, format string, limit int) string {
	var text string
	if strings.EqualFold(strings.TrimSpace(format), "markdown") {
		text = markdownToText(content)
	} else {
		text = htmlToText(content)
	}
	text = strings.TrimSpace(seoSpacePattern.ReplaceAllString(text, " "))
	return truncateRunes(text, limit)
}

// sanitizeArticleHTML 清理正文中的脚本与事件属性
func sanitizeArticleHTML(raw string) string {
	clean := seoBlockTagPattern.ReplaceAllString(raw, " ")
	clean = seoEventAttrPattern.ReplaceAllString(clean, "")
	clean = seoSchemeHrefPattern.ReplaceAllString(clean, `$1="#"`)
	return strings.TrimSpace(clean)
}

// escapeHTML 转义 HTML 特殊字符
func escapeHTML(text string) string {
	return html.EscapeString(text)
}

// truncateRunes 按字符截断文本
func truncateRunes(text string, limit int) string {
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}
	return string(runes[:limit])
}

// localDate 把时间戳格式化为日期，时间戳非法时返回空字符串
func localDate(timestamp int) string {
	text := isoTime(timestamp)
	if len(text) < 10 {
		return ""
	}
	return text[:10]
}
