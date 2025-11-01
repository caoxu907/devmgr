// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package test

import (
	"context"

	"devmgr/api/test/v1"
)

type ITestV1 interface {
	FileGet(ctx context.Context, req *v1.FileGetReq) (res *v1.FileGetRes, err error)
	FileUpload(ctx context.Context, req *v1.FileUploadReq) (res *v1.FileUploadRes, err error)
	FileUpload2(ctx context.Context, req *v1.FileUpload2Req) (res *v1.FileUpload2Res, err error)
}
