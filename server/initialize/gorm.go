package initialize

import (
	"feilongBlog/model/blog"
	"os"

	"feilongBlog/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Gorm 初始化数据库并产生数据库全局变量
// Author SliverHorn
func Gorm() *gorm.DB {
	return GormMysql()
}

// RegisterTables 注册数据库表专用
// Author SliverHorn
func RegisterTables(db *gorm.DB) {
	err := db.AutoMigrate(
		// 系统模块表
		blog.User{},
		blog.ArticleLike{},
		blog.BlogStatistic{},
		blog.BlogVisitor{},
		blog.Setting{},
		blog.Tag{}, blog.ArticleTag{}, blog.ArticleRevision{},
	)
	if err != nil {
		global.BLOG_LOG.Error("register table failed", zap.Error(err))
		os.Exit(1)
	}
	if !db.Migrator().HasColumn(&blog.Article{}, "ContentFormat") {
		if err = db.Migrator().AddColumn(&blog.Article{}, "ContentFormat"); err != nil {
			global.BLOG_LOG.Error("add article content format failed", zap.Error(err))
			os.Exit(1)
		}
	}
	if err = MigratePublishing(db); err != nil {
		global.BLOG_LOG.Error("publishing migration failed", zap.Error(err))
		os.Exit(1)
	}
	global.BLOG_LOG.Info("register table success")
}

// MigratePublishing 只增加字段，并按需扩容原有标题/正文；可重复运行。
func MigratePublishing(db *gorm.DB) error {
	for _, field := range []string{"Excerpt", "SeoTitle", "SeoDescription", "PublishAt", "Version"} {
		if !db.Migrator().HasColumn(&blog.Article{}, field) {
			if err := db.Migrator().AddColumn(&blog.Article{}, field); err != nil {
				return err
			}
		}
	}
	if !db.Migrator().HasIndex(&blog.Article{}, "idx_article_schedule") {
		if err := db.Migrator().CreateIndex(&blog.Article{}, "idx_article_schedule"); err != nil {
			return err
		}
	}
	columns, err := db.Migrator().ColumnTypes(&blog.Article{})
	if err != nil {
		return err
	}
	for _, column := range columns {
		if column.Name() == "title" {
			if length, ok := column.Length(); ok && length < 200 {
				if err := db.Migrator().AlterColumn(&blog.Article{}, "Title"); err != nil {
					return err
				}
			}
		}
		if (column.Name() == "ctime" || column.Name() == "edittime") && column.DatabaseTypeName() != "bigint" && column.DatabaseTypeName() != "BIGINT" {
			field := "Ctime"
			if column.Name() == "edittime" {
				field = "EditTime"
			}
			if err := db.Migrator().AlterColumn(&blog.Article{}, field); err != nil {
				return err
			}
		}
		if column.Name() == "content" && column.DatabaseTypeName() != "longtext" && column.DatabaseTypeName() != "LONGTEXT" {
			if err := db.Migrator().AlterColumn(&blog.Article{}, "Content"); err != nil {
				return err
			}
		}
	}
	return nil
}
