package initialize

import (
	"context"
	devapi "devmgr/internal/logic/devapi"
	mqapi "devmgr/internal/logic/mq"
	otaapi "devmgr/internal/logic/ota"
	"devmgr/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func InitApi(ctx context.Context) error {
	initIDevApi()
	initIOta()
	initIMQClient(ctx)
	return nil
}

func initIDevApi() {
	service.RegisterDevApi(devapi.New())
}

func initIOta() {
	service.RegisterOta(otaapi.New())
}

func initIMQClient(ctx context.Context) {
	brokersAddress := g.Cfg().MustGet(ctx, "kafka.brokers")
	groupId := g.Cfg().MustGet(ctx, "kafka.groupId")
	g.Log().Info(ctx, "brokers:", brokersAddress, "group:", groupId)
	service.RegisterMQClient(mqapi.New(brokersAddress.String(), groupId.String()))

}
