package cmd

import (
	"context"
	"devmgr/internal/controller/devmgr"
	"devmgr/internal/controller/devmgr_http"
	"devmgr/internal/initialize"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"go.lanniu.top/nebula-cloud/go-scaffold/sframe"
	"google.golang.org/grpc"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start both HTTP and gRPC servers",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {

			//api初始化
			err = initialize.InitApi(ctx)
			if err != nil {
				g.Log().Fatal(ctx, "API初始化失败:", err)
				return err
			}

			//线程初始化
			err = initialize.InitThread(ctx)
			defer initialize.StopThread()
			if err != nil {
				g.Log().Fatal(ctx, "线程初始化失败:", err)
				return err
			}

			// 启动 HTTP 服务
			httpServer := g.Server("devmgr-http")
			httpServer.SetAddr(sframe.Endpoint)
			httpServer.Use(ghttp.MiddlewareHandlerResponse)
			httpServer.Group("/", func(group *ghttp.RouterGroup) {
				group.Bind(devmgr_http.NewV1())

			})

			//运行http服务
			go func() {
				httpServer.Run()
			}()
			
			// 停止http服务
			defer func(httpServer *ghttp.Server) {
				err := httpServer.Shutdown()
				if err != nil {
					g.Log().Error(ctx, "HTTP服务关闭失败:", err)
				}
			}(httpServer)

			config := grpcx.Server.NewConfig()
			config.Options = append(config.Options, []grpc.ServerOption{
				grpcx.Server.ChainUnary(grpcx.Server.UnaryValidate),
			}...)
			grpcServer := grpcx.Server.New(config)
			devmgr.Register(grpcServer)
			grpcServer.Run()

			return nil
		},
	}
)
