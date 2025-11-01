package dbres

import (
	"context"
	v1 "devmgr/api/dbres/v1"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

type Controller struct {
	v1.UnimplementedDbresServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterDbresServer(s.Server, &Controller{})
}

func (*Controller) CreateInstance(ctx context.Context, req *v1.CreateInstanceReq) (res *v1.CreateInstanceRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) UpdateInstance(ctx context.Context, req *v1.UpdateInstanceReq) (res *v1.UpdateInstanceRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) CreateDatabase(ctx context.Context, req *v1.CreateDatabaseReq) (res *v1.CreateDatabaseRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) UpdateDatabase(ctx context.Context, req *v1.UpdateDatabaseReq) (res *v1.UpdateDatabaseRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) CreateRDBSchema(ctx context.Context, req *v1.CreateRDBSchemaReq) (res *v1.CreateRDBSchemaRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) UpdateRDBSchema(ctx context.Context, req *v1.UpdateRDBSchemaReq) (res *v1.UpdateRDBSchemaRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) CreateRDBTable(ctx context.Context, req *v1.CreateRDBTableReq) (res *v1.CreateRDBTableRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) UpdateRDBTable(ctx context.Context, req *v1.UpdateRDBTableReq) (res *v1.UpdateRDBTableRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) CreateHisdbTable(ctx context.Context, req *v1.CreateHisdbTableReq) (res *v1.CreateHisdbTableRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) UpdateHisdbTable(ctx context.Context, req *v1.UpdateHisdbTableReq) (res *v1.UpdateHisdbTableRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) GetRdbResTree(ctx context.Context, req *v1.GetRdbResTreeReq) (res *v1.GetRdbResTreeRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) GetHisdbResTree(ctx context.Context, req *v1.GetHisdbResTreeReq) (res *v1.GetHisdbResTreeRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) GetField(ctx context.Context, req *v1.GetFieldReq) (res *v1.GetFieldRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) CreateBinding(ctx context.Context, req *v1.CreateBindingReq) (res *v1.CreateBindingRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) UpdateBinding(ctx context.Context, req *v1.UpdateBindingReq) (res *v1.UpdateBindingRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) GetAllBindings(ctx context.Context, req *v1.GetAllBindingsReq) (res *v1.GetAllBindingsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) ReadData(ctx context.Context, req *v1.ReadDataReq) (res *v1.ReadDataRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
