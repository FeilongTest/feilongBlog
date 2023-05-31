package utils

var (
	IdVerify       = Rules{"ID": []string{NotEmpty()}}
	LoginVerify    = Rules{"Username": {NotEmpty()}, "Password": {NotEmpty()}}
	CategoryVerify = Rules{"Name": {NotEmpty()}}
	ArticleVerify  = Rules{"Title": {NotEmpty()}}
)
