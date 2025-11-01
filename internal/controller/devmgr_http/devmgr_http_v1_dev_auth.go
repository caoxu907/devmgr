package devmgr_http

import (
	"context"
	"devmgr/internal/dao"
	"devmgr/internal/model/do"
	"devmgr/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"

	"devmgr/api/devmgr_http/v1"
)

func (c *ControllerV1) DevAuth(ctx context.Context, req *v1.DevAuthReq) (res *v1.DevAuthRes, err error) {
	res = &v1.DevAuthRes{Result: "deny"}
	resHandle := g.RequestFromCtx(ctx).Response
	resHandle.WriteStatus(200)

	var devices []entity.Device
	if err = dao.Device.Ctx(ctx).
		Where(do.Device{
			MachineCode:  req.Username,
			DeviceSecret: req.Password,
			Deleted:      false,
		}).
		Scan(&devices); err != nil {
		g.Log().Line().Error(ctx, err.Error())
		// 直接写出原始 JSON 并退出，避免中间件再次封装
		resHandle.WriteJsonExit(res)
		return nil, nil
	}

	// 认证失败时直接返回 deny
	if len(devices) != 1 {
		g.Log().Line().Error(ctx, "认证失败", len(devices))
		resHandle.WriteJsonExit(res)
		return nil, nil
	}

	// 认证通过，直接返回 allow
	res = &v1.DevAuthRes{Result: "allow", ClientAttrs: v1.ClientAttr{DeviceKey: devices[0].DeviceKey}}
	resHandle.WriteJsonExit(res)
	return nil, nil
}
