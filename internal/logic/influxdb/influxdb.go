package influxdb

import (
	"context"
	"devmgr/internal/service"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func New() service.IInfluxWriter {
	var ctx = gctx.New()

	url := g.Cfg().MustGet(ctx, "influxdb.url")
	g.Log().Line().Info(ctx, "InfluxDB URL:", url)
	token := g.Cfg().MustGet(ctx, "influxdb.token")
	g.Log().Line().Info(ctx, "InfluxDB Token:", token)
	org := g.Cfg().MustGet(ctx, "influxdb.org")
	bucket := g.Cfg().MustGet(ctx, "influxdb.bucket")

	client := influxdb2.NewClient(url.String(), token.String())
	writeAPI := client.WriteAPI(org.String(), bucket.String())

	instance := &sInfluxWriter{
		writeAPI: writeAPI,
		ch:       make(chan []*write.Point),
		wg:       &sync.WaitGroup{},
	}

	// 只创建一次错误监听goroutine
	go func() {
		for err := range writeAPI.Errors() {
			g.Log().Line().Error(context.Background(), "InfluxDB写入错误: ", err)
		}
	}()

	return instance
}

type sInfluxWriter struct {
	writeAPI api.WriteAPI
	ch       chan []*write.Point
	wg       *sync.WaitGroup
}

func (iw *sInfluxWriter) Push(points []*write.Point) {
	// 记录尝试写入的点数
	g.Log().Debug(context.Background(), "尝试写入数据点数量:", len(points))

	for _, pt := range points {
		iw.writeAPI.WritePoint(pt)
	}

	// 添加日志确认写入请求已发送
	g.Log().Debug(context.Background(), "数据点已提交到WriteAPI")

	// Flush操作会阻塞直到所有点都被写入
	iw.writeAPI.Flush()
	g.Log().Debug(context.Background(), "Flush操作完成")
}
