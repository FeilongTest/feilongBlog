package blog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"feilongBlog/global"
	"feilongBlog/initialize"
	model "feilongBlog/model/blog"
	service "feilongBlog/service/blog"
	"feilongBlog/utils"
	driver "github.com/go-sql-driver/mysql"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Opt-in real MySQL/MariaDB contract tests. Only isolated databases with this
// prefix are accepted; the fixture tables are reset at the start of each test.
func publishingDatabase(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("BLOG_TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("set BLOG_TEST_MYSQL_DSN to run real database contracts")
	}
	parsed, err := driver.ParseDSN(dsn)
	if err != nil || !strings.HasPrefix(parsed.DBName, "feilong_integration_") {
		t.Fatal("integration DSN must name an isolated feilong_integration_ database")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, Logger: logger.Default.LogMode(logger.Error)})
	if err != nil {
		t.Fatal(err)
	}
	previousDB, previousConfig, previousLog := global.BLOG_DB, global.BLOG_CONFIG, global.BLOG_LOG
	global.BLOG_DB = db
	global.BLOG_LOG = zap.NewNop()
	global.BLOG_CONFIG.Seo.AutoPush = false
	global.BLOG_CONFIG.JWT.SigningKey = "integration-test-key-never-use-in-production"
	global.BLOG_CONFIG.Local.Path = "/uploads/file"
	global.BLOG_CONFIG.Local.StorePath = t.TempDir()
	global.BLOG_CONFIG.Seo.SiteUrl = "https://example.com"
	global.BLOG_CONFIG.Seo.WebRoot = t.TempDir()
	if err := os.WriteFile(filepath.Join(global.BLOG_CONFIG.Seo.WebRoot, "index.html"), []byte(`<html><head><title>Test</title><!--seo:head--></head><body><div id="app"></div></body></html>`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.Article{}, &model.Category{}, &model.Tag{}, &model.ArticleTag{}, &model.ArticleRevision{}, &model.ArticleLike{}, &model.Comment{}, &model.Setting{}); err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"blog_article_tag", "blog_article_revision", "blog_tag", "blog_article_like", "blog_comment", "blog_article", "blog_category", "blog_setting"} {
		if err := db.Exec("DELETE FROM " + table).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Create(&model.Category{ID: 1, Name: "技术"}).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		global.BLOG_DB = previousDB
		global.BLOG_CONFIG = previousConfig
		global.BLOG_LOG = previousLog
		raw, _ := db.DB()
		raw.Close()
	})
	return db
}

func sampleArticle(status int) model.Article {
	return model.Article{Title: "从草稿到发布", Content: "# 原始正文\n\n独立的正文关键词", ContentFormat: "markdown", Fid: 1, Type: 1, Status: status, TagNames: []string{"Go", "系统", "go"}, Excerpt: "作者自定义摘要", SeoTitle: "面向搜索的标题", SeoDescription: "独立的搜索描述"}
}

func TestPublishingDraftsHistoryAndConflicts(t *testing.T) {
	db := publishingDatabase(t)
	articles := service.ArticleService{}
	tags := service.TagService{}
	article := sampleArticle(model.ArticleDraft)
	if err := articles.CreateArticle(&article); err != nil {
		t.Fatal(err)
	}
	if article.ID == 0 || article.Version != 1 || len(article.Tags) != 2 {
		t.Fatalf("created article=%+v", article)
	}
	if _, err := articles.GetArticle(int(article.ID)); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("draft public read=%v", err)
	}
	publicTags, err := tags.List(true)
	if err != nil || len(publicTags) != 0 {
		t.Fatalf("draft tags leaked: %+v %v", publicTags, err)
	}
	html, status, err := (&service.SeoService{}).RenderPage(fmt.Sprintf("/content/1/%d", article.ID), "")
	if err != nil || status != 404 || strings.Contains(html, "独立的正文关键词") {
		t.Fatalf("draft SSR status=%d error=%v", status, err)
	}
	stale := article
	article.Status = model.ArticlePublished
	article.Title = "已发布新标题"
	if err := articles.UpdateArticle(&article); err != nil {
		t.Fatal(err)
	}
	if article.Version != 2 {
		t.Fatalf("version=%d", article.Version)
	}
	stale.Title = "旧页面覆盖"
	if err := articles.UpdateArticle(&stale); !errors.Is(err, service.ErrArticleConflict) {
		t.Fatalf("stale save=%v", err)
	}
	public, err := articles.GetArticle(int(article.ID))
	if err != nil || public.Title != "已发布新标题" {
		t.Fatalf("public=%+v %v", public, err)
	}
	revisions, err := articles.ListRevisions(article.ID)
	if err != nil || len(revisions) != 2 {
		t.Fatalf("revisions=%+v %v", revisions, err)
	}
	snapshot, err := articles.GetRevision(article.ID, revisions[1].ID)
	if err != nil || snapshot.Title != "从草稿到发布" {
		t.Fatalf("snapshot=%+v %v", snapshot, err)
	}
	if _, err := articles.GetRevision(article.ID+1, revisions[1].ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("cross-article history read=%v", err)
	}
	top := 1
	patched, err := articles.PatchArticleState(model.ArticleStatePatch{ID: article.ID, Version: article.Version, IsTop: &top})
	if err != nil || patched.Content != article.Content || patched.View != public.View || patched.Version != 3 {
		t.Fatalf("state patch=%+v %v", patched, err)
	}
	var count int64
	db.Model(&model.Article{}).Where("id = ? AND title = ?", article.ID, "已发布新标题").Count(&count)
	if count != 1 {
		t.Fatal("stale update changed persisted article")
	}
}

func TestPublishingScheduleIsPrivateAtomicAndRestartSafe(t *testing.T) {
	publishingDatabase(t)
	articles := service.ArticleService{}
	now := time.Now().Unix()
	article := sampleArticle(model.ArticleScheduled)
	article.PublishAt = now + 120
	if err := articles.CreateArticle(&article); err != nil {
		t.Fatal(err)
	}
	if err := articles.PublishedArticleExists(int(article.ID)); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatal("scheduled article exposed before publish")
	}
	if count, err := articles.PublishDue(context.Background(), now+119); err != nil || count != 0 {
		t.Fatalf("early publish=%d %v", count, err)
	}
	var wg sync.WaitGroup
	results := make(chan int, 2)
	failures := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count, err := articles.PublishDue(context.Background(), now+121)
			results <- count
			failures <- err
		}()
	}
	wg.Wait()
	close(results)
	close(failures)
	total := 0
	for count := range results {
		total += count
	}
	for err := range failures {
		if err != nil {
			t.Fatal(err)
		}
	}
	if total != 1 {
		t.Fatalf("published %d times", total)
	}
	public, err := articles.GetArticle(int(article.ID))
	if err != nil || public.Status != model.ArticlePublished || public.Ctime != int(article.PublishAt) || public.Version != 2 {
		t.Fatalf("published=%+v %v", public, err)
	}
	if count, err := articles.PublishDue(context.Background(), now+300); err != nil || count != 0 {
		t.Fatalf("retry publish=%d %v", count, err)
	}
}

func TestPublishingTagsMetadataAndLiteralSearch(t *testing.T) {
	publishingDatabase(t)
	articles := service.ArticleService{}
	tags := service.TagService{}
	article := sampleArticle(model.ArticlePublished)
	article.Title = "100% Go"
	if err := articles.CreateArticle(&article); err != nil {
		t.Fatal(err)
	}
	other := sampleArticle(model.ArticlePublished)
	other.Title = "100XX Go"
	other.Content = "完全不同的内容"
	other.Excerpt = "另一个摘要"
	other.TagNames = []string{"Frontend"}
	if err := articles.CreateArticle(&other); err != nil {
		t.Fatal(err)
	}
	info := model.ArticleSearch{}
	info.Page = 1
	info.PageSize = 10
	info.Keyword = "%"
	list, total, err := articles.GetArticleList(info)
	if err != nil || total != 1 || list[0].ID != article.ID {
		t.Fatalf("literal wildcard search=%+v total=%d %v", list, total, err)
	}
	info.Keyword = "正文关键词"
	info.Scope = "all"
	list, total, err = articles.GetArticleList(info)
	if err != nil || total != 1 || list[0].Summary != "作者自定义摘要" || list[0].Content != "" {
		t.Fatalf("body search=%+v total=%d %v", list, total, err)
	}
	info.Keyword = ""
	info.TagID = article.Tags[0].ID
	_, total, err = articles.GetArticleList(info)
	if err != nil || total != 1 {
		t.Fatalf("tag filter total=%d %v", total, err)
	}
	directory, err := tags.List(true)
	if err != nil || len(directory) != 3 {
		t.Fatalf("tag directory=%+v %v", directory, err)
	}
	tag, err := tags.GetPublic(article.Tags[0].ID)
	if err != nil || tag.Count != 1 {
		t.Fatalf("public tag=%+v %v", tag, err)
	}
	if err := tags.Delete(tag.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := articles.GetArticle(int(article.ID)); err != nil {
		t.Fatalf("deleting tag deleted article: %v", err)
	}
}

func TestPublishingRevisionRetention(t *testing.T) {
	db := publishingDatabase(t)
	articles := service.ArticleService{}
	article := sampleArticle(model.ArticleDraft)
	if err := articles.CreateArticle(&article); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 55; i++ {
		article.Title = fmt.Sprintf("版本 %d", i)
		if err := articles.UpdateArticle(&article); err != nil {
			t.Fatal(err)
		}
	}
	var count int64
	db.Model(&model.ArticleRevision{}).Where("article_id = ?", article.ID).Count(&count)
	if count != 50 {
		t.Fatalf("revision count=%d", count)
	}
}

func TestPublishingMigrationPreservesLegacyContent(t *testing.T) {
	db := publishingDatabase(t)
	if err := db.Migrator().DropTable(&model.Article{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("CREATE TABLE blog_article (id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT, title VARCHAR(40) NOT NULL, content TEXT NOT NULL, ctime INT NOT NULL DEFAULT 0, edittime INT NOT NULL DEFAULT 0, status INT NOT NULL DEFAULT 0)").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("INSERT INTO blog_article (title,content) VALUES (?,?)", "旧站标题", "旧站正文").Error; err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := initialize.MigratePublishing(db); err != nil {
			t.Fatal(err)
		}
	}
	var record struct {
		Title   string
		Content string
		Version uint64
		Excerpt string
	}
	if err := db.Table("blog_article").First(&record).Error; err != nil {
		t.Fatal(err)
	}
	if record.Title != "旧站标题" || record.Content != "旧站正文" || record.Version != 1 {
		t.Fatalf("legacy data=%+v", record)
	}
	if err := db.Table("blog_article").Where("id > 0").Updates(map[string]interface{}{"ctime": int64(2200000000), "edittime": int64(2200000001)}).Error; err != nil {
		t.Fatalf("time columns were not widened beyond 2038: %v", err)
	}
	// Restore full fixture schema for any subsequent contracts.
	if err := db.AutoMigrate(&model.Article{}); err != nil {
		t.Fatal(err)
	}
}

func TestPublishingHTTPContracts(t *testing.T) {
	publishingDatabase(t)
	router := initialize.Routers()
	token, err := utils.NewJWT().CreateToken(model.CustomClaims{BaseClaims: model.BaseClaims{ID: 1, Admin: 1}, StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour).Unix()}})
	if err != nil {
		t.Fatal(err)
	}
	request := func(method, path string, body interface{}, authenticated bool) (int, map[string]interface{}) {
		raw, _ := json.Marshal(body)
		req := httptest.NewRequest(method, path, bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		if authenticated {
			req.Header.Set("x-token", token)
		}
		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)
		var payload map[string]interface{}
		if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
			t.Fatalf("response %s", response.Body.String())
		}
		return response.Code, payload
	}
	status, _ := request("GET", "/admin/article/revisions?id=1", nil, false)
	if status != 401 {
		t.Fatalf("unprotected revisions=%d", status)
	}
	status, payload := request("POST", "/admin/article/createArticle", sampleArticle(model.ArticleDraft), true)
	if status != 200 || payload["code"] != float64(0) {
		t.Fatalf("create=%d %+v", status, payload)
	}
	saved := payload["data"].(map[string]interface{})
	id := uint(saved["id"].(float64))
	version := uint64(saved["version"].(float64))
	if id == 0 || version != 1 {
		t.Fatalf("create contract=%+v", payload)
	}
	status, payload = request("POST", "/base/getArticle", map[string]interface{}{"id": id}, false)
	if payload["code"] == float64(0) {
		t.Fatal("draft exposed through public HTTP")
	}
	status, payload = request("PATCH", "/admin/article/state", map[string]interface{}{"id": id, "version": version, "status": model.ArticlePublished}, true)
	if status != 200 || payload["code"] != float64(0) {
		t.Fatalf("publish=%d %+v", status, payload)
	}
	status, _ = request("PATCH", "/admin/article/state", map[string]interface{}{"id": id, "version": version, "status": model.ArticleHidden}, true)
	if status != 409 {
		t.Fatalf("stale patch=%d", status)
	}
}

func TestPublishingTagRenameInvalidatesStaleEditors(t *testing.T) {
	publishingDatabase(t)
	articles := service.ArticleService{}
	tags := service.TagService{}
	article := sampleArticle(model.ArticleDraft)
	if err := articles.CreateArticle(&article); err != nil {
		t.Fatal(err)
	}
	stale := article
	tag := article.Tags[0]
	tag.Name = "Golang"
	if err := tags.Save(&tag); err != nil {
		t.Fatal(err)
	}
	if err := articles.UpdateArticle(&stale); !errors.Is(err, service.ErrArticleConflict) {
		t.Fatalf("old editor restored renamed tags: %v", err)
	}
	updated, err := articles.GetArticleAdmin(int(article.ID))
	if err != nil {
		t.Fatal(err)
	}
	if updated.Version != article.Version+1 || !strings.Contains(strings.Join(updated.TagNames, ","), "Golang") {
		t.Fatalf("taxonomy update=%+v", updated)
	}
}

func TestPublishingDeleteCleansPrivateHistoryAndReadingList(t *testing.T) {
	db := publishingDatabase(t)
	articles := service.ArticleService{}
	article := sampleArticle(model.ArticlePublished)
	if err := articles.CreateArticle(&article); err != nil {
		t.Fatal(err)
	}
	if err := articles.DeleteArticle(article); err != nil {
		t.Fatal(err)
	}
	for _, row := range []struct {
		model interface{}
		field string
	}{{&model.ArticleTag{}, "article_id"}, {&model.ArticleRevision{}, "article_id"}, {&model.ArticleLike{}, "aid"}, {&model.Comment{}, "aid"}} {
		var count int64
		db.Model(row.model).Where(row.field+" = ?", article.ID).Count(&count)
		if count != 0 {
			t.Fatalf("orphaned %s records", row.field)
		}
	}
	list, err := articles.ReadingArticles([]uint{article.ID})
	if err != nil || len(list) != 0 {
		t.Fatalf("deleted reading item=%+v %v", list, err)
	}
}

func TestPublishingMetadataRenderedInPublicHTML(t *testing.T) {
	publishingDatabase(t)
	articles := service.ArticleService{}
	article := sampleArticle(model.ArticlePublished)
	if err := articles.CreateArticle(&article); err != nil {
		t.Fatal(err)
	}
	html, status, err := (&service.SeoService{}).RenderPage(fmt.Sprintf("/content/1/%d", article.ID), "")
	if err != nil || status != 200 || !strings.Contains(html, "<title>面向搜索的标题</title>") || !strings.Contains(html, `name="description" content="独立的搜索描述"`) {
		t.Fatalf("article SEO status=%d error=%v html=%s", status, err, html)
	}
	topic, status, err := (&service.SeoService{}).RenderPage(fmt.Sprintf("/tag/%d", article.Tags[0].ID), "")
	if err != nil || status != 200 || !strings.Contains(topic, article.Title) || !strings.Contains(topic, article.Excerpt) {
		t.Fatalf("topic status=%d error=%v", status, err)
	}
	siteMap := (&service.SeoService{}).Sitemap()
	if !strings.Contains(siteMap, fmt.Sprintf("/tag/%d", article.Tags[0].ID)) {
		t.Fatal("public topic missing from sitemap")
	}
}

func TestPublishingValidationDoesNotPersistInvalidContent(t *testing.T) {
	db := publishingDatabase(t)
	articles := service.ArticleService{}
	for _, edit := range []func(*model.Article){func(a *model.Article) { a.Status = 99 }, func(a *model.Article) { a.Title = strings.Repeat("字", 201) }, func(a *model.Article) { a.ContentFormat = "unknown" }, func(a *model.Article) { a.Fid = 999 }, func(a *model.Article) { a.Status = model.ArticleScheduled; a.PublishAt = time.Now().Unix() - 1 }, func(a *model.Article) { a.TagNames = []string{strings.Repeat("标", 33)} }} {
		article := sampleArticle(model.ArticlePublished)
		edit(&article)
		if err := articles.CreateArticle(&article); err == nil {
			t.Fatalf("invalid article accepted=%+v", article)
		}
	}
	var count int64
	db.Model(&model.Article{}).Count(&count)
	if count != 0 {
		t.Fatalf("invalid writes=%d", count)
	}
	draft := sampleArticle(model.ArticleDraft)
	draft.Title = ""
	draft.Fid = 0
	if err := articles.CreateArticle(&draft); err != nil || draft.Title != "未命名草稿" {
		t.Fatalf("untitled draft=%+v %v", draft, err)
	}
}
