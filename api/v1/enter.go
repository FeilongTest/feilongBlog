package v1

import "feilongBlog/api/v1/blog"

type ApiGroup struct {
	BlogApiGroup    blog.UserApi
	ArticleApiGroup blog.ArticleApi
}

var ApiGroupApp = new(ApiGroup)
