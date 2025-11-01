package test

import (
	"context"
	"devmgr/internal/consts"
	"devmgr/internal/model"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

func (h *Handler) StartHeartbeat(mqttClient *MqttClient, deviceKey string) {
	// 如果已有心跳在运行，先不重复启动
	if h.hbCancel != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	h.hbCancel = cancel

	topic := fmt.Sprintf("up/egw/%s/notify", deviceKey)

	payloads := model.DeviceMessage{
		MessageType: consts.Edge_heartbeat_up,
		Content: model.DeviceMessageHeartbeat{
			Version:   version,
			Cpu:       32.1,
			Mem:       1024,
			MemTotal:  4096,
			Disk:      1024,
			DiskTotal: 4096,
		},
	}

	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		msgData, err := json.Marshal(payloads)
		if err != nil {
			g.Log().Line().Error(ctx, err.Error())
		}

		// 立即发送一次
		mqttClient.Publish(topic, msgData)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				msgData, err = json.Marshal(payloads)
				if err != nil {
					g.Log().Line().Error(ctx, err.Error())
				}
				mqttClient.Publish(topic, msgData)
			}
		}
	}()
}

func (h *Handler) StopHeartbeat() {
	if h.hbCancel == nil {
		return
	}
	h.hbCancel()
	h.hbCancel = nil
}
