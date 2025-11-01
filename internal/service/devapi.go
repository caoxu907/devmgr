// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"devmgr/internal/model/do"
)

type (
	IDevApi interface {
		// DevStatusCheck 检查设备状态
		DevStatusCheck(ctx context.Context)
		// 检查命令执行超时
		CommandTimeoutCheck(ctx context.Context)
		// ParseExcelFile  解析Excel文件并返回设备列表
		ParseExcelFile(ctx context.Context, fileBytes []byte) (devices []do.Device, err error)
		// ParseDeviceMessage 解析设备消息
		ParseDeviceMessage(ctx context.Context, deviceKey string, js string) (err error)
	}
)

var (
	localDevApi IDevApi
)

func DevApi() IDevApi {
	if localDevApi == nil {
		panic("implement not found for interface IDevApi, forgot register?")
	}
	return localDevApi
}

func RegisterDevApi(i IDevApi) {
	localDevApi = i
}
