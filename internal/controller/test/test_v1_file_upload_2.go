package test

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"devmgr/api/test/v1"
)

func (c *ControllerV1) FileUpload2(ctx context.Context, req *v1.FileUpload2Req) (res *v1.FileUpload2Res, err error) {

	names, err := req.File.Save("/home/shr/temp/")
	if err != nil {
		return nil, err
	}

	g.Log().Line().Info(ctx, "productId:", req.ProductId)
	g.Log().Line().Info(ctx, "names:", names)

	return
}
