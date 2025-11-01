package test

import (
	"context"
	"devmgr/internal/consts"
	"devmgr/internal/model"
	"devmgr/internal/utility"
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gogf/gf/v2/frame/g"
)

func (h *Handler) otaMsgHandler(ctx context.Context, topic string, mqttClient mqtt.Client, message model.DeviceMessage) (err error) {
	var commandSend model.DeviceMessageOtaSend
	if err = utility.ParseContent(ctx, message.Content, &commandSend, "解析ota命令失败"); err != nil {
		return
	}
	g.Log().Line().Info(ctx, fmt.Sprintf("收到云端下发OTA命令: Id=%d, Type=%s, Url=%s, PackageVersion=%s",
		commandSend.Id, commandSend.MessageType, commandSend.Url, commandSend.PackageVersion))

	var msgData []byte

	switch commandSend.MessageType {
	case consts.Ota_command_type_cancel:
		res := model.DeviceMessage{
			MessageType: consts.Edge_ota_up,
			Content: model.DeviceMessageOtaResponse{
				MessageType: commandSend.MessageType,
				Id:          commandSend.Id,
				Level:       consts.Msg_level_info,
				Content:     "OTA取消命令已接收",
				Status:      consts.Command_status_failed,
			},
		}
		msgData, err = json.Marshal(res)
		if err != nil {
			g.Log().Line().Error(ctx, err.Error())
			return
		}
		mqttClient.Publish(topic, 1, false, msgData)

	case consts.Ota_command_type_upgrade:
		// 模拟执行OTA命令
		res := model.DeviceMessage{
			MessageType: consts.Edge_ota_up,
			Content: model.DeviceMessageOtaResponse{
				MessageType: commandSend.MessageType,
				Id:          commandSend.Id,
				Level:       consts.Msg_level_info,
				Content:     fmt.Sprintf("OTA命令[%d]已接收，正在执行...", commandSend.Id),
				Status:      consts.Command_status_executing,
			},
		}
		msgData, err = json.Marshal(res)
		if err != nil {
			g.Log().Line().Error(ctx, err.Error())
			return
		}
		mqttClient.Publish(topic, 1, false, msgData)

		//模拟执行完成或者失败，根据时间随机
		var status = int32(time.Now().Unix() % 2)
		var level int32
		var content string
		if status == 0 {
			status = consts.Command_status_completed
			level = consts.Msg_level_info
			content = fmt.Sprintf("OTA命令[%d]执行成功", commandSend.Id)
		} else {
			status = consts.Command_status_failed
			level = consts.Msg_level_error
			content = fmt.Sprintf("OTA命令[%d]执行失败", commandSend.Id)
		}

		time.Sleep(10 * time.Second)
		res = model.DeviceMessage{
			MessageType: consts.Edge_ota_up,
			Content: model.DeviceMessageOtaResponse{
				MessageType: commandSend.MessageType,
				Id:          commandSend.Id,
				Level:       level,
				Content:     content,
				Status:      status,
			},
		}
		msgData, err = json.Marshal(res)
		if err != nil {
			g.Log().Line().Error(ctx, err.Error())
			return
		}
		mqttClient.Publish(topic, 1, false, msgData)
	}

	return
}
