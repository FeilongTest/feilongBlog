package blog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ValidationError 可直接反馈给编辑者；数据库和网络错误不向公开接口暴露。
type ValidationError struct{ Message string }

func (e *ValidationError) Error() string  { return e.Message }
func invalidArticle(message string) error { return &ValidationError{Message: message} }

var ErrArticleConflict = errors.New("文章已在其他页面更新，请重新载入后再保存；本机草稿仍然保留")

const maxRevisions = 50

func validateArticle(article *model.Article, now int64) error {
	article.Title = strings.TrimSpace(article.Title)
	article.Excerpt = strings.TrimSpace(article.Excerpt)
	article.SeoTitle = strings.TrimSpace(article.SeoTitle)
	article.SeoDescription = strings.TrimSpace(article.SeoDescription)
	if article.Title == "" && article.Status == model.ArticleDraft {
		article.Title = "未命名草稿"
	}
	if article.Title == "" || utf8.RuneCountInString(article.Title) > 200 {
		return invalidArticle("标题需要 1 到 200 个字符")
	}
	if len(article.Content) > 2<<20 {
		return invalidArticle("正文不能超过 2 MB")
	}
	if utf8.RuneCountInString(article.Excerpt) > 300 || utf8.RuneCountInString(article.SeoDescription) > 300 || utf8.RuneCountInString(article.SeoTitle) > 200 {
		return invalidArticle("摘要和 SEO 描述最多 300 字，SEO 标题最多 200 字")
	}
	if article.ContentFormat == "" {
		article.ContentFormat = "html"
	}
	if article.ContentFormat != "html" && article.ContentFormat != "markdown" {
		return invalidArticle("不支持的正文格式")
	}
	if article.Type != 0 && article.Type != 1 {
		return invalidArticle("文章样式不合法")
	}
	if article.Status < model.ArticlePublished || article.Status > model.ArticleScheduled {
		return invalidArticle("文章状态不合法")
	}
	if article.IsTop != 0 && article.IsTop != 1 {
		return invalidArticle("置顶参数不合法")
	}
	if article.Status == model.ArticleScheduled && article.PublishAt <= now {
		return invalidArticle("定时发布时间必须晚于当前时间")
	}
	if article.Status != model.ArticleScheduled {
		article.PublishAt = 0
	}
	if article.Fid < 0 || (article.Fid == 0 && article.Status != model.ArticleDraft) {
		return invalidArticle("发布前请选择文章分类")
	}
	if len(article.Pic) > 500 {
		return invalidArticle("缩略图地址过长")
	}
	names, err := normalizeTagNames(article.TagNames)
	if err != nil {
		return err
	}
	article.TagNames = names
	return nil
}

func normalizeTagNames(names []string) ([]string, error) {
	if names == nil {
		return nil, nil
	}
	result := make([]string, 0, len(names))
	seen := map[string]bool{}
	for _, name := range names {
		name = strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(name), "#"))
		if name == "" {
			continue
		}
		if utf8.RuneCountInString(name) > 32 || strings.IndexFunc(name, unicode.IsControl) >= 0 {
			return nil, invalidArticle("标签最多 32 字且不能包含控制字符")
		}
		key := strings.ToLower(name)
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, name)
	}
	if len(result) > 10 {
		return nil, invalidArticle("每篇文章最多添加 10 个标签")
	}
	return result, nil
}

func validateArticleCategory(tx *gorm.DB, article model.Article) error {
	if article.Fid == 0 {
		return nil
	}
	var category model.Category
	if err := tx.Select("id").First(&category, article.Fid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invalidArticle("所选分类不存在")
		}
		return err
	}
	return nil
}

func attachArticleTags(tx *gorm.DB, articles ...*model.Article) error {
	if len(articles) == 0 {
		return nil
	}
	ids := make([]uint, 0, len(articles))
	byID := map[uint]*model.Article{}
	for _, a := range articles {
		a.Tags = []model.Tag{}
		a.TagNames = []string{}
		ids = append(ids, a.ID)
		byID[a.ID] = a
	}
	var rows []struct {
		ArticleID   uint
		ID          uint
		Name        string
		Description string
	}
	if err := tx.Table("blog_article_tag AS link").Select("link.article_id, tag.id, tag.name, tag.description").Joins("JOIN blog_tag AS tag ON tag.id = link.tag_id").Where("link.article_id IN ?", ids).Order("tag.name asc, tag.id asc").Scan(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		if a := byID[row.ArticleID]; a != nil {
			a.Tags = append(a.Tags, model.Tag{ID: row.ID, Name: row.Name, Description: row.Description})
			a.TagNames = append(a.TagNames, row.Name)
		}
	}
	return nil
}

func replaceArticleTags(tx *gorm.DB, id uint, names []string) error {
	if names == nil {
		return nil
	}
	if err := tx.Where("article_id = ?", id).Delete(&model.ArticleTag{}).Error; err != nil {
		return err
	}
	seenIDs := map[uint]bool{}
	for _, name := range names {
		tag := model.Tag{Name: name}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&tag).Error; err != nil {
			return err
		}
		if err := tx.Where("name = ?", name).First(&tag).Error; err != nil {
			return err
		}
		if seenIDs[tag.ID] {
			continue
		}
		seenIDs[tag.ID] = true
		if err := tx.Create(&model.ArticleTag{ArticleID: id, TagID: tag.ID}).Error; err != nil {
			return err
		}
	}
	return nil
}

func articleSnapshot(a model.Article) model.ArticleSnapshot {
	return model.ArticleSnapshot{Title: a.Title, Content: a.Content, ContentFormat: a.ContentFormat, Pic: a.Pic, Fid: a.Fid, Type: a.Type, Excerpt: a.Excerpt, SeoTitle: a.SeoTitle, SeoDescription: a.SeoDescription, TagNames: a.TagNames}
}

func recordRevision(tx *gorm.DB, a model.Article, reason string) error {
	payload, err := json.Marshal(articleSnapshot(a))
	if err != nil {
		return err
	}
	if a.Version == 0 {
		a.Version = 1
	}
	revision := model.ArticleRevision{ArticleID: a.ID, Version: a.Version, Title: a.Title, SavedAt: time.Now().Unix(), Reason: reason, Snapshot: string(payload)}
	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&revision).Error; err != nil {
		return err
	}
	var cutoff []uint64
	if err := tx.Model(&model.ArticleRevision{}).Where("article_id = ?", a.ID).Order("version desc").Offset(maxRevisions-1).Limit(1).Pluck("version", &cutoff).Error; err != nil {
		return err
	}
	if len(cutoff) > 0 {
		return tx.Where("article_id = ? AND version < ?", a.ID, cutoff[0]).Delete(&model.ArticleRevision{}).Error
	}
	return nil
}

func (s *ArticleService) createWithWorkflow(a *model.Article) error {
	now := time.Now().Unix()
	if err := validateArticle(a, now); err != nil {
		return err
	}
	a.ID = 0
	a.Ctime = int(now)
	a.EditTime = int(now)
	a.Version = 1
	a.View = 0
	if a.Uid <= 0 {
		a.Uid = 1
	}
	err := global.BLOG_DB.Transaction(func(tx *gorm.DB) error {
		if err := validateArticleCategory(tx, *a); err != nil {
			return err
		}
		if err := tx.Create(a).Error; err != nil {
			return err
		}
		if err := replaceArticleTags(tx, a.ID, a.TagNames); err != nil {
			return err
		}
		if err := attachArticleTags(tx, a); err != nil {
			return err
		}
		return recordRevision(tx, *a, "create")
	})
	if err == nil {
		s.afterArticleChange(*a)
	}
	return err
}

func (s *ArticleService) updateWithWorkflow(input *model.Article) error {
	now := time.Now().Unix()
	if err := validateArticle(input, now); err != nil {
		return err
	}
	var saved model.Article
	err := global.BLOG_DB.Transaction(func(tx *gorm.DB) error {
		var current model.Article
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&current, input.ID).Error; err != nil {
			return err
		}
		if current.Version != input.Version || input.Version == 0 {
			return ErrArticleConflict
		}
		if err := validateArticleCategory(tx, *input); err != nil {
			return err
		}
		if err := attachArticleTags(tx, &current); err != nil {
			return err
		}
		if err := recordRevision(tx, current, "before-update"); err != nil {
			return err
		}
		saved = current
		saved.Title = input.Title
		saved.Content = input.Content
		saved.ContentFormat = input.ContentFormat
		saved.Pic = input.Pic
		saved.Fid = input.Fid
		saved.Type = input.Type
		saved.IsTop = input.IsTop
		saved.Status = input.Status
		saved.PublishAt = input.PublishAt
		saved.Excerpt = input.Excerpt
		saved.SeoTitle = input.SeoTitle
		saved.SeoDescription = input.SeoDescription
		saved.Version++
		saved.EditTime = int(now)
		if (current.Status == model.ArticleDraft || current.Status == model.ArticleScheduled) && saved.Status == model.ArticlePublished {
			saved.Ctime = int(now)
		}
		if err := tx.Save(&saved).Error; err != nil {
			return err
		}
		if err := replaceArticleTags(tx, saved.ID, input.TagNames); err != nil {
			return err
		}
		if err := attachArticleTags(tx, &saved); err != nil {
			return err
		}
		return recordRevision(tx, saved, "update")
	})
	if err == nil {
		*input = saved
		s.afterArticleChange(saved)
	}
	return err
}

func (s *ArticleService) afterArticleChange(a model.Article) {
	seoService.ResetSitemap()
	seoService.ResetFeed()
	if a.Status == model.ArticlePublished {
		seoService.PushArticles([]string{seoService.ArticleURL(a.Fid, a.ID)})
	}
}

func (s *ArticleService) PatchArticleState(patch model.ArticleStatePatch) (model.Article, error) {
	var a model.Article
	if patch.Status == nil && patch.IsTop == nil {
		return a, invalidArticle("请选择需要调整的状态")
	}
	err := global.BLOG_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&a, patch.ID).Error; err != nil {
			return err
		}
		if a.Version != patch.Version || patch.Version == 0 {
			return ErrArticleConflict
		}
		if err := attachArticleTags(tx, &a); err != nil {
			return err
		}
		if err := recordRevision(tx, a, "before-state-change"); err != nil {
			return err
		}
		previous := a.Status
		if patch.Status != nil {
			a.Status = *patch.Status
			if a.Status == model.ArticleScheduled {
				return invalidArticle("请在编辑页设置定时发布时间")
			}
		}
		if patch.IsTop != nil {
			a.IsTop = *patch.IsTop
		}
		if err := validateArticle(&a, time.Now().Unix()); err != nil {
			return err
		}
		if err := validateArticleCategory(tx, a); err != nil {
			return err
		}
		if (previous == model.ArticleDraft || previous == model.ArticleScheduled) && a.Status == model.ArticlePublished {
			a.Ctime = int(time.Now().Unix())
		}
		a.EditTime = int(time.Now().Unix())
		a.Version++
		if err := tx.Save(&a).Error; err != nil {
			return err
		}
		return recordRevision(tx, a, "state-change")
	})
	if err == nil {
		s.afterArticleChange(a)
	}
	return a, err
}

func (s *ArticleService) ListRevisions(id uint) ([]model.ArticleRevision, error) {
	revisions := []model.ArticleRevision{}
	err := global.BLOG_DB.Select("id, article_id, version, title, saved_at, reason").Where("article_id = ?", id).Order("version desc, id desc").Limit(maxRevisions).Find(&revisions).Error
	return revisions, err
}

func (s *ArticleService) GetRevision(articleID, revisionID uint) (model.ArticleSnapshot, error) {
	var revision model.ArticleRevision
	var snapshot model.ArticleSnapshot
	if err := global.BLOG_DB.Where("id = ? AND article_id = ?", revisionID, articleID).First(&revision).Error; err != nil {
		return snapshot, err
	}
	err := json.Unmarshal([]byte(revision.Snapshot), &snapshot)
	return snapshot, err
}

// PublishDue 原子认领到期文章；多实例、重启补发均不会重复发布。
func (s *ArticleService) PublishDue(ctx context.Context, now int64) (int, error) {
	var ids []uint
	if err := global.BLOG_DB.WithContext(ctx).Model(&model.Article{}).Where("status = ? AND publish_at > 0 AND publish_at <= ?", model.ArticleScheduled, now).Order("publish_at asc, id asc").Limit(100).Pluck("id", &ids).Error; err != nil {
		return 0, err
	}
	count := 0
	for _, id := range ids {
		var a model.Article
		published := false
		err := global.BLOG_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&a, id).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil
				}
				return err
			}
			if a.Status != model.ArticleScheduled || a.PublishAt <= 0 || a.PublishAt > now {
				return nil
			}
			if err := attachArticleTags(tx, &a); err != nil {
				return err
			}
			if err := recordRevision(tx, a, "before-scheduled-publish"); err != nil {
				return err
			}
			a.Status = model.ArticlePublished
			a.Ctime = int(a.PublishAt)
			a.EditTime = int(now)
			a.Version++
			if err := tx.Save(&a).Error; err != nil {
				return err
			}
			if err := recordRevision(tx, a, "scheduled-publish"); err != nil {
				return err
			}
			published = true
			return nil
		})
		if err != nil {
			return count, fmt.Errorf("发布定时文章 %d: %w", id, err)
		}
		if published {
			count++
			s.afterArticleChange(a)
		}
	}
	return count, nil
}

func RunScheduledPublisher(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	run := func() {
		service := ArticleService{}
		if _, err := service.PublishDue(ctx, time.Now().Unix()); err != nil && ctx.Err() == nil {
			global.BLOG_LOG.Error("定时发布失败，将在下次检查重试", zap.Error(err))
		}
	}
	run()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			run()
		}
	}
}

func (s *ArticleService) ReadingArticles(ids []uint) ([]model.ArticleBrief, error) {
	articles := []model.ArticleBrief{}
	if len(ids) == 0 || len(ids) > 100 {
		return articles, invalidArticle("一次最多读取 100 篇文章")
	}
	err := global.BLOG_DB.Model(&model.Article{}).Select("id, fid, title, pic, view, ctime").Where("id IN ? AND status = ?", ids, model.ArticlePublished).Find(&articles).Error
	for i := range articles {
		articles[i].Pic = publicMediaURL(articles[i].Pic)
	}
	return articles, err
}

func (s *ArticleService) deleteArticles(ids []int) error {
	if len(ids) == 0 || len(ids) > 100 {
		return invalidArticle("请选择 1 到 100 篇文章")
	}
	for _, id := range ids {
		if id <= 0 {
			return invalidArticle("文章编号不合法")
		}
	}
	err := global.BLOG_DB.Transaction(func(tx *gorm.DB) error {
		// Lock articles in a stable order before removing associated rows.
		var existing []model.Article
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id IN ?", ids).Order("id asc").Find(&existing).Error; err != nil {
			return err
		}
		for _, row := range []struct {
			model interface{}
			field string
		}{{&model.ArticleTag{}, "article_id"}, {&model.ArticleRevision{}, "article_id"}, {&model.ArticleLike{}, "aid"}, {&model.Comment{}, "aid"}} {
			if err := tx.Where(row.field+" IN ?", ids).Delete(row.model).Error; err != nil {
				return err
			}
		}
		return tx.Delete(&model.Article{}, ids).Error
	})
	if err == nil {
		seoService.ResetSitemap()
		seoService.ResetFeed()
	}
	return err
}
