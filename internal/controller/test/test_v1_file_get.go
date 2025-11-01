package test

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"devmgr/api/test/v1"
)

func (c *ControllerV1) FileGet(ctx context.Context, req *v1.FileGetReq) (res *v1.FileGetRes, err error) {
	//将本目录下的test.json
	g.RequestFromCtx(ctx).Response.ServeFileDownload("test.json")
	return
}
