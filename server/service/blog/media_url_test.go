package blog

import (
	"feilongBlog/global"
	"testing"
)

func TestPublicMediaURLRewritesLegacyAvatar(t *testing.T) {
	previous := global.BLOG_CONFIG.R2.PublicURL
	global.BLOG_CONFIG.R2.PublicURL = "https://img.teshh.com/"
	t.Cleanup(func() { global.BLOG_CONFIG.R2.PublicURL = previous })

	got := publicMediaURL("./Public/Uploads/UserPic/1.png")
	want := "https://img.teshh.com/legacy/Public/Uploads/UserPic/1.png"
	if got != want {
		t.Fatalf("legacy avatar = %q, want %q", got, want)
	}

	absolute := "https://img.teshh.com/avatars/new.png"
	if publicMediaURL(absolute) != absolute {
		t.Fatalf("absolute url was rewritten")
	}
	if publicMediaURL("  ") != "" {
		t.Fatalf("blank url was not cleared")
	}
}

func TestPublicMediaURLKeepsLocalUploadsOnAPIOrigin(t *testing.T) {
	for _, raw := range []string{"uploads/file/2023/image.png", "public/uploads/file/2023/image.png", "/public/uploads/file/2023/image.png"} {
		if got := publicMediaURL(raw); got != "/blog/uploads/file/2023/image.png" {
			t.Fatalf("%s -> %s", raw, got)
		}
	}
}
