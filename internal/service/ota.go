// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IOta interface {
		// CheckUpgradeTaskStatus 定时检查升级任务状态
		CheckUpgradeTaskStatus(ctx context.Context)
	}
)

var (
	localOta IOta
)

func Ota() IOta {
	if localOta == nil {
		panic("implement not found for interface IOta, forgot register?")
	}
	return localOta
}

func RegisterOta(i IOta) {
	localOta = i
}
