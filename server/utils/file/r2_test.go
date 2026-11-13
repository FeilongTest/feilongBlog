package file

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"feilongBlog/global"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Exercise the real SDK over HTTP so its defaults cannot silently break R2 uploads
// after a security upgrade (e.g. unsupported optional CRC/checksum trailers).
func TestR2UploadAndDeleteProtocol(t *testing.T) {
	previous := global.BLOG_CONFIG.R2
	t.Cleanup(func() { global.BLOG_CONFIG.R2 = previous })
	global.BLOG_CONFIG.R2.AccountID = "test-account"
	global.BLOG_CONFIG.R2.AccessKeyID = "test-key"
	global.BLOG_CONFIG.R2.SecretAccessKey = "test-secret"
	global.BLOG_CONFIG.R2.Bucket = "test-bucket"
	global.BLOG_CONFIG.R2.PublicURL = "https://images.example.com/"
	payload := []byte("image-payload-for-r2-protocol-contract")
	events := make(chan string, 2)
	remote := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/test-bucket/images/") || !strings.HasSuffix(r.URL.Path, ".png") {
			t.Errorf("unexpected object path: %s", r.URL.Path)
		}
		if !strings.HasPrefix(r.Header.Get("Authorization"), "AWS4-HMAC-SHA256 Credential=test-key/") {
			t.Error("request did not use configured SigV4 credentials")
		}
		if r.Method == http.MethodPut {
			events <- r.URL.Path
			body, err := io.ReadAll(r.Body)
			if err != nil || !bytes.Equal(body, payload) {
				t.Errorf("object bytes changed: %q, %v", body, err)
			}
			if r.Header.Get("Content-Type") != "image/png" {
				t.Error("image MIME type lost")
			}
			for name := range r.Header {
				if strings.HasPrefix(strings.ToLower(name), "x-amz-checksum-") || strings.EqualFold(name, "X-Amz-Sdk-Checksum-Algorithm") {
					t.Errorf("optional checksum unsupported by R2: %s", name)
				}
			}
			if strings.Contains(r.Header.Get("Content-Encoding"), "aws-chunked") {
				t.Error("object unexpectedly requires AWS checksum trailers")
			}
			w.Header().Set("ETag", `"test-etag"`)
			w.WriteHeader(http.StatusOK)
		} else if r.Method == http.MethodDelete {
			events <- r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		} else {
			t.Errorf("unexpected method: %s", r.Method)
			w.WriteHeader(405)
		}
	}))
	defer remote.Close()
	storage := NewR2()
	storage.client = s3.New(storage.client.Options(), func(options *s3.Options) {
		options.BaseEndpoint = aws.String(remote.URL)
		options.UsePathStyle = true
		options.HTTPClient = remote.Client()
	})
	var form bytes.Buffer
	writer := multipart.NewWriter(&form)
	part, err := writer.CreateFormFile("file", "image.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(payload); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/upload", &form).WithContext(context.Background())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	if err := request.ParseMultipartForm(1024); err != nil {
		t.Fatal(err)
	}
	file := request.MultipartForm.File["file"][0]
	file.Header.Set("Content-Type", "image/png")
	url, key, err := storage.UploadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	uploaded := <-events
	if url != "https://images.example.com/"+key || uploaded != "/test-bucket/"+key {
		t.Fatalf("wrong public/object URLs: %s %s %s", url, key, uploaded)
	}
	if err := storage.DeleteFile(key); err != nil {
		t.Fatal(err)
	}
	deleted := <-events
	if deleted != uploaded {
		t.Fatalf("delete targeted a different object: %s", deleted)
	}
}
