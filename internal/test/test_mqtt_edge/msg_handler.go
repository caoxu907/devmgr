package test

import (
	"context"
	"devmgr/internal/consts"
	"devmgr/internal/model"
	"devmgr/internal/utility"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gogf/gf/v2/frame/g"
)

type mqttJob struct {
	topic   string
	payload []byte
}

type Handler struct {
	running  bool
	wg       sync.WaitGroup
	hbCancel context.CancelFunc
}

func NewHandler() *Handler {
	handler := &Handler{
		running: true,
	}
	return handler
}

func (h *Handler) Close() {
	h.running = false
	h.wg.Wait()
}

// CloudMsgHandler Handles downstream message from cloud (received via MQ), and forwards to MQTT
func (h *Handler) CloudMsgHandler(mqttClient mqtt.Client, msg mqtt.Message, deviceKey string) {
	ctx := context.Background()
	topic := fmt.Sprintf("up/egw/%s/notify", deviceKey)
	gatewayID := strings.Split(msg.Topic(), "/")[2]
	if gatewayID == "" {
		log.Println("gatewayID not found")
	}
	log.Printf("Receive Msg, topic: %s, value: %q", msg.Topic(), msg.Payload())

	// 解析JSON数据
	var message model.DeviceMessage
	err := json.Unmarshal([]byte(msg.Payload()), &message)
	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return
	}
	// 检查消息类型
	if message.MessageType == consts.Edge_command_down {
		_ = h.commandMsgHandler(ctx, topic, mqttClient, message)
	} else if message.MessageType == consts.Edge_ota_down {
		_ = h.otaMsgHandler(ctx, topic, mqttClient, message)
	} else {
		g.Log().Line().Warning(ctx, "未知的设备消息类型:", message.MessageType)
	}

}

func (h *Handler) commandMsgHandler(ctx context.Context, topic string, mqttClient mqtt.Client, message model.DeviceMessage) (err error) {
	var commandSend model.DeviceMessageCommandSend
	if err = utility.ParseContent(ctx, message.Content, &commandSend, "解析设备命令响应失败"); err != nil {
		return
	}
	g.Log().Line().Info(ctx, fmt.Sprintf("收到云端下发命令: CommandId=%d, Command=%s, Param=%s", commandSend.CommandId, commandSend.Command, commandSend.Param))

	// 模拟执行命令
	res := model.DeviceMessage{
		MessageType: consts.Edge_command_up,
		Content: model.DeviceMessageCommandResponse{
			CommandId: commandSend.CommandId,
			Level:     consts.Msg_level_info,
			Content:   fmt.Sprintf("[%s]命令已接收，正在执行...", commandSend.Command),
			Status:    consts.Command_status_executing,
		},
	}
	msgData, err := json.Marshal(res)
	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return
	}
	mqttClient.Publish(topic, 1, false, msgData)

	//模拟执行完成或者失败，根据时间随机
	var status int32 = int32(time.Now().Unix() % 2)
	var content string
	if status == 0 {
		status = consts.Command_status_completed
		content = fmt.Sprintf("[%s]命令执行成功", commandSend.Command)
	} else {
		status = consts.Command_status_failed
		content = fmt.Sprintf("[%s]命令执行失败", commandSend.Command)
	}

	time.Sleep(3 * time.Second)
	res = model.DeviceMessage{
		MessageType: consts.Edge_command_up,
		Content: model.DeviceMessageCommandResponse{
			CommandId: commandSend.CommandId,
			Level:     consts.Msg_level_info,
			Content:   content,
			Status:    status,
		},
	}
	msgData, err = json.Marshal(res)
	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return
	}
	mqttClient.Publish(topic, 1, false, msgData)

	return
}
