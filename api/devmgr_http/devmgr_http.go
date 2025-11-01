// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package devmgr_http

import (
	"context"

	"devmgr/api/devmgr_http/v1"
)

type IDevmgrHttpV1 interface {
	DevAuth(ctx context.Context, req *v1.DevAuthReq) (res *v1.DevAuthRes, err error)
	Health(ctx context.Context, req *v1.HealthReq) (res *v1.HealthRes, err error)
}
