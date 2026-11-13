package initialize

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"feilongBlog/utils"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

func TestLocalUploadIsUniqueAndUsesImageExtension(t *testing.T) {
	oldConfig, oldLog := global.BLOG_CONFIG, global.BLOG_LOG
	t.Cleanup(func() { global.BLOG_CONFIG = oldConfig; global.BLOG_LOG = oldLog })
	global.BLOG_LOG = zap.NewNop()
	global.BLOG_CONFIG.Local.Path = "/uploads/file"
	global.BLOG_CONFIG.Local.StorePath = t.TempDir()
	global.BLOG_CONFIG.System.OssType = "local"
	global.BLOG_CONFIG.JWT.SigningKey = "test-local-image-upload-authentication-key"
	router := Routers()
	token, err := utils.NewJWT().CreateToken(model.CustomClaims{BaseClaims: model.BaseClaims{ID: 1, Admin: 1}, StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour).Unix()}})
	if err != nil {
		t.Fatal(err)
	}
	image, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+aDBkAAAAASUVORK5CYII=")
	if err != nil {
		t.Fatal(err)
	}
	var previous string
	for i := 0; i < 2; i++ {
		var body bytes.Buffer
		form := multipart.NewWriter(&body)
		file, err := form.CreateFormFile("file", "same-name.html")
		if err != nil {
			t.Fatal(err)
		}
		file.Write(image)
		form.Close()
		req := httptest.NewRequest("POST", "/admin/file/upload", &body)
		req.Header.Set("Content-Type", form.FormDataContentType())
		req.Header.Set("x-token", token)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)
		var payload struct {
			Code int
			Data struct {
				URL string `json:"url"`
			}
		}
		if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
			t.Fatal(err)
		}
		if payload.Code != 0 || !strings.HasSuffix(payload.Data.URL, ".png") || payload.Data.URL == previous {
			t.Fatalf("unsafe or duplicate upload=%s", response.Body.String())
		}
		previous = payload.Data.URL
		get := httptest.NewRecorder()
		router.ServeHTTP(get, httptest.NewRequest("GET", strings.TrimPrefix(previous, "/blog"), nil))
		if get.Code != 200 || get.Header().Get("Content-Type") != "image/png" || !bytes.Equal(get.Body.Bytes(), image) {
			t.Fatalf("image response=%d %s", get.Code, get.Header())
		}
	}
}
