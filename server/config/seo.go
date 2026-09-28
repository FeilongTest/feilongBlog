package config

// Seo 站点搜索引擎优化相关配置
type Seo struct {
	SiteUrl     string `mapstructure:"site-url" json:"site-url" yaml:"site-url"`             // 站点主域名，例如 https://teshh.com
	WebRoot     string `mapstructure:"web-root" json:"web-root" yaml:"web-root"`             // 前端构建产物目录，配置后开启服务端 SEO 渲染
	SitemapTTL  int    `mapstructure:"sitemap-ttl" json:"sitemap-ttl" yaml:"sitemap-ttl"`    // 站点地图缓存时间（秒），默认 600
	BaiduToken  string `mapstructure:"baidu-token" json:"baidu-token" yaml:"baidu-token"`    // 百度普通收录推送 token
	IndexNowKey string `mapstructure:"indexnow-key" json:"indexnow-key" yaml:"indexnow-key"` // IndexNow 密钥，支持 Bing、Yandex 等
	AutoPush    bool   `mapstructure:"auto-push" json:"auto-push" yaml:"auto-push"`          // 发布或更新文章后自动推送给搜索引擎
}
