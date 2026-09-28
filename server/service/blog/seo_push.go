package blog

import (
	"bytes"
	"encoding/json"
	"feilongBlog/global"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

// 推送搜索引擎时使用的客户端
var seoHTTPClient = &http.Client{Timeout: 15 * time.Second}

// PushArticles 异步把文章地址推送给搜索引擎，未配置时自动跳过
func (s *SeoService) PushArticles(urls []string) {
	urls = normalizeURLs(urls)
	if len(urls) == 0 || !global.BLOG_CONFIG.Seo.AutoPush {
		return
	}
	if global.BLOG_CONFIG.Seo.BaiduToken == "" && global.BLOG_CONFIG.Seo.IndexNowKey == "" {
		return
	}
	go func() {
		s.pushBaidu(urls)
		s.pushIndexNow(urls)
	}()
}

// ArticleURL 拼接文章的前台地址
func (s *SeoService) ArticleURL(fid int, id uint) string {
	return s.absoluteURL(fmt.Sprintf("/content/%d/%d", fid, id))
}

// notifySearchEngines 内容变化后清理站点地图缓存并推送地址给搜索引擎
func notifySearchEngines(urls ...string) {
	service := SeoService{}
	service.ResetSitemap()
	service.PushArticles(urls)
}

// IndexNowKey 返回当前的 IndexNow 密钥
func (s *SeoService) IndexNowKey() string {
	return strings.TrimSpace(global.BLOG_CONFIG.Seo.IndexNowKey)
}

// normalizeURLs 过滤空地址并去重
func normalizeURLs(urls []string) []string {
	seen := make(map[string]struct{}, len(urls))
	result := make([]string, 0, len(urls))
	for _, item := range urls {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}

// pushBaidu 使用百度普通收录接口推送地址
func (s *SeoService) pushBaidu(urls []string) {
	siteURL := s.SiteURL()
	token := strings.TrimSpace(global.BLOG_CONFIG.Seo.BaiduToken)
	if siteURL == "" || token == "" {
		return
	}
	endpoint := fmt.Sprintf("http://data.zz.baidu.com/urls?site=%s&token=%s", url.QueryEscape(siteURL), url.QueryEscape(token))
	body := strings.Join(urls, "\n")
	resp, err := seoHTTPClient.Post(endpoint, "text/plain", bytes.NewBufferString(body))
	if err != nil {
		global.BLOG_LOG.Warn("推送百度收录失败", zap.Error(err))
		return
	}
	defer resp.Body.Close()
	result, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	global.BLOG_LOG.Info("推送百度收录完成", zap.Int("status", resp.StatusCode), zap.String("result", string(result)))
}

// pushIndexNow 使用 IndexNow 推送给 Bing、Yandex 等搜索引擎
func (s *SeoService) pushIndexNow(urls []string) {
	siteURL := s.SiteURL()
	key := s.IndexNowKey()
	if siteURL == "" || key == "" {
		return
	}
	host := strings.TrimPrefix(strings.TrimPrefix(siteURL, "https://"), "http://")
	payload := map[string]any{
		"host":        host,
		"key":         key,
		"keyLocation": fmt.Sprintf("%s/%s.txt", siteURL, key),
		"urlList":     urls,
	}
	content, err := json.Marshal(payload)
	if err != nil {
		return
	}
	resp, err := seoHTTPClient.Post("https://api.indexnow.org/indexnow", "application/json", bytes.NewReader(content))
	if err != nil {
		global.BLOG_LOG.Warn("推送 IndexNow 失败", zap.Error(err))
		return
	}
	defer resp.Body.Close()
	result, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	global.BLOG_LOG.Info("推送 IndexNow 完成", zap.Int("status", resp.StatusCode), zap.String("result", string(result)))
}
