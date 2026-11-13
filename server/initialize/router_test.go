package initialize

import (
	"feilongBlog/global"
	model "feilongBlog/model/blog"
	"feilongBlog/utils"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestArticleManagementRequiresAdmin(t *testing.T) {
	previousConfig, previousLog := global.BLOG_CONFIG, global.BLOG_LOG
	t.Cleanup(func() { global.BLOG_CONFIG = previousConfig; global.BLOG_LOG = previousLog })
	global.BLOG_CONFIG.Local.Path = "/uploads/file"
	global.BLOG_CONFIG.Local.StorePath = t.TempDir()
	global.BLOG_CONFIG.JWT.SigningKey = "test-signing-key-not-a-production-secret"
	global.BLOG_LOG = zap.NewNop()
	router := Routers()
	makeToken := func(admin int) string {
		claims := model.CustomClaims{BaseClaims: model.BaseClaims{ID: 1, Admin: admin}, StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour).Unix()}}
		token, err := utils.NewJWT().CreateToken(claims)
		if err != nil {
			t.Fatal(err)
		}
		return token
	}
	// No requests in this test may reach a database or mutate article data.
	for _, route := range []struct{ method, path string }{
		{"GET", "getArticleList"}, {"POST", "getArticle"}, {"POST", "createArticle"},
		{"PUT", "updateArticle"}, {"DELETE", "delArticle"}, {"DELETE", "delArticleByIds"},
	} {
		for _, prefix := range []string{"/article/", "/blog/article/"} {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(route.method, prefix+route.path, strings.NewReader("{}"))
			router.ServeHTTP(w, req)
			if w.Code != 404 || !strings.Contains(w.Header().Get("Content-Type"), "application/json") {
				t.Fatalf("public %s: %d %s", req.URL, w.Code, w.Body.String())
			}
		}
		for _, test := range []struct {
			token  string
			status int
		}{{"", 401}, {"invalid", 401}, {makeToken(0), 403}} {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(route.method, "/admin/article/"+route.path, strings.NewReader("{}"))
			req.Header.Set("x-token", test.token)
			router.ServeHTTP(w, req)
			if w.Code != test.status {
				t.Fatalf("%s: got %d want %d", req.URL, w.Code, test.status)
			}
		}
	}
	// A valid administrator reaches request validation, rather than being blocked.
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/admin/article/createArticle", strings.NewReader("invalid-json"))
	req.Header.Set("x-token", makeToken(1))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	if w.Code != 200 || strings.Contains(w.Body.String(), "仅管理员") {
		t.Fatalf("admin request rejected: %d %s", w.Code, w.Body.String())
	}
}
