package blog

import (
	"bytes"
	htmlparser "golang.org/x/net/html"
	"html"
	"io"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
)

// 正文清洗与纯文本提取使用的正则
// 注意：Go 的正则不支持反向引用，这里对标签名做宽松匹配
var (
	seoBlockTagPattern        = regexp.MustCompile(`(?is)<(?:script|style|iframe|noscript|form)\b[^>]*>.*?</(?:script|style|iframe|noscript|form)>`)
	seoTagPattern             = regexp.MustCompile(`(?s)<[^>]*>`)
	seoSpacePattern           = regexp.MustCompile(`\s+`)
	seoMarkdownCodeFence      = regexp.MustCompile("(?s)```.*?```")
	seoMarkdownImage          = regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)`)
	seoMarkdownLink           = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	seoMarkdownHeading        = regexp.MustCompile(`(?m)^\s{0,3}#{1,6}\s*`)
	seoMarkdownQuoteAndBullet = regexp.MustCompile(`(?m)^\s{0,3}(>|[-*+]\s+|\d+\.\s+)`)
	seoMarkdownTable          = regexp.MustCompile(`(?m)^\s{0,3}\|`)
	seoMarkdownEmphasis       = regexp.MustCompile("[*_`~]+")
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

// articleHTMLPolicy 使用 HTML 解析器和允许列表，实体解码后的 URL 也必须是安全协议。
var articleHTMLPolicy = func() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("id").Matching(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)).OnElements("h1", "h2", "h3", "h4", "h5", "h6")
	p.AllowAttrs("class").Matching(regexp.MustCompile(`^language-[a-zA-Z0-9_-]+$`)).OnElements("code", "pre")
	return p
}()

func sanitizeArticleHTML(raw string) string {
	return normalizeArticleMedia(strings.TrimSpace(articleHTMLPolicy.Sanitize(raw)))
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

// normalizeArticleMedia 修正正文中旧的相对图片，避免在 /content/... 下解析出错误地址。
func normalizeArticleMedia(content string) string {
	tokenizer := htmlparser.NewTokenizer(strings.NewReader(content))
	var output strings.Builder
	for {
		kind := tokenizer.Next()
		if kind == htmlparser.ErrorToken {
			if tokenizer.Err() == io.EOF {
				return output.String()
			}
			return content
		}
		if kind == htmlparser.StartTagToken || kind == htmlparser.SelfClosingTagToken {
			token := tokenizer.Token()
			if token.Data == "img" {
				for i := range token.Attr {
					if token.Attr[i].Key == "src" {
						token.Attr[i].Val = publicMediaURL(token.Attr[i].Val)
					}
				}
				output.WriteString(token.String())
				continue
			}
		}
		output.Write(tokenizer.Raw())
	}
}
