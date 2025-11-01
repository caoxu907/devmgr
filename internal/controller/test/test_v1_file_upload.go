package test

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"devmgr/api/test/v1"
)

func (c *ControllerV1) FileUpload(ctx context.Context, req *v1.FileUploadReq) (res *v1.FileUploadRes, err error) {
	request := g.RequestFromCtx(ctx)
	if request == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "invalid request")
	}

	files := request.GetUploadFiles("file")
	productId := request.MultipartForm.Value["productId"]

	names, err := files.Save("/home/shr/temp/")
	if err != nil {
		//request.Response.WriteExit(err)
		return nil, err
	}

	//request.Response.WriteExit("upload successfully: ", names)

	g.Log().Line().Info(ctx, "productId:", productId)
	g.Log().Line().Info(ctx, "names:", names)
	return
}
