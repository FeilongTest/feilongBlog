package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type appConfig struct {
	Mysql struct {
		Path     string `yaml:"path"`
		Port     string `yaml:"port"`
		Config   string `yaml:"config"`
		DBName   string `yaml:"db-name"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"mysql"`
	R2 struct {
		PublicURL string `yaml:"public-url"`
	} `yaml:"r2"`
}

type article struct {
	ID      uint   `gorm:"column:id;primaryKey" json:"id"`
	Content string `gorm:"column:content" json:"content"`
	Pic     string `gorm:"column:pic" json:"pic"`
}

func (article) TableName() string { return "blog_article" }

type backupArticle struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
	Pic     string `json:"pic"`
}

type rewriteResult struct {
	Content      string
	Pic          string
	Replacements int
	AssignedPic  bool
}

var (
	srcPattern    = regexp.MustCompile(`(?i)(src\s*=\s*["'])([^"']+)(["'])`)
	uploadPattern = regexp.MustCompile(`(?i)/Public/Uploads/([^?#"']+\.(?:jpg|jpeg|png|gif|webp|bmp|svg))`)
)

func main() {
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	sourcePath := flag.String("source", "", "旧Public/Uploads目录")
	publicURL := flag.String("public-url", "", "R2公开访问基地址，默认读取配置")
	backupPath := flag.String("backup", "", "应用前JSON备份路径")
	apply := flag.Bool("apply", false, "实际写入数据库")
	assignPic := flag.Bool("assign-pic", true, "无缩略图时使用正文首张已迁移图片")
	flag.Parse()

	if *sourcePath == "" {
		fatal(errors.New("必须提供-source"))
	}
	cfg, err := loadConfig(*configPath)
	if err != nil {
		fatal(err)
	}
	baseURL := strings.TrimRight(*publicURL, "/")
	if baseURL == "" {
		baseURL = strings.TrimRight(cfg.R2.PublicURL, "/")
	}
	if baseURL == "" {
		fatal(errors.New("public-url不能为空"))
	}

	fileKeys, err := collectFileKeys(*sourcePath)
	if err != nil {
		fatal(err)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", cfg.Mysql.Username, cfg.Mysql.Password, cfg.Mysql.Path, cfg.Mysql.Port, cfg.Mysql.DBName, cfg.Mysql.Config)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fatal(fmt.Errorf("连接数据库失败: %w", err))
	}

	var articles []article
	if err = db.Find(&articles).Error; err != nil {
		fatal(err)
	}
	updates := make(map[uint]rewriteResult)
	var totalReplacements, assignedPics, emptyPics int
	for _, item := range articles {
		if strings.TrimSpace(item.Pic) == "" {
			emptyPics++
		}
		result := rewriteArticle(item, fileKeys, baseURL, *assignPic)
		if result.Content != item.Content || result.Pic != item.Pic {
			updates[item.ID] = result
			totalReplacements += result.Replacements
			if result.AssignedPic {
				assignedPics++
			}
		}
	}
	fmt.Printf("文章: %d, 空缩略图: %d, 待更新: %d, URL替换: %d, 补充缩略图: %d\n", len(articles), emptyPics, len(updates), totalReplacements, assignedPics)
	if !*apply {
		fmt.Println("当前为预演模式，数据库未修改")
		return
	}
	if len(updates) == 0 {
		return
	}
	if *backupPath == "" {
		fatal(errors.New("应用模式必须提供-backup"))
	}
	if err = writeBackup(*backupPath, articles, updates); err != nil {
		fatal(err)
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		ids := make([]int, 0, len(updates))
		for id := range updates {
			ids = append(ids, int(id))
		}
		sort.Ints(ids)
		for _, rawID := range ids {
			id := uint(rawID)
			result := updates[id]
			if err := tx.Model(&article{}).Where("ID = ?", id).Updates(map[string]any{
				"content": result.Content,
				"pic":     result.Pic,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		fatal(fmt.Errorf("事务失败，已回滚: %w", err))
	}
	fmt.Printf("数据库更新完成，备份: %s\n", *backupPath)
}

func loadConfig(path string) (appConfig, error) {
	var cfg appConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}

func collectFileKeys(root string) (map[string]string, error) {
	keys := make(map[string]string)
	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() {
			return walkErr
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relative = strings.ReplaceAll(relative, "\\", "/")
		keys[strings.ToLower(relative)] = "legacy/Public/Uploads/" + relative
		return nil
	})
	return keys, err
}

func rewriteArticle(item article, keys map[string]string, baseURL string, assignPic bool) rewriteResult {
	result := rewriteResult{Content: item.Content, Pic: item.Pic}
	var firstMigrated string
	result.Content = srcPattern.ReplaceAllStringFunc(item.Content, func(attribute string) string {
		parts := srcPattern.FindStringSubmatch(attribute)
		if len(parts) != 4 {
			return attribute
		}
		newURL, ok := migratedURL(parts[2], keys, baseURL)
		if !ok {
			return attribute
		}
		result.Replacements++
		if firstMigrated == "" {
			firstMigrated = newURL
		}
		return parts[1] + newURL + parts[3]
	})
	if newPic, ok := migratedURL(item.Pic, keys, baseURL); ok {
		result.Pic = newPic
		result.Replacements++
	} else if assignPic && strings.TrimSpace(item.Pic) == "" && firstMigrated != "" {
		result.Pic = firstMigrated
		result.AssignedPic = true
	}
	return result
}

func migratedURL(raw string, keys map[string]string, baseURL string) (string, bool) {
	decoded := html.UnescapeString(strings.TrimSpace(raw))
	match := uploadPattern.FindStringSubmatch(decoded)
	if len(match) != 2 {
		return "", false
	}
	key, ok := keys[strings.ToLower(strings.ReplaceAll(match[1], "\\", "/"))]
	if !ok {
		return "", false
	}
	return baseURL + "/" + key, true
}

func writeBackup(path string, articles []article, updates map[uint]rewriteResult) error {
	backup := make([]backupArticle, 0, len(updates))
	for _, item := range articles {
		if _, ok := updates[item.ID]; ok {
			backup = append(backup, backupArticle{ID: item.ID, Content: item.Content, Pic: item.Pic})
		}
	}
	data, err := json.MarshalIndent(struct {
		CreatedAt string          `json:"createdAt"`
		Articles  []backupArticle `json:"articles"`
	}{CreatedAt: time.Now().Format(time.RFC3339), Articles: backup}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
