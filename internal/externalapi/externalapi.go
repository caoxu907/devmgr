package externalapi

import (
	filemgr "devmgr/api/filemgr/v1"
	"sync"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

var (
	svcNameFilemgr = "filemgr"
)

var (
	filemgrMutex  sync.Mutex
	filemgrClient filemgr.FileServiceClient
)

func GetFilemgrClient() (filemgr.FileServiceClient, error) {
	filemgrMutex.Lock()
	defer filemgrMutex.Unlock()

	if filemgrClient != nil {
		return filemgrClient, nil
	}
	conn, err := grpcx.Client.NewGrpcClientConn(svcNameFilemgr, grpcx.Balancer.WithRandom())
	if err != nil {
		return nil, err
	}
	filemgrClient = filemgr.NewFileServiceClient(conn)
	return filemgrClient, nil
}
