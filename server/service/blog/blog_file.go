package blog

import (
	"feilongBlog/utils/file"
	"mime/multipart"
)

type FileService struct{}

//@function: UploadFile
//@description: 根据配置文件判断是文件上传到本地或者七牛云

func (e *FileService) UploadFile(header *multipart.FileHeader) (filePath, key string, err error) {
	oss := file.NewOss()
	filePath, key, err = oss.UploadFile(header)
	return
}

func (e *FileService) DelFile() error {
	oss := file.NewOss()
	return oss.DeleteFile("")
}
