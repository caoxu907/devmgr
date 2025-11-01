// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type (
	IInfluxWriter interface {
		Push(points []*write.Point)
	}
)

var (
	localInfluxWriter IInfluxWriter
)

func InfluxWriter() IInfluxWriter {
	if localInfluxWriter == nil {
		panic("implement not found for interface IInfluxWriter, forgot register?")
	}
	return localInfluxWriter
}

func RegisterInfluxWriter(i IInfluxWriter) {
	localInfluxWriter = i
}
