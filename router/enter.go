package router

import (
	"feilongBlog/router/blog"
)

type RouterGroup struct {
	Blog blog.ApiRouter
}

var RouterGroupApp = new(RouterGroup)
