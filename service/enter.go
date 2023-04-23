package service

import "feilongBlog/service/blog"

type ServiceGroup struct {
	BlogService blog.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
