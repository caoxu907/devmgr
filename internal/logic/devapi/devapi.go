package devapi

import (
	"bytes"
	"context"
	"devmgr/internal/consts"
	"devmgr/internal/dao"
	"devmgr/internal/model"
	"devmgr/internal/model/do"
	"devmgr/internal/model/entity"
	"devmgr/internal/service"
	"devmgr/internal/utility"
	"encoding/json"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/xuri/excelize/v2"
)

func New() service.IDevApi {
	return &sDevApi{}
}

type sDevApi struct{}

// DevStatusCheck 检查设备状态
func (s *sDevApi) DevStatusCheck(ctx context.Context) {
	// 分页查询设备
	pageSize := 100
	page := 1

	for {
		var devices []entity.Device
		err := dao.Device.Ctx(ctx).
			Page(page, pageSize).
			Order("id ASC").
			Scan(&devices)

		if err != nil {
			g.Log().Line().Error(ctx, "查询设备失败", err.Error())
			return
		}

		// 没有更多数据，退出循环
		if len(devices) == 0 {
			break
		}

		// 处理当前批次的设备
		for _, device := range devices {
			// 设置目标状态
			targetStatus := consts.DeviceStatusOffline
			if device.TenantId == 0 {
				targetStatus = consts.DeviceStatusNotActivated
			} else {
				if device.LastOnlineTime == nil {
					targetStatus = consts.DeviceStatusOffline
				} else {
					if time.Since(device.LastOnlineTime.Time) > time.Duration(consts.DeviceOfflineTime)*time.Second {
						targetStatus = consts.DeviceStatusOffline
					} else {
						targetStatus = consts.DeviceStatusOnline
					}
				}
			}

			// 检查设备状态是否需要更新
			if device.Status != int(targetStatus) {
				_, err = dao.Device.Ctx(ctx).Data(do.Device{
					Status: targetStatus,
				}).WherePri(device.Id).Update()

				if err != nil {
					g.Log().Line().Error(ctx, "更新设备状态失败", err.Error())
					continue
				}

				g.Log().Line().Infof(ctx, "设备 %s 状态从 %s 更新为 %s", device.DeviceKey, consts.GetDeviceStatusStr(int32(device.Status)), consts.GetDeviceStatusStr(targetStatus))
			}
		}

		// 继续查询下一页
		page++
	}

	return
}

// 检查命令执行超时
func (s *sDevApi) CommandTimeoutCheck(ctx context.Context) {
	// 分页查询正在执行的命令
	pageSize := 100
	page := 1

	for {
		var commands []entity.CommandRecord
		err := dao.CommandRecord.Ctx(ctx).
			Where("status", consts.Command_status_executing).
			Page(page, pageSize).
			Order("id ASC").
			Scan(&commands)

		if err != nil {
			g.Log().Line().Error(ctx, "查询命令记录失败", err.Error())
			return
		}

		// 没有更多数据，退出循环
		if len(commands) == 0 {
			break
		}

		for _, command := range commands {

			if command.CreateTime == nil {
				g.Log().Line().Warningf(ctx, "命令记录 %d 执行时间为空，跳过超时检查", command.Id)
				continue
			}
			if time.Since(command.CreateTime.Time) > time.Duration(command.TimeoutSeconds)*time.Second {
				// 更新命令状态为超时
				_, err = dao.CommandRecord.Ctx(ctx).Data(do.CommandRecord{
					Status: int(consts.Command_status_timeout),
				}).WherePri(command.Id).Update()

				if err != nil {
					g.Log().Line().Error(ctx, "更新命令状态失败", err.Error())
					continue
				}

				g.Log().Line().Infof(ctx, "命令 %d 超时，状态更新为超时", command.Id)
			}
		}

		// 继续查询下一页
		page++
	}

	return
}

// ParseExcelFile  解析Excel文件并返回设备列表
func (s *sDevApi) ParseExcelFile(ctx context.Context, fileBytes []byte) (devices []do.Device, err error) {
	// 从字节数据创建Excel文件
	reader := bytes.NewReader(fileBytes)
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}
	defer func(f *excelize.File) {
		err := f.Close()
		if err != nil {
			g.Log().Line().Error(ctx, "关闭Excel文件失败:", err.Error())
		}
	}(f)

	// 获取第一个工作表
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, gerror.New("Excel文件中没有工作表")
	}

	// 读取第一个工作表的所有行
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}

	// 跳过标题行，从第二行开始读取数据
	if len(rows) <= 1 {
		return nil, gerror.New("Excel文件中没有数据")
	}

	// 处理每一行数据
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		// 确保行有足够的列数
		if len(row) < 2 {
			continue
		}

		// 从Excel行创建设备对象
		// 假设第一列是设备名称，第二列是设备序列号
		device := do.Device{
			DeviceDesc:   row[0],
			DeviceSecret: row[1],
		}

		devices = append(devices, device)
	}

	if len(devices) == 0 {
		return nil, gerror.New("没有有效的设备数据")
	}

	return devices, nil
}

// ParseDeviceMessage 解析设备消息
func (s *sDevApi) ParseDeviceMessage(ctx context.Context, deviceKey string, js string) (err error) {
	// 解析JSON数据
	var message model.DeviceMessage
	err = json.Unmarshal([]byte(js), &message)
	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return gerror.New("解析设备消息JSON失败")
	}
	// 检查消息类型
	if message.MessageType == consts.Edge_command_up {
		_ = s.commandMsgHandle(ctx, deviceKey, message.Content)
	} else if message.MessageType == consts.Edge_heartbeat_up {
		_ = s.heartbeatHandle(ctx, deviceKey, message.Content)
	} else if message.MessageType == consts.Edge_ota_up {
		_ = s.otaMsgHandle(ctx, deviceKey, message.Content)
	} else {
		g.Log().Line().Warning(ctx, "未知的设备消息类型:", message.MessageType)
	}

	return nil
}

// commandMsgHandle 记录设备命令响应
func (s *sDevApi) commandMsgHandle(ctx context.Context, _deviceKey string, content interface{}) (err error) {
	// 解析数据
	var commandRes model.DeviceMessageCommandResponse
	if err = utility.ParseContent(ctx, content, &commandRes, "解析设备命令响应失败"); err != nil {
		return err
	}

	// 是否需要校验设备与执行id的一致性

	//todo 向前端推送命令执行日志

	_, err = dao.CommandLog.Ctx(ctx).Data(do.CommandLog{
		CommandId: commandRes.CommandId,
		Level:     commandRes.Level,
		Content:   commandRes.Content,
	}).Insert()
	if err != nil {
		g.Log().Line().Error(ctx, "插入命令日志失败:", err.Error())
		return err
	}

	// 更新命令记录状态
	var targetStatus int
	if commandRes.Status == consts.Command_status_completed {
		targetStatus = int(consts.Command_status_completed)
	} else if commandRes.Status == consts.Command_status_failed {
		targetStatus = int(consts.Command_status_failed)
	} else if commandRes.Status == consts.Command_status_executing {
		targetStatus = int(consts.Command_status_executing)
	} else {
		g.Log().Line().Error(ctx, "未知的命令状态:", commandRes.Status)
		return nil
	}

	//todo 向前端推送命令执行结果

	_, err = dao.CommandRecord.Ctx(ctx).Data(do.CommandRecord{Status: targetStatus}).WherePri(commandRes.CommandId).Update()
	if err != nil {
		g.Log().Line().Error(ctx, "更新命令记录状态失败:", err.Error())
		return err
	}

	return nil
}

// heartbeatHandle  记录设备心跳
func (s *sDevApi) heartbeatHandle(ctx context.Context, _deviceKey string, content interface{}) (err error) {
	// 解析数据
	var devHeartbeat model.DeviceMessageHeartbeat
	if err = utility.ParseContent(ctx, content, &devHeartbeat, "解析设备心跳失败"); err != nil {
		return err
	}

	// 更新设备最后在线时间、版本、CPU、内存、磁盘等信息
	_, err = dao.Device.Ctx(ctx).Data(do.Device{
		LastOnlineTime:  gtime.Now(),
		HardwareVersion: devHeartbeat.Version,
		Cpu:             devHeartbeat.Cpu,
		Mem:             devHeartbeat.Mem,
		MemTotal:        devHeartbeat.MemTotal,
		Disk:            devHeartbeat.Disk,
		DiskTotal:       devHeartbeat.DiskTotal,
	}).Where(do.Device{DeviceKey: _deviceKey}).Update()
	if err != nil {
		g.Log().Line().Error(ctx, "更新设备心跳信息失败:", err.Error())
		return err
	}

	return nil
}

func (s *sDevApi) otaMsgHandle(ctx context.Context, _deviceKey string, content interface{}) (err error) {
	// 解析数据
	var otaRes model.DeviceMessageOtaResponse
	if err = utility.ParseContent(ctx, content, &otaRes, "解析设备OTA响应失败"); err != nil {
		return err
	}

	// 查找对应的升级记录
	var record entity.UpgradeRecord
	err = dao.UpgradeRecord.Ctx(ctx).Where(do.UpgradeRecord{
		Id: otaRes.Id,
	}).OrderDesc("id").Limit(1).Scan(&record)

	if err != nil {
		g.Log().Line().Error(ctx, "查询升级记录失败:", err.Error())
		return err
	}

	// 更新升级状态
	status := otaRes.Status
	logContent := otaRes.Content
	logLevel := otaRes.Level

	// 更新升级记录
	_, err = dao.UpgradeRecord.Ctx(ctx).WherePri(record.Id).Data(do.UpgradeRecord{UpgradeStatus: status}).Update()
	if err != nil {
		g.Log().Line().Error(ctx, "更新升级记录失败:", err.Error())
		return err
	}

	// 记录日志
	if logContent != "" {
		_, err = dao.UpgradeLog.Ctx(ctx).Data(do.UpgradeLog{
			RecordId: record.Id,
			Content:  logContent,
			Level:    logLevel,
		}).Insert()
		if err != nil {
			g.Log().Line().Error(ctx, "创建升级日志失败:", err.Error())
		}
	}

	return
}
