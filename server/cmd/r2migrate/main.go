package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gopkg.in/yaml.v2"
)

type appConfig struct {
	R2 struct {
		AccountID       string `yaml:"account-id"`
		AccessKeyID     string `yaml:"access-key-id"`
		SecretAccessKey string `yaml:"secret-access-key"`
		Bucket          string `yaml:"bucket"`
	} `yaml:"r2"`
}

type migrationItem struct {
	Source string `json:"source"`
	Key    string `json:"key"`
	Size   int64  `json:"size"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

var imageExtensions = map[string]struct{}{
	".jpg": {}, ".jpeg": {}, ".png": {}, ".gif": {}, ".webp": {}, ".bmp": {}, ".svg": {},
}

func main() {
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	sourcePath := flag.String("source", "", "待迁移图片目录")
	prefix := flag.String("prefix", "legacy/Public/Uploads", "R2对象Key前缀")
	manifestPath := flag.String("manifest", "", "可选的迁移结果JSON路径")
	workers := flag.Int("workers", 4, "并发上传数量")
	dryRun := flag.Bool("dry-run", false, "仅扫描，不上传")
	flag.Parse()

	if *sourcePath == "" {
		fatal(errors.New("必须提供-source"))
	}
	if *workers < 1 || *workers > 16 {
		fatal(errors.New("workers必须在1到16之间"))
	}

	cfg, err := loadConfig(*configPath)
	if err != nil {
		fatal(err)
	}
	files, err := collectImages(*sourcePath, *prefix)
	if err != nil {
		fatal(err)
	}
	if *dryRun {
		for i := range files {
			files[i].Status = "dry-run"
		}
		printSummary(files)
		writeManifest(*manifestPath, files)
		return
	}

	if cfg.R2.AccountID == "" || cfg.R2.AccessKeyID == "" || cfg.R2.SecretAccessKey == "" || cfg.R2.Bucket == "" {
		fatal(errors.New("R2配置不完整"))
	}

	client := s3.New(s3.Options{
		Region:       "auto",
		BaseEndpoint: aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2.AccountID)),
		Credentials:  credentials.NewStaticCredentialsProvider(cfg.R2.AccessKeyID, cfg.R2.SecretAccessKey, ""),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	jobs := make(chan int)
	var completed atomic.Int64
	var wg sync.WaitGroup
	for worker := 0; worker < *workers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				migrateFile(ctx, client, cfg.R2.Bucket, &files[index])
				done := completed.Add(1)
				if done%25 == 0 || done == int64(len(files)) {
					fmt.Printf("进度: %d/%d\n", done, len(files))
				}
			}
		}()
	}
	for index := range files {
		jobs <- index
	}
	close(jobs)
	wg.Wait()

	printSummary(files)
	writeManifest(*manifestPath, files)
	for _, item := range files {
		if item.Status == "failed" {
			os.Exit(1)
		}
	}
}

func loadConfig(path string) (appConfig, error) {
	var cfg appConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("读取配置失败: %w", err)
	}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("解析配置失败: %w", err)
	}
	return cfg, nil
}

func collectImages(root, prefix string) ([]migrationItem, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	var items []migrationItem
	err = filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}
		if _, ok := imageExtensions[strings.ToLower(filepath.Ext(path))]; !ok {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		items = append(items, migrationItem{
			Source: path,
			Key:    strings.Trim(strings.ReplaceAll(filepath.Join(prefix, relative), "\\", "/"), "/"),
			Size:   info.Size(),
		})
		return nil
	})
	sort.Slice(items, func(i, j int) bool { return items[i].Key < items[j].Key })
	return items, err
}

func migrateFile(ctx context.Context, client *s3.Client, bucket string, item *migrationItem) {
	head, err := client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(bucket), Key: aws.String(item.Key)})
	if err == nil && head.ContentLength == item.Size {
		item.Status = "skipped"
		return
	}

	file, err := os.Open(item.Source)
	if err != nil {
		item.Status, item.Error = "failed", err.Error()
		return
	}
	defer file.Close()

	contentType := mime.TypeByExtension(strings.ToLower(filepath.Ext(item.Source)))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       aws.String(bucket),
		Key:          aws.String(item.Key),
		Body:         file,
		ContentType:  aws.String(contentType),
		CacheControl: aws.String("public, max-age=86400"),
	})
	if err != nil {
		item.Status, item.Error = "failed", err.Error()
		return
	}
	item.Status = "uploaded"
}

func printSummary(items []migrationItem) {
	counts := make(map[string]int)
	var bytes int64
	for _, item := range items {
		counts[item.Status]++
		bytes += item.Size
	}
	fmt.Printf("扫描: %d, 大小: %.2f MiB, 上传: %d, 跳过: %d, 失败: %d, 预演: %d\n",
		len(items), float64(bytes)/(1024*1024), counts["uploaded"], counts["skipped"], counts["failed"], counts["dry-run"])
}

func writeManifest(path string, items []migrationItem) {
	if path == "" {
		return
	}
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		fatal(err)
	}
	if err = os.WriteFile(path, data, 0600); err != nil {
		fatal(err)
	}
	fmt.Printf("迁移清单: %s\n", path)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
