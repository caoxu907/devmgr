package filemgr

import (
	"context"
	v1 "devmgr/api/filemgr/v1"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

type Controller struct {
	v1.UnimplementedFileServiceServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterFileServiceServer(s.Server, &Controller{})
}

func (*Controller) GeneratePresignedDownloadURL(ctx context.Context, req *v1.PresignDownloadReq) (res *v1.PresignDownloadRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) GeneratePresignedUploadURL(ctx context.Context, req *v1.PresignUploadReq) (res *v1.PresignUploadRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) GetFileMetadata(ctx context.Context, req *v1.GetFileMetadataReq) (res *v1.GetFileMetadataRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
