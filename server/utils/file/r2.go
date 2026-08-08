package file

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"feilongBlog/global"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2 struct {
	client *s3.Client
}

// NewR2 R2的实例化方法
func NewR2() *R2 {
	cfg := global.BLOG_CONFIG.R2
	client := s3.New(s3.Options{
		Region:       "auto",
		BaseEndpoint: aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID)),
		Credentials:  credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	})
	return &R2{client: client}
}

//@object: *R2
//@function: UploadFile
//@description: 上传文件到Cloudflare R2
//@param: file *multipart.FileHeader
//@return: string, string, error

func (r *R2) UploadFile(file *multipart.FileHeader) (string, string, error) {
	f, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	ext := strings.ToLower(filepath.Ext(file.Filename))
	key := fmt.Sprintf("images/%s/%d%s", time.Now().Format("2006/01/02"), time.Now().UnixNano(), ext)
	_, err = r.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(global.BLOG_CONFIG.R2.Bucket),
		Key:         aws.String(key),
		Body:        f,
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", "", fmt.Errorf("上传文件到R2失败, err: %w", err)
	}
	filePath := strings.TrimRight(global.BLOG_CONFIG.R2.PublicURL, "/") + "/" + key
	return filePath, key, nil
}

//@object: *R2
//@function: DeleteFile
//@description: 删除Cloudflare R2中的文件
//@param: key string
//@return: error

func (r *R2) DeleteFile(key string) error {
	_, err := r.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(global.BLOG_CONFIG.R2.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("删除R2文件失败, err: %w", err)
	}
	return nil
}
