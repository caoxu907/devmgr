package ota

import (
	"context"
	filemgr "devmgr/api/filemgr/v1"
	"devmgr/internal/consts"
	"devmgr/internal/dao"
	"devmgr/internal/externalapi"
	"devmgr/internal/model"
	"devmgr/internal/model/do"
	"devmgr/internal/model/entity"
	"devmgr/internal/service"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func New() service.IOta {
	return &sOta{}
}

type sOta struct {
}

// CheckUpgradeTaskStatus 定时检查升级任务状态
func (s *sOta) CheckUpgradeTaskStatus(ctx context.Context) {
	// 检查正在执行的任务状态
	s.checkExecutingOtaTaskStatus(ctx)

	// 检查未开始的任务状态
	s.checkPendingOtaTaskStatus(ctx)

	//检查待推送的升级任务
	s.checkPendingPushTasks(ctx)
}

// 检查未开始的任务状态
func (s *sOta) checkPendingOtaTaskStatus(ctx context.Context) {
	// 处理待处理的批次（灰度发布自动创建的批次）
	var pendingBatches []entity.UpgradeBatch
	err := dao.UpgradeBatch.Ctx(ctx).Where(do.UpgradeBatch{
		Batchstatus: consts.Ota_batch_status_pending,
	}).Scan(&pendingBatches)
	if err != nil {
		g.Log().Line().Error(ctx, "查询待处理批次失败:", err.Error())
		return
	}

	for _, batch := range pendingBatches {
		// 更新状态为处理中
		_, err := dao.UpgradeBatch.Ctx(ctx).WherePri(batch.Id).Data(do.UpgradeBatch{
			Batchstatus: consts.Ota_batch_status_executing,
		}).Update()
		if err != nil {
			g.Log().Line().Error(ctx, "更新批次状态失败:", err.Error())
			continue
		}

		// 获取升级包信息
		var pkg entity.UpgradePackage
		err = dao.UpgradePackage.Ctx(ctx).WherePri(batch.PackageId).Scan(&pkg)
		if err != nil {
			g.Log().Line().Error(ctx, "查询升级包失败:", err.Error())
			s.markBatchFailed(ctx, batch.Id, "查询升级包失败: "+err.Error())
			continue
		}

		// 查询要升级的设备
		var targetDevices []entity.Device
		if batch.ScopeType == int(consts.Ota_scope_type_all) {
			// 查询所有符合条件的设备
			err = dao.Device.Ctx(ctx).Where(do.Device{
				ProductId: pkg.ProductId,
				Deleted:   false,
				Status:    int(consts.DeviceStatusOnline), // 只升级在线设备
			}).Scan(&targetDevices)
			if err != nil {
				g.Log().Line().Error(ctx, "查询设备列表失败:", err.Error())
				s.markBatchFailed(ctx, batch.Id, "查询设备信息失败: "+err.Error())

				return
			}
		} else if batch.ScopeType == int(consts.Ota_scope_type_manual) {
			// 先将 JSON 字符串解析为设备 ID 数组
			var deviceIds []int64
			err = json.Unmarshal([]byte(batch.ScopeDevices), &deviceIds)
			if err != nil {
				g.Log().Line().Error(ctx, "解析设备ID列表失败:", err.Error())
				s.markBatchFailed(ctx, batch.Id, "解析设备ID列表失败: "+err.Error())
				return
			}

			// 使用解析后的数组进行查询
			err = dao.Device.Ctx(ctx).WhereIn("id", deviceIds).
				Where(do.Device{
					Deleted: false,
					Status:  int(consts.DeviceStatusOnline),
				}).Scan(&targetDevices)
			if err != nil {
				g.Log().Line().Error(ctx, "查询设备列表失败:", err.Error())
				s.markBatchFailed(ctx, batch.Id, "查询设备信息失败: "+err.Error())
				return
			}
		}
		totalDevices := len(targetDevices)
		if totalDevices == 0 {
			g.Log().Line().Info(ctx, "没有符合条件的设备需要升级，批次ID:", batch.Id)
			// 更新批次状态为已完成
			_, err := dao.UpgradeBatch.Ctx(ctx).WherePri(batch.Id).Data(do.UpgradeBatch{
				Batchstatus:    consts.Ota_batch_status_completed,
				TotalDevices:   0,
				SuccessDevices: 0,
				FailureDevices: 0,
				PendingDevices: 0,
				CompleteTime:   gtime.New(),
			}).Update()
			if err != nil {
				g.Log().Line().Error(ctx, "更新批次状态失败:", err.Error())
			}
			continue
		}

		// 灰度发布处理
		var devicesToUpgrade []entity.Device
		if batch.StrategyType == int(consts.Ota_strategy_type_gray) {
			// 计算灰度数量
			grayCount := int(float64(totalDevices) * float64(batch.GrayPercentage) / 100.0)
			if grayCount < 1 {
				grayCount = 1 // 至少选择一个设备
			}

			// 随机选择指定比例的设备
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(targetDevices), func(i, j int) {
				targetDevices[i], targetDevices[j] = targetDevices[j], targetDevices[i]
			})

			// 取前N个设备作为灰度设备
			if grayCount > len(targetDevices) {
				grayCount = len(targetDevices)
			}
			devicesToUpgrade = targetDevices[:grayCount]

			g.Log().Line().Info(ctx, "灰度升级设备数量:", grayCount, "总设备数:", totalDevices)
		} else {
			// 非灰度策略，所有设备都升级
			devicesToUpgrade = targetDevices
		}

		// 处理批次
		err = s.processNewBatch(ctx, batch.Id, devicesToUpgrade, pkg)
		if err != nil {
			g.Log().Line().Error(ctx, "处理批次失败:", err.Error())
			s.markBatchFailed(ctx, batch.Id, "处理批次失败: "+err.Error())
			continue
		}
	}
}

// 检查正在执行的任务状态
func (s *sOta) checkExecutingOtaTaskStatus(ctx context.Context) {
	// 获取所有正在执行的批次
	var batches []entity.UpgradeBatch
	err := dao.UpgradeBatch.Ctx(ctx).Where(do.UpgradeBatch{
		Batchstatus: consts.Ota_batch_status_executing,
	}).Scan(&batches)
	if err != nil {
		g.Log().Line().Error(ctx, "查询升级批次失败:", err.Error())
		return
	}
	for _, batch := range batches {
		// 检查每个批次中的记录
		s.checkBatchRecords(ctx, batch)

		// 检查批次是否已完成
		completed := s.updateBatchStatus(ctx, batch.Id)

		// 如果批次已完成
		if completed {
			// 如果是灰度发布且灰度阶段未完成，检查是否需要扩大发布范围
			if batch.StrategyType == int(consts.Ota_strategy_type_gray) && !batch.GrayCompleted {
				err := s.checkGrayDeployment(ctx, batch)
				if err != nil {
					g.Log().Line().Error(ctx, "检查灰度发布状态失败:", err.Error())
					continue
				}
				// checkGrayDeployment 中会处理灰度扩展并设置 GrayCompleted
				// 如果灰度扩展，批次会保持执行中状态，等待剩余设备升级完成
			} else {
				// 灰度阶段已完成或非灰度批次，直接更新批次状态为已完成
				_, err := dao.UpgradeBatch.Ctx(ctx).WherePri(batch.Id).Data(do.UpgradeBatch{
					Batchstatus:  consts.Ota_batch_status_completed,
					CompleteTime: gtime.New(),
				}).Update()
				if err != nil {
					g.Log().Line().Error(ctx, "更新批次状态失败:", err.Error())
				}
			}
		}
	}
}

// 检查批次中的记录
func (s *sOta) checkBatchRecords(ctx context.Context, batch entity.UpgradeBatch) {
	// 获取批次中的所有记录
	var records []entity.UpgradeRecord
	err := dao.UpgradeRecord.Ctx(ctx).Where(do.UpgradeRecord{
		BatchId:  batch.Id,
		IsActive: true,
	}).Scan(&records)
	if err != nil {
		g.Log().Line().Error(ctx, "查询升级记录失败:", err.Error())
		return
	}

	mqClient := service.MQClient()
	if mqClient == nil {
		g.Log().Line().Error(ctx, "MQClient未初始化")
		return
	}

	// 检查每条记录
	for _, record := range records {
		// 检查超时的记录
		if record.UpgradeStatus == int(consts.Ota_record_status_pushed) ||
			record.UpgradeStatus == int(consts.Ota_record_status_upgrading) {

			// 如果开始时间不为空，检查是否超时
			if record.StartTime != nil {
				startTime := record.StartTime.Time
				currentTime := time.Now()

				// 超过30分钟视为超时
				if currentTime.Sub(startTime) > 30*time.Minute {
					// 尝试重试
					if record.RetryTimes < batch.MaxRetryCount {
						s.retryUpgrade(ctx, record, batch, mqClient)
					} else {
						// 重试次数已达上限，标记为失败
						s.markUpgradeFailed(ctx, record.Id, "升级超时，重试次数已达上限")
					}
				}
			} else if record.UpgradeStatus == int(consts.Ota_record_status_pushed) {
				// 如果推送后一段时间还未开始下载，考虑重试
				currentTime := time.Now()
				createTime := record.CreateTime.Time

				// 推送后5分钟还未开始，尝试重试
				if currentTime.Sub(createTime) > 5*time.Minute && record.RetryTimes < batch.MaxRetryCount {
					s.retryUpgrade(ctx, record, batch, mqClient)
				}
			}
		}
	}
}

// 重试升级
func (s *sOta) retryUpgrade(ctx context.Context, record entity.UpgradeRecord, batch entity.UpgradeBatch, mqClient service.IMQClient) {
	// 获取设备信息
	var device entity.Device
	err := dao.Device.Ctx(ctx).WherePri(record.DeviceId).Scan(&device)
	if err != nil {
		g.Log().Line().Error(ctx, "查询设备信息失败:", err.Error())
		return
	}

	// 检查设备是否在线
	if device.Status != int(consts.DeviceStatusOnline) {
		s.markUpgradeFailed(ctx, record.Id, "设备不在线，无法重试升级")
		return
	}

	// 获取升级包信息
	var pkg entity.UpgradePackage
	err = dao.UpgradePackage.Ctx(ctx).WherePri(batch.PackageId).Scan(&pkg)
	if err != nil {
		g.Log().Line().Error(ctx, "查询升级包失败:", err.Error())
		return
	}

	// 获取下载URL
	fileClient, err := externalapi.GetFilemgrClient()
	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, consts.Grpc_Timeout*time.Second)
	defer cancel()
	presignDownloadRes, err := fileClient.GeneratePresignedDownloadURL(timeoutCtx, &filemgr.PresignDownloadReq{
		Bucket:   consts.File_bucket,
		Key:      pkg.PackagePath,
		Expiry:   consts.File_download_expiry,
		Intranet: false,
	})
	if err != nil {
		g.Log().Line().Error(ctx, "获取下载URL失败:", err.Error())
		return
	}

	// 构造升级命令
	msg := model.DeviceMessage{
		MessageType: consts.Edge_ota_down,
		Content: model.DeviceMessageOtaSend{
			MessageType:    consts.Ota_command_type_upgrade,
			Id:             int32(record.Id),
			Url:            presignDownloadRes.Url,
			PackageVersion: pkg.Version,
			PackageHash:    pkg.PackageHash,
			PackageSize:    int32(pkg.PackageSize),
		},
	}

	msgData, err := json.Marshal(msg)
	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return
	}

	// 更新重试次数
	_, err = dao.UpgradeRecord.Ctx(ctx).WherePri(record.Id).Data(do.UpgradeRecord{
		RetryTimes: record.RetryTimes + 1,
	}).Update()
	if err != nil {
		g.Log().Line().Error(ctx, "更新重试次数失败:", err.Error())
		return
	}

	// 记录重试日志
	_, err = dao.UpgradeLog.Ctx(ctx).Data(do.UpgradeLog{
		RecordId: record.Id,
		Content:  fmt.Sprintf("第%d次重试升级任务", record.RetryTimes+1),
		Level:    consts.Msg_level_warn,
	}).Insert()
	if err != nil {
		g.Log().Line().Error(ctx, "创建升级日志失败:", err.Error())
	}

	// 推送升级命令
	g.Log().Line().Info(ctx, "重试推送升级消息给设备:", device.DeviceKey)
	err = mqClient.Publish(consts.Kafka_topic_notify_down, []byte(device.DeviceKey), msgData)
	if err != nil {
		g.Log().Line().Error(ctx, "推送消息失败:", err.Error(), "设备:", device.DeviceKey)

		// 记录推送失败日志
		_, _ = dao.UpgradeLog.Ctx(ctx).Data(do.UpgradeLog{
			RecordId: record.Id,
			Content:  "消息推送失败: " + err.Error(),
			Level:    consts.Msg_level_error,
		}).Insert()
	}
}

// 更新重试次数+1，检查是否超过最大重试次数，超过则标记失败
func (s *sOta) checkRetryTimes(ctx context.Context, record entity.UpgradeRecord, maxRetryCount int) {
	if record.RetryTimes < maxRetryCount {
		_, err := dao.UpgradeRecord.Ctx(ctx).WherePri(record.Id).Data(do.UpgradeRecord{
			RetryTimes: record.RetryTimes + 1,
		}).Update()
		if err != nil {
			g.Log().Line().Error(ctx, "更新升级状态失败:", err.Error())
			return
		}
	} else {
		s.markUpgradeFailed(ctx, record.Id, "重试次数已达上限")
	}
	return
}

// 标记升级失败
func (s *sOta) markUpgradeFailed(ctx context.Context, recordId int, reason string) {
	// 更新记录状态
	_, err := dao.UpgradeRecord.Ctx(ctx).WherePri(recordId).Data(do.UpgradeRecord{
		UpgradeStatus: consts.Ota_record_status_failed,
	}).Update()
	if err != nil {
		g.Log().Line().Error(ctx, "更新升级状态失败:", err.Error())
		return
	}

	// 记录失败日志
	_, err = dao.UpgradeLog.Ctx(ctx).Data(do.UpgradeLog{
		RecordId: recordId,
		Content:  reason,
		Level:    consts.Msg_level_error,
	}).Insert()
	if err != nil {
		g.Log().Line().Error(ctx, "创建升级日志失败:", err.Error())
	}
}

// 更新批次状态
func (s *sOta) updateBatchStatus(ctx context.Context, batchId int) bool {

	// 获取批次中的记录统计
	var total, pending, success, failed int
	total, err := dao.UpgradeRecord.Ctx(ctx).Where(do.UpgradeRecord{BatchId: batchId}).Count()
	if err != nil {
		g.Log().Line().Error(ctx, "统计升级记录失败:", err.Error())
		return false
	}

	pending, err = dao.UpgradeRecord.Ctx(ctx).Where(do.UpgradeRecord{
		BatchId: batchId,
		UpgradeStatus: []int64{
			int64(consts.Ota_record_status_pending),
			int64(consts.Ota_record_status_to_be_pushed),
			int64(consts.Ota_record_status_pushed),
			int64(consts.Ota_record_status_upgrading),
		},
	}).Count()
	if err != nil {
		g.Log().Line().Error(ctx, "统计待升级记录失败:", err.Error())
		return false
	}

	success, err = dao.UpgradeRecord.Ctx(ctx).Where(do.UpgradeRecord{
		BatchId:       batchId,
		UpgradeStatus: consts.Ota_record_status_success,
	}).Count()
	if err != nil {
		g.Log().Line().Error(ctx, "统计成功升级记录失败:", err.Error())
		return false
	}

	failed, err = dao.UpgradeRecord.Ctx(ctx).Where(do.UpgradeRecord{
		BatchId: batchId,
		UpgradeStatus: []int64{
			int64(consts.Ota_record_status_failed),
			int64(consts.Ota_record_status_canceled), // 失败和取消都算作失败
		},
	}).Count()
	if err != nil {
		g.Log().Line().Error(ctx, "统计失败升级记录失败:", err.Error())
		return false
	}

	// 更新批次状态数据
	updateData := do.UpgradeBatch{
		TotalDevices:   total,
		SuccessDevices: success,
		FailureDevices: failed,
		PendingDevices: pending,
	}

	// 如果所有设备都已完成升级，标记批次为已完成
	var completed bool = false
	if pending == 0 {
		completed = true
	}

	// 更新批次数据
	_, err = dao.UpgradeBatch.Ctx(ctx).WherePri(batchId).Data(updateData).Update()
	if err != nil {
		g.Log().Line().Error(ctx, "更新批次状态失败:", err.Error())
	}

	return completed
}

// 检查灰度发布状态并决策是否扩大发布范围
func (s *sOta) checkGrayDeployment(ctx context.Context, batch entity.UpgradeBatch) (err error) {
	// 计算灰度成功率
	successRate := float64(batch.SuccessDevices) / float64(batch.TotalDevices) * 100

	// 如果成功率达到阈值，扩大发布范围
	if successRate >= float64(batch.GraySuccessThreshold) {
		// 处理剩余设备
		err = s.upgradeRemainingDevices(ctx, batch)
		if err != nil {
			g.Log().Line().Error(ctx, "扩大发布范围失败:", err.Error())
			return err
		}

		// 标记灰度阶段已完成
		_, err := dao.UpgradeBatch.Ctx(ctx).WherePri(batch.Id).Data(do.UpgradeBatch{
			GrayCompleted: true,
			// 保持批次状态为执行中，等待剩余设备升级完成
			Batchstatus: consts.Ota_batch_status_executing,
		}).Update()
		if err != nil {
			g.Log().Line().Error(ctx, "更新批次灰度状态失败:", err.Error())
			return err
		}

		return nil
	} else {
		// 成功率不达标，记录并停止后续发布
		g.Log().Line().Warning(ctx, "灰度发布成功率不达标，停止后续发布。批次ID:", batch.Id,
			"成功率:", successRate, "阈值:", batch.GraySuccessThreshold)

		// 将批次标记为已完成，不处理剩余设备
		_, err := dao.UpgradeBatch.Ctx(ctx).WherePri(batch.Id).Data(do.UpgradeBatch{
			Batchstatus:   consts.Ota_batch_status_failed,
			GrayCompleted: true, // 虽然失败但灰度阶段已完成
			CompleteTime:  gtime.New(),
		}).Update()
		if err != nil {
			g.Log().Line().Error(ctx, "更新批次状态失败:", err.Error())
		}
	}

	return nil
}

// 升级剩余设备
func (s *sOta) upgradeRemainingDevices(ctx context.Context, batch entity.UpgradeBatch) (err error) {
	// 获取升级包信息
	var pkg entity.UpgradePackage
	err = dao.UpgradePackage.Ctx(ctx).WherePri(batch.PackageId).Scan(&pkg)
	if err != nil {
		g.Log().Line().Error(ctx, "查询升级包失败:", err.Error())
		return err
	}

	// 获取已升级的设备ID列表
	var upgradedDeviceIds []int64
	err = dao.UpgradeRecord.Ctx(ctx).
		Where(do.UpgradeRecord{BatchId: batch.Id}).
		Fields("device_id").
		Scan(&upgradedDeviceIds)
	if err != nil {
		g.Log().Line().Error(ctx, "查询已升级设备失败:", err.Error())
		return err
	}

	// 查询剩余符合条件但尚未升级的设备
	var remainingDevices []entity.Device

	// 如果是产品下所有设备
	if batch.ScopeType == int(consts.Ota_scope_type_all) {
		err = dao.Device.Ctx(ctx).
			Where(do.Device{
				ProductId: pkg.ProductId,
				Deleted:   false,
				Status:    int(consts.DeviceStatusOnline),
			}).
			WhereNotIn("id", upgradedDeviceIds).
			Scan(&remainingDevices)
		if err != nil {
			g.Log().Line().Error(ctx, "查询剩余设备失败:", err.Error())
			return err
		}
	} else if batch.ScopeType == int(consts.Ota_scope_type_manual) {
		// 如果是手动选择的设备，解析原始选择的设备
		var originalScopeDevices []int64
		err = json.Unmarshal([]byte(batch.ScopeDevices), &originalScopeDevices)
		if err != nil {
			g.Log().Line().Error(ctx, "解析设备ID列表失败:", err.Error())
			return err
		}

		// 过滤出未升级的设备
		var remainingIds []int64
		for _, id := range originalScopeDevices {
			if !s.containsInt64(upgradedDeviceIds, id) {
				remainingIds = append(remainingIds, id)
			}
		}

		if len(remainingIds) > 0 {
			err = dao.Device.Ctx(ctx).
				WhereIn("id", remainingIds).
				Where(do.Device{
					Deleted: false,
					Status:  int(consts.DeviceStatusOnline),
				}).
				Scan(&remainingDevices)
			if err != nil {
				g.Log().Line().Error(ctx, "查询剩余设备失败:", err.Error())
				return err
			}
		}
	}

	if len(remainingDevices) == 0 {
		g.Log().Line().Info(ctx, "没有剩余设备需要升级，灰度发布已完成。批次ID:", batch.Id)

		// 标记批次为已完成
		_, err := dao.UpgradeBatch.Ctx(ctx).WherePri(batch.Id).Data(do.UpgradeBatch{
			Batchstatus:  consts.Ota_batch_status_completed,
			CompleteTime: gtime.New(),
		}).Update()
		if err != nil {
			g.Log().Line().Error(ctx, "更新批次状态失败:", err.Error())
		}

		return nil
	}

	// 更新批次的设备总数(增加剩余设备数量)
	newTotalDevices := int64(batch.TotalDevices) + int64(len(remainingDevices))
	newPendingDevices := int64(batch.PendingDevices) + int64(len(remainingDevices))

	_, err = dao.UpgradeBatch.Ctx(ctx).WherePri(batch.Id).Data(do.UpgradeBatch{
		TotalDevices:   newTotalDevices,
		PendingDevices: newPendingDevices,
		// 保持批次状态为执行中
		Batchstatus: consts.Ota_batch_status_executing,
	}).Update()
	if err != nil {
		g.Log().Line().Error(ctx, "更新批次设备数量失败:", err.Error())
		return err
	}

	// 直接处理剩余设备的升级任务
	err = s.processDevicesUpgrade(ctx, batch.Id, remainingDevices, pkg, false)
	if err != nil {
		g.Log().Line().Error(ctx, "处理剩余设备升级失败:", err.Error())
		return err
	}

	g.Log().Line().Info(ctx, "灰度发布成功，已开始处理剩余设备。批次ID:", batch.Id,
		"剩余设备数:", len(remainingDevices))

	return nil
}

func (s *sOta) containsInt64(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// 标记批次处理失败
func (s *sOta) markBatchFailed(ctx context.Context, batchId int, reason string) {
	// 更新批次状态为失败，并记录失败原因
	g.Log().Line().Info(ctx, "标记批次为失败，批次ID:", batchId, "原因:", reason)
	_, err := dao.UpgradeBatch.Ctx(ctx).WherePri(batchId).Data(do.UpgradeBatch{
		Batchstatus: consts.Ota_batch_status_failed,
		//FailReason:   reason,
		CompleteTime: gtime.New(), // 记录完成时间
	}).Update()
	if err != nil {
		g.Log().Line().Error(ctx, "更新批次状态失败:", err.Error())
		return
	}

	// 将该批次下所有未完成的设备升级记录标记为取消
	_, err = dao.UpgradeRecord.Ctx(ctx).
		Where(do.UpgradeRecord{
			BatchId:  batchId,
			IsActive: true,
		}).
		WhereIn("upgrade_status", []int32{
			consts.Ota_record_status_pending,
			consts.Ota_record_status_to_be_pushed,
			consts.Ota_record_status_pushed,
			consts.Ota_record_status_upgrading,
		}).
		Data(do.UpgradeRecord{
			UpgradeStatus: consts.Ota_record_status_canceled,
		}).
		Update()
}

func (s *sOta) processNewBatch(ctx context.Context, batchId int, devices []entity.Device,
	pkg entity.UpgradePackage) error {

	g.Log().Line().Info(ctx, "处理新批次，设备数量:", len(devices), "批次ID:", batchId)

	// 更新总设备数和待升级设备数
	totalDevices := len(devices)
	_, err := dao.UpgradeBatch.Ctx(ctx).WherePri(batchId).Data(do.UpgradeBatch{
		TotalDevices:   totalDevices,
		PendingDevices: totalDevices,
		Batchstatus:    consts.Ota_batch_status_executing,
	}).Update()
	if err != nil {
		g.Log().Line().Error(ctx, "更新批次设备数量失败:", err.Error())
		return fmt.Errorf("更新批次设备数量失败: %w", err)
	}

	// 判断是否为灰度发布
	isGray := false
	var batch entity.UpgradeBatch
	err = dao.UpgradeBatch.Ctx(ctx).WherePri(batchId).Scan(&batch)
	if err == nil && batch.StrategyType == int(consts.Ota_strategy_type_gray) {
		isGray = true
	}

	// 处理设备升级
	err = s.processDevicesUpgrade(ctx, batchId, devices, pkg, isGray)
	if err != nil {
		g.Log().Line().Error(ctx, "处理设备升级失败:", err.Error())
		return fmt.Errorf("处理设备升级失败: %w", err)
	}

	return nil
}

// 处理设备升级
// 处理设备升级，分批次推送
func (s *sOta) processDevicesUpgrade(ctx context.Context, batchId int, devices []entity.Device,
	pkg entity.UpgradePackage, isGray bool) error {

	// 获取批次信息
	var batch entity.UpgradeBatch
	err := dao.UpgradeBatch.Ctx(ctx).WherePri(batchId).Scan(&batch)
	if err != nil {
		g.Log().Line().Error(ctx, "查询批次信息失败:", err.Error())
		return err
	}

	// 设置默认值
	maxConcurrent := 10 // 默认并发数
	if batch.MaxConcurrent > 0 {
		maxConcurrent = batch.MaxConcurrent
	}

	// 随机延迟范围
	randomDelay := 60 // 默认最大延迟60秒
	if batch.RandomDelay > 0 {
		randomDelay = batch.RandomDelay
	}

	// 计算分组数量
	totalDevices := len(devices)
	groupCount := (totalDevices + maxConcurrent - 1) / maxConcurrent

	g.Log().Line().Info(ctx, "设备升级分组处理，设备总数:", totalDevices,
		"并发数:", maxConcurrent,
		"分组数:", groupCount)

	// 为每个设备创建记录，设置不同的预计开始时间
	for i, device := range devices {
		// 计算所属组和组内索引
		groupIndex := i / maxConcurrent

		// 计算延迟时间
		// 1. 组间基础延迟
		baseDelay := groupIndex * batch.GroupInterval

		// 2. 组内随机延迟 (0-randomDelay秒)
		rand.Seed(time.Now().UnixNano() + int64(i))
		groupRandomDelay := rand.Intn(randomDelay)

		totalDelay := baseDelay + groupRandomDelay

		// 计划执行时间
		planTime := time.Now().Add(time.Duration(totalDelay) * time.Second)
		planTimeGTime := gtime.NewFromTime(planTime)

		// 创建升级记录
		recordData := do.UpgradeRecord{
			BatchId:       batchId,
			DeviceId:      device.Id,
			RetryTimes:    0,
			PlanTime:      planTimeGTime,
			UpgradeStatus: consts.Ota_record_status_to_be_pushed, // 待推送状态
			FromVersion:   device.Version,
			ToVersion:     pkg.Version,
			IsActive:      true,
		}

		recordResult, err := dao.UpgradeRecord.Ctx(ctx).Data(recordData).Insert()
		if err != nil {
			g.Log().Line().Error(ctx, "创建升级记录失败:", err.Error(), "设备ID:", device.Id)
			continue
		}

		recordId, err := recordResult.LastInsertId()
		if err != nil {
			g.Log().Line().Error(ctx, "获取记录ID失败:", err.Error())
			continue
		}

		// 记录日志
		logContent := fmt.Sprintf("设备已排队等待升级，计划时间: %s", planTimeGTime.String())
		if isGray {
			logContent = fmt.Sprintf("设备已排队等待灰度升级，计划时间: %s", planTimeGTime.String())
		}

		_, err = dao.UpgradeLog.Ctx(ctx).Data(do.UpgradeLog{
			RecordId: recordId,
			Content:  logContent,
			Level:    consts.Msg_level_info,
		}).Insert()
		if err != nil {
			g.Log().Line().Error(ctx, "创建升级日志失败:", err.Error())
		}
	}

	g.Log().Line().Info(ctx, "设备升级任务已创建，等待定时任务处理。批次ID:", batchId)
	return nil
}

// 检查待推送的升级任务
func (s *sOta) checkPendingPushTasks(ctx context.Context) {
	// 获取当前时间
	now := time.Now()

	// 获取计划时间已到但尚未推送的升级记录
	var records []entity.UpgradeRecord
	err := dao.UpgradeRecord.Ctx(ctx).
		Where(do.UpgradeRecord{
			UpgradeStatus: consts.Ota_record_status_to_be_pushed,
			IsActive:      true,
		}).
		WhereLTE("plan_time", gtime.NewFromTime(now)).
		Limit(50). // 每次处理50条
		Scan(&records)

	if err != nil {
		g.Log().Line().Error(ctx, "查询待推送升级记录失败:", err.Error())
		return
	}

	if len(records) == 0 {
		return // 没有待处理的记录
	}

	g.Log().Line().Info(ctx, "开始处理待推送升级任务，数量:", len(records))

	// 获取当前升级中的设备数量
	upgradingCount, err := dao.UpgradeRecord.Ctx(ctx).
		Where(do.UpgradeRecord{
			IsActive: true,
		}).
		WhereIn("upgrade_status", []int32{
			consts.Ota_record_status_pushed,
			consts.Ota_record_status_upgrading,
		}).
		Count()

	if err != nil {
		g.Log().Line().Error(ctx, "查询当前升级中设备数量失败:", err.Error())
		return
	}

	// 设置最大并发数量
	maxConcurrent := 50 // 系统全局最大并发升级数

	// 如果当前已经达到最大并发，暂不处理
	if upgradingCount >= maxConcurrent {
		g.Log().Line().Info(ctx, "当前升级中设备数量已达上限，等待下次处理:", upgradingCount)
		return
	}

	// 计算可以处理的数量
	canProcessCount := maxConcurrent - upgradingCount
	if canProcessCount > len(records) {
		canProcessCount = len(records)
	}

	// 处理可推送的记录
	recordsToProcess := records[:canProcessCount]

	mqClient := service.MQClient()
	if mqClient == nil {
		g.Log().Line().Error(ctx, "MQClient未初始化")
		return
	}

	// 按批次分组处理
	batchMap := make(map[int][]entity.UpgradeRecord)
	for _, record := range recordsToProcess {
		batchMap[record.BatchId] = append(batchMap[record.BatchId], record)
	}

	// 处理每个批次
	for batchId, batchRecords := range batchMap {
		// 获取批次信息
		var batch entity.UpgradeBatch
		err = dao.UpgradeBatch.Ctx(ctx).WherePri(batchId).Scan(&batch)
		if err != nil {
			g.Log().Line().Error(ctx, "查询批次信息失败:", err.Error(), "批次ID:", batchId)
			continue
		}

		// 获取升级包信息
		var pkg entity.UpgradePackage
		err = dao.UpgradePackage.Ctx(ctx).WherePri(batch.PackageId).Scan(&pkg)
		if err != nil {
			g.Log().Line().Error(ctx, "查询升级包失败:", err.Error())
			continue
		}

		// 获取下载URL
		var fileClient filemgr.FileServiceClient
		fileClient, err = externalapi.GetFilemgrClient()
		if err != nil {
			errInfo := "获取文件服务客户端失败: " + err.Error()
			g.Log().Line().Error(ctx, errInfo)
			s.markBatchFailed(ctx, batch.Id, errInfo) //不可恢复错误，标记批次失败
			continue
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, consts.Grpc_Timeout*time.Second)

		var presignDownloadRes *filemgr.PresignDownloadRes
		presignDownloadRes, err = fileClient.GeneratePresignedDownloadURL(timeoutCtx, &filemgr.PresignDownloadReq{
			Bucket:   consts.File_bucket,
			Key:      pkg.PackagePath,
			Expiry:   consts.File_download_expiry,
			Intranet: false,
		})
		cancel()

		if err != nil {
			errInfo := "获取下载URL失败: " + err.Error()
			g.Log().Line().Error(ctx, errInfo)
			s.markBatchFailed(ctx, batch.Id, errInfo) //不可恢复错误，标记批次失败
			continue
		}

		// 推送每个设备
		for _, record := range batchRecords {
			var device entity.Device
			err = dao.Device.Ctx(ctx).WherePri(record.DeviceId).Scan(&device)
			if err != nil {
				g.Log().Line().Error(ctx, "查询设备信息失败:", err.Error())
				continue
			}

			// 检查设备是否在线
			if device.Status != int(consts.DeviceStatusOnline) {
				errorInfo := fmt.Sprintf("设备[%s]不在线，无法推送升级", device.DeviceKey)
				g.Log().Line().Warning(ctx, errorInfo)
				s.markUpgradeFailed(ctx, record.Id, errorInfo)
				continue
			}

			// 构造升级命令
			msg := model.DeviceMessage{
				MessageType: consts.Edge_ota_down,
				Content: model.DeviceMessageOtaSend{
					MessageType:    consts.Ota_command_type_upgrade,
					Id:             int32(record.Id),
					Url:            presignDownloadRes.Url,
					PackageVersion: pkg.Version,
					PackageHash:    pkg.PackageHash,
					PackageSize:    int32(pkg.PackageSize),
				},
			}

			var msgData []byte
			msgData, err = json.Marshal(msg)
			if err != nil {
				g.Log().Line().Error(ctx, err.Error())
				continue
			}

			// 更新记录状态为已推送
			_, err = dao.UpgradeRecord.Ctx(ctx).WherePri(record.Id).Data(do.UpgradeRecord{
				UpgradeStatus: consts.Ota_record_status_pushed,
			}).Update()
			if err != nil {
				g.Log().Line().Error(ctx, "更新记录状态失败:", err.Error())
				continue
			}

			// 记录日志
			isGray := batch.StrategyType == int(consts.Ota_strategy_type_gray) && !batch.GrayCompleted
			logContent := "已推送升级任务(全量发布)"
			if isGray {
				logContent = "已推送升级任务(灰度发布)"
			}

			_, err = dao.UpgradeLog.Ctx(ctx).Data(do.UpgradeLog{
				RecordId: record.Id,
				Content:  logContent,
				Level:    consts.Msg_level_info,
			}).Insert()
			if err != nil {
				g.Log().Line().Error(ctx, "创建升级日志失败:", err.Error())
			}

			// 推送升级命令
			g.Log().Line().Info(ctx, "推送升级消息给设备:", device.DeviceKey)
			err = mqClient.Publish(consts.Kafka_topic_notify_down, []byte(device.DeviceKey), msgData)
			if err != nil {
				g.Log().Line().Error(ctx, "推送消息失败:", err.Error(), "设备:", device.DeviceKey)
				// 记录失败日志但继续处理其他设备
				_, _ = dao.UpgradeLog.Ctx(ctx).Data(do.UpgradeLog{
					RecordId: record.Id,
					Content:  "消息推送失败: " + err.Error(),
					Level:    consts.Msg_level_error,
				}).Insert()
			}
		}
	}
}

// RegisterUpgradeTasksChecker 在应用启动时注册定时任务
func RegisterUpgradeTasksChecker() {
	// 每分钟检查一次升级任务状态
	//gtime.SetInterval(time.Minute, func() {
	//	ctx := gctx.New()
	//	CheckUpgradeTaskStatus(ctx)
	//})
}

//
//func (*Controller) DeviceUpgradeReport(ctx context.Context, req *v1.DeviceUpgradeReportReq) (res *v1.DeviceUpgradeReportRes, err error) {
//	// 查找对应的升级记录
//	var record entity.UpgradeRecord
//	err = dao.UpgradeRecord.Ctx(ctx).Where(do.UpgradeRecord{
//		DeviceId:  req.DeviceId,
//		IsActive:  true,
//		ToVersion: req.TargetVersion,
//	}).OrderDesc("id").Limit(1).Scan(&record)
//
//	if err != nil {
//		g.Log().Line().Error(ctx, "查询升级记录失败:", err.Error())
//		return nil, err
//	}
//
//	if record.Id == 0 {
//		g.Log().Line().Warning(ctx, "未找到对应的升级记录，设备ID:", req.DeviceId)
//		return &v1.DeviceUpgradeReportRes{}, nil
//	}
//
//	// 更新升级状态
//	updateData := do.UpgradeRecord{}
//
//	// 记录日志内容
//	logContent := ""
//	logLevel := consts.Msg_level_info
//
//	switch req.Status {
//	case consts.Upgrade_record_status_downloading:
//		updateData.UpgradeStatus = consts.Upgrade_record_status_downloading
//		if record.StartTime == nil {
//			updateData.StartTime = gtime.New()
//		}
//		logContent = "设备开始下载升级包"
//
//	case consts.Ota_record_status_upgrading:
//		updateData.UpgradeStatus = consts.Ota_record_status_upgrading
//		logContent = "设备开始安装升级包"
//
//	case consts.Ota_record_status_success:
//		updateData.UpgradeStatus = consts.Ota_record_status_success
//		updateData.Duration = int(time.Since(record.StartTime.Time).Seconds())
//		logContent = "设备升级成功"
//
//		// 更新设备版本信息
//		_, err = dao.Device.Ctx(ctx).WherePri(req.DeviceId).Data(do.Device{
//			Version: req.TargetVersion,
//		}).Update()
//		if err != nil {
//			g.Log().Line().Error(ctx, "更新设备版本失败:", err.Error())
//		}
//
//	case consts.Ota_record_status_failed:
//		updateData.UpgradeStatus = consts.Ota_record_status_failed
//		logContent = "设备升级失败: " + req.Message
//		logLevel = consts.Msg_level_error
//	}
//
//	// 更新升级记录
//	if updateData.UpgradeStatus > 0 {
//		_, err = dao.UpgradeRecord.Ctx(ctx).WherePri(record.Id).Data(updateData).Update()
//		if err != nil {
//			g.Log().Line().Error(ctx, "更新升级记录失败:", err.Error())
//			return nil, err
//		}
//	}
//
//	// 记录日志
//	if logContent != "" {
//		_, err = dao.UpgradeLog.Ctx(ctx).Data(do.UpgradeLog{
//			RecordId: record.Id,
//			Content:  logContent,
//			Level:    logLevel,
//		}).Insert()
//		if err != nil {
//			g.Log().Line().Error(ctx, "创建升级日志失败:", err.Error())
//		}
//	}
//
//	// 获取批次信息，检查是否需要更新批次状态
//	updateBatchStatus(ctx, record.BatchId)
//
//	return &v1.DeviceUpgradeReportRes{}, nil
//}
