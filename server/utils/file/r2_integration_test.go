package file

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// 设置R2_INTEGRATION=1并提供测试凭据后运行。
func TestR2Integration(t *testing.T) {
	if os.Getenv("R2_INTEGRATION") != "1" {
		t.Skip("R2集成测试未启用")
	}
	accountID, accessKey, secretKey, bucket := os.Getenv("BLOG_R2_ACCOUNT_ID"), os.Getenv("BLOG_R2_ACCESS_KEY_ID"), os.Getenv("BLOG_R2_SECRET_ACCESS_KEY"), os.Getenv("BLOG_R2_BUCKET")
	if accountID == "" || accessKey == "" || secretKey == "" || bucket == "" {
		t.Fatal("R2测试环境变量不完整")
	}
	client := s3.New(s3.Options{Region: "auto", BaseEndpoint: aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)), Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")})
	key := fmt.Sprintf("integration-tests/codex-%d.txt", time.Now().UnixNano())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := client.PutObject(ctx, &s3.PutObjectInput{Bucket: &bucket, Key: &key, Body: bytes.NewReader([]byte("r2 connectivity ok")), ContentType: aws.String("text/plain")}); err != nil {
		t.Fatalf("put object: %v", err)
	}
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()
		if _, err := client.DeleteObject(cleanupCtx, &s3.DeleteObjectInput{Bucket: &bucket, Key: &key}); err != nil {
			t.Errorf("delete test object: %v", err)
		}
	})
	result, err := client.GetObject(ctx, &s3.GetObjectInput{Bucket: &bucket, Key: &key})
	if err != nil {
		t.Fatalf("get object: %v", err)
	}
	result.Body.Close()
}
