package devmgr_http

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"devmgr/api/devmgr_http/v1"
)

func (c *ControllerV1) Health(ctx context.Context, req *v1.HealthReq) (res *v1.HealthRes, err error) {
	g.RequestFromCtx(ctx).Response.WriteStatusExit(200)
	return nil, nil
}
