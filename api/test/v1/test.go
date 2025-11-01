package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type FileGetReq struct {
	g.Meta `path:"/test/template" method:"get" tags:"test" summary:"file get test"`
}
type FileGetRes struct{}

type FileUploadReq struct {
	g.Meta `path:"/test/import" method:"post" tags:"test" summary:"file upload test"`
}
type FileUploadRes struct{}

type FileUpload2Req struct {
	g.Meta    `path:"/test/import2" method:"post" mime:"multipart/form-data" tags:"test" summary:"file upload test"`
	File      *ghttp.UploadFile `v:"required" form:"file" dc:"File to upload"`
	ProductId int64             `v:"required" form:"productId" dc:"Product ID"`
}
type FileUpload2Res struct{}
