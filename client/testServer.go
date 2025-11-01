package main

import (
	"devmgr/internal/controller/test"
	"devmgr/internal/tool"

	"github.com/gogf/gf/contrib/registry/nacos/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {

	ctx := gctx.New()

	address := g.Cfg().MustGet(ctx, "nacos.address")
	network := g.Cfg().MustGet(ctx, "network")
	ip, err := tool.GetFirstIPInSubnet(network.String())
	if err != nil {
		g.Log().Fatal(ctx, "获取网段IP失败:", err)
		return
	}
	ip = ip + ":0"

	// 启动 HTTP 服务
	gsvc.SetRegistry(nacos.New(address.String()))
	httpServer := g.Server("devmgr-test")
	httpServer.SetAddr(ip)
	httpServer.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(tool.MiddlewareHandlerResponse)
		group.Bind(test.NewV1())
	})
	go func() { httpServer.Run() }()
	defer httpServer.Shutdown()
}
