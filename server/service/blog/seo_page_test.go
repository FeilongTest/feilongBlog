package blog

import (
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"strings"
	"testing"
)

func TestPlainSeoDoesNotCanonicalizeMissingPages(t *testing.T) {
	service := SeoService{}
	seo := service.plainSeo(model.Setting{
		SiteName:    "飞龙博客",
		Description: "首页描述",
		Keywords:    "博客",
		Author:      "飞龙",
	}, "/about")

	if seo.Canonical != "" {
		t.Fatalf("missing page canonical = %q", seo.Canonical)
	}
	if seo.Description == "首页描述" || strings.Contains(seo.Description, "首页描述") {
		t.Fatalf("missing page reused homepage description: %q", seo.Description)
	}
	if seo.Title == "飞龙博客" || !strings.Contains(seo.Title, "页面不存在") {
		t.Fatalf("missing page title = %q", seo.Title)
	}
	if seo.Robots != "noindex, nofollow" {
		t.Fatalf("robots = %q", seo.Robots)
	}
}

func TestPlainSeoKeepsPrivatePagesOutOfIndex(t *testing.T) {
	service := SeoService{}
	seo := service.plainSeo(model.Setting{SiteName: "飞龙博客", Description: "首页描述"}, "/sign-in")
	if seo.Canonical != "" {
		t.Fatalf("private page canonical = %q", seo.Canonical)
	}
	if strings.Contains(seo.Title, "页面不存在") {
		t.Fatalf("login page was treated as missing: %q", seo.Title)
	}
	if !strings.Contains(seo.Robots, "noindex") {
		t.Fatalf("robots = %q", seo.Robots)
	}
}

func TestRobotsAllowsSearchEnginesBeforeWildcard(t *testing.T) {
	previous := global.BLOG_CONFIG.Seo.SiteUrl
	global.BLOG_CONFIG.Seo.SiteUrl = "https://teshh.com"
	t.Cleanup(func() { global.BLOG_CONFIG.Seo.SiteUrl = previous })

	service := SeoService{}
	text := service.Robots()
	baidu := strings.Index(text, "User-agent: Baiduspider\nAllow: /")
	wildcard := strings.Index(text, "User-agent: *")
	if baidu < 0 || wildcard < 0 || baidu > wildcard {
		t.Fatalf("search engine allow rules must come before the wildcard:\n%s", text)
	}
	if !strings.Contains(text, "Sitemap: https://teshh.com/sitemap.xml") {
		t.Fatalf("sitemap missing:\n%s", text)
	}
}

func TestCanonicalRedirectNormalizesAboutAndArchive(t *testing.T) {
	service := SeoService{}
	if got := service.CanonicalRedirect("/about/"); got != "/about" {
		t.Fatalf("about redirect = %q", got)
	}
	if got := service.CanonicalRedirect("/ABOUT"); got != "/about" {
		t.Fatalf("ABOUT redirect = %q", got)
	}
	if got := service.CanonicalRedirect("/about"); got != "" {
		t.Fatalf("canonical about should not redirect, got %q", got)
	}
	if got := service.CanonicalRedirect("/archive/"); got != "/archive" {
		t.Fatalf("archive redirect = %q", got)
	}
	if got := service.CanonicalRedirect("/content/1/2"); got != "" {
		t.Fatalf("article path redirected: %q", got)
	}
}
