package main

import (
	v1 "devmgr/api/devmgr/v1"

	"github.com/gogf/gf/contrib/registry/nacos/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	{
		grpcx.Resolver.Register(nacos.New("127.0.0.1:18848"))

		var ctx = gctx.New()
		var conn = grpcx.Client.MustNewGrpcClientConn("devmgr", grpcx.Balancer.WithRandom())
		var client = v1.NewDevmgrClient(conn)

		res, err := client.DevAuth(ctx, &v1.DevAuthReq{
			MachineCode:  "device_key_20",
			DeviceSecret: "device_secret_20",
		})
		if err != nil {
			g.Log().Line().Error(ctx, err.Error())
			return
		} else {
			g.Log().Line().Info(ctx, "ok, Response:", res.String())
		}
	}

	//{
	//	gsvc.SetRegistry(nacos.New("127.0.0.1:18848"))
	//
	//	var (
	//		ctx    = gctx.New()
	//		client = g.Client()
	//	)
	//	client.SetDiscovery(gsvc.GetRegistry())
	//	res := client.GetContent(ctx, `http://devmgr-http/product`)
	//	g.Log().Line().Info(ctx, res)
	//}

}
