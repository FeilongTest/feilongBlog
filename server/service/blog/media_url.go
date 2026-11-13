package blog

import (
	"feilongBlog/global"
	"strings"
)

// publicMediaURL 把旧站 ./Public/Uploads 头像和图片改写成对象存储上的公开地址。
// 已经是绝对地址的保持不变。无法识别的相对路径原样返回，避免丢掉新上传的本地路径。
func publicMediaURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") || strings.HasPrefix(raw, "//") {
		return raw
	}

	cleaned := strings.TrimPrefix(raw, "./")
	cleaned = strings.TrimPrefix(cleaned, "/")
	local := strings.TrimPrefix(cleaned, "public/")
	if strings.HasPrefix(local, "uploads/file/") {
		return "/blog/" + local
	}
	lower := strings.ToLower(cleaned)
	const marker = "public/uploads/"
	index := strings.Index(lower, marker)
	if index < 0 {
		return raw
	}

	base := strings.TrimRight(strings.TrimSpace(global.BLOG_CONFIG.R2.PublicURL), "/")
	if base == "" {
		return raw
	}
	return base + "/legacy/Public/Uploads/" + cleaned[index+len(marker):]
}
