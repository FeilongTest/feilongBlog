package blog

import (
	"errors"
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"regexp"
	"strings"
	"testing"
)

func mockDatabase(t *testing.T) sqlmock.Sqlmock {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	orm, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{DisableAutomaticPing: true, Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	prev := global.BLOG_DB
	global.BLOG_DB = orm
	t.Cleanup(func() {
		global.BLOG_DB = prev
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
		db.Close()
	})
	return mock
}

func TestPublicArticleCannotReadHiddenArticle(t *testing.T) {
	mock := mockDatabase(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `blog_article` WHERE id = ? AND status = ? ORDER BY `blog_article`.`id` LIMIT 1")).WithArgs(7, 0).WillReturnRows(sqlmock.NewRows([]string{"id", "status"}))
	_, err := (&ArticleService{}).GetArticle(7)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("error=%v", err)
	}
}

func TestAdminCanEditHiddenArticleWithoutIncrementingViews(t *testing.T) {
	mock := mockDatabase(t)
	mock.ExpectQuery("SELECT .* WHERE .*id.*").WithArgs(7).WillReturnRows(sqlmock.NewRows([]string{"id", "status", "view"}).AddRow(7, 1, 9))
	article, err := (&ArticleService{}).GetArticleAdmin(7)
	if err != nil || article.Status != 1 || article.View != 9 {
		t.Fatalf("article=%+v error=%v", article, err)
	}
}

func TestArticleSanitizationRejectsExecutableMarkup(t *testing.T) {
	for _, raw := range []string{
		`<a href="java&#x73;cript:alert(1)">link</a>`,
		`<img src="x" onerror="alert(1)">`,
		`<svg><a xlink:href="javascript:alert(1)">x</a></svg>`,
		`<iframe srcdoc="<script>alert(1)</script>"></iframe>`,
		`<script>alert(1)`,
		`<a href=javascript:alert(1)>x</a>`,
	} {
		clean := sanitizeArticleHTML(raw)
		if strings.Contains(clean, "alert(1)") || strings.Contains(clean, "onerror") || strings.Contains(clean, "<script") || strings.Contains(clean, "<iframe") {
			t.Fatalf("unsafe output for %q: %q", raw, clean)
		}
	}
	safe := `<h2 id="section-1">标题</h2><p><a href="https://example.com">链接</a></p><pre><code class="language-go">fmt.Println(1)</code></pre><img src="https://example.com/x.png">`
	clean := sanitizeArticleHTML(safe)
	for _, want := range []string{`id="section-1"`, `language-go`, `https://example.com/x.png`, `<pre>`} {
		if !strings.Contains(clean, want) {
			t.Fatalf("safe content lost: %s", clean)
		}
	}
}

func TestSeoDatabaseErrorsAreNot404(t *testing.T) {
	mock := mockDatabase(t)
	mock.ExpectQuery("SELECT .* WHERE id = ").WithArgs(1).WillReturnError(errors.New("database unavailable"))
	_, status, err := (&SeoService{}).buildPageSeo("/content/1/7", "")
	if status != 500 || err == nil {
		t.Fatalf("status=%d error=%v", status, err)
	}
}

func TestSearchAndPaginationSeo(t *testing.T) {
	mock := mockDatabase(t)
	previous := global.BLOG_CONFIG.Seo.SiteUrl
	global.BLOG_CONFIG.Seo.SiteUrl = "https://example.com"
	t.Cleanup(func() { global.BLOG_CONFIG.Seo.SiteUrl = previous })
	mock.ExpectQuery("SELECT .* WHERE status = .*title LIKE.*LIMIT 10 OFFSET 10").WithArgs(0, "%Go%").WillReturnRows(sqlmock.NewRows([]string{"id", "title", "fid", "status"}).AddRow(2, "Go page two", 3, 0))
	mock.ExpectQuery("SELECT .* ORDER BY sort desc, id asc").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))
	seo := &PageSeo{Title: "Blog", Robots: "index, follow"}
	err := (&SeoService{}).applyListingQuery(seo, model.Setting{SiteName: "Blog"}, 0, "keyword=Go&page=2")
	if err != nil {
		t.Fatal(err)
	}
	if seo.Robots != "noindex, follow" || seo.Canonical != "https://example.com/?keyword=Go&page=2" || !strings.Contains(seo.Body, "Go page two") || !strings.Contains(seo.Title, "第 2 页") {
		t.Fatalf("seo=%+v", seo)
	}
}

func TestArticleSeoDistinguishesMissingFromDatabaseFailure(t *testing.T) {
	for _, test := range []struct {
		name    string
		failure error
		status  int
	}{{"missing", gorm.ErrRecordNotFound, 404}, {"database failure", errors.New("connection lost"), 500}} {
		t.Run(test.name, func(t *testing.T) {
			mock := mockDatabase(t)
			mock.ExpectQuery("SELECT .* WHERE id = ").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "site_name"}).AddRow(1, "Blog"))
			mock.ExpectQuery("SELECT .* WHERE id = .*status = ").WithArgs(7, 0).WillReturnError(test.failure)
			_, status, err := (&SeoService{}).buildPageSeo("/content/1/7", "")
			if status != test.status {
				t.Fatalf("status=%d error=%v", status, err)
			}
			if test.status == 500 && err == nil {
				t.Fatal("database error was swallowed")
			}
		})
	}
}

func TestHiddenArticlesCannotBeLiked(t *testing.T) {
	mock := mockDatabase(t)
	mock.ExpectQuery("SELECT .* FROM `blog_article` WHERE id = .*status = ").WithArgs(7, 0).WillReturnRows(sqlmock.NewRows([]string{"id"}))
	_, _, err := (&ArticleService{}).LikeArticle(model.ArticleLikeRequest{Aid: 7, VisitorID: "visitor"})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("error=%v", err)
	}
}

func TestEmbeddedLocalMediaIsRootAbsolute(t *testing.T) {
	got := sanitizeArticleHTML(`<p><img src="public/uploads/file/2023/image.png"></p>`)
	if !strings.Contains(got, `src="/blog/uploads/file/2023/image.png"`) {
		t.Fatalf("unexpected media path: %s", got)
	}
}

func TestCreateArticleKeepsAssignedIDForPublishResponse(t *testing.T) {
	mock := mockDatabase(t)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `blog_article`").WillReturnResult(sqlmock.NewResult(456, 1))
	mock.ExpectCommit()
	article := model.Article{Title: "新文章", Fid: 1, ContentFormat: "html"}
	if err := (&ArticleService{}).CreateArticle(&article); err != nil {
		t.Fatal(err)
	}
	if article.ID != 456 {
		t.Fatalf("created ID=%d", article.ID)
	}
}

func TestPublicFeedUsesSummaryWithoutDroppingTextUpdates(t *testing.T) {
	mock := mockDatabase(t)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `blog_article`").WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	mock.ExpectQuery("SELECT \\* FROM `blog_article`.*LIMIT 10").WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"id", "type", "content", "content_format"}).AddRow(1, 1, "# 技术文章\n\n正文说明", "markdown").AddRow(2, 0, "<p>这是一条动态</p>", "html"))
	mock.ExpectQuery("SELECT aid, count\\(\\*\\) as count FROM `blog_comment`").WillReturnRows(sqlmock.NewRows([]string{"aid", "count"}))
	mock.ExpectQuery("SELECT aid, count\\(\\*\\) as count FROM `blog_article_like`").WillReturnRows(sqlmock.NewRows([]string{"aid", "count"}))
	info := model.ArticleSearch{}
	info.Page = 1
	info.PageSize = 10
	articles, total, err := (&ArticleService{}).GetArticleList(info)
	if err != nil {
		t.Fatal(err)
	}
	if total != 2 || len(articles) != 2 {
		t.Fatalf("count=%d articles=%d", total, len(articles))
	}
	if articles[0].Content != "" || !strings.Contains(articles[0].Summary, "技术文章") || strings.Contains(articles[0].Summary, "#") {
		t.Fatalf("summary=%+v", articles[0])
	}
	if articles[1].Content != "<p>这是一条动态</p>" {
		t.Fatal("text update lost its body")
	}
}
