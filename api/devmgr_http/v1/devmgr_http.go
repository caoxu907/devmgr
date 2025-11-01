package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type DevAuthReq struct {
	g.Meta   `path:"/device/auth" method:"post" tags:"device" summary:"device auth"`
	Username string `v:"required"  dc:"username"`
	Password string `v:"required"  dc:"password"`
}
type DevAuthRes struct {
	Result      string     `json:"result" dc:"result"`
	ClientAttrs ClientAttr `json:"client_attrs" dc:"client_attrs"`
}

type ClientAttr struct {
	DeviceKey string `json:"device_key" dc:"device_key"`
}

type HealthReq struct {
	g.Meta `path:"/health" method:"get" tags:"devmgr" summary:"devmgr health"`
}
type HealthRes struct {
}
