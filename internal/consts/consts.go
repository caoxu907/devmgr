package consts

// Grpc_Timeout grpc请求超时
const (
	Grpc_Timeout           = 10 //seconds
	Upgrade_max_concurrent = 10 //最大升级并发数
	Upgrade_group_interval = 60 //组间基础延迟（秒）
	Upgrade_random_delay   = 60 //组内随机最大延迟 （秒）
)

// 设备状态类型定义
const (
	DeviceStatusUnKnown      int32 = 0
	DeviceStatusNotActivated int32 = 1
	DeviceStatusOffline      int32 = 2
	DeviceStatusOnline       int32 = 3
)

func GetDeviceStatusStr(dt int32) string {
	switch dt {
	case DeviceStatusNotActivated:
		return "未激活"
	case DeviceStatusOffline:
		return "离线"
	case DeviceStatusOnline:
		return "在线"
	default:
		return "unknown"
	}
}
func GetDeviceStatusValue(str string) int32 {
	switch str {
	case "未激活":
		return DeviceStatusNotActivated
	case "离线":
		return DeviceStatusOffline
	case "在线":
		return DeviceStatusOnline
	default:
		return DeviceStatusUnKnown
	}
}

// 设备检查相关常量
const (
	DeviceOfflineTime   int32 = 180 //设备离线时间，单位秒
	DeviceCheckInterval int32 = 10  //设备周期检查时间间隔
)

// MessageHandler Kafka消息处理函数定义
type MessageHandler func(topic string, key []byte, payload []byte)

// kafka主题定义
const (
	Kafka_topic_notify_up   string = "up.egw.notify"   // Edge上报其他数据到云端
	Kafka_topic_notify_down string = "down.egw.notify" // 云端下发命令到Edge
)

// 日志消息级别定义
const (
	Msg_level_info  int32 = 1
	Msg_level_warn  int32 = 2
	Msg_level_error int32 = 3
)

// 云<-->边 消息类型定义
const (
	Edge_command_down string = "edge_command_down" // 下发设备命令
	Edge_command_up   string = "edge_command_up"   // 下发设备命令
	Edge_ota_down     string = "edge_ota_down"
	Edge_ota_up       string = "edge_ota_up"
	Edge_heartbeat_up string = "edge_heartbeat_up"
)

// 设备命令执行状态 1-执行中, 2-完成, 3-失败, 4-超时
const (
	Command_status_executing int32 = 1 // 执行中
	Command_status_completed int32 = 2 // 完成
	Command_status_failed    int32 = 3 // 失败
	Command_status_timeout   int32 = 4 // 超时
)

// OTA命令类型
const (
	Ota_command_type_upgrade string = "Ota_command_type_upgrade"
	Ota_command_type_cancel  string = "Ota_command_type_cancel"
)

// 升级包状态
const (
	Ota_package_status_draft   int32 = 1 //草稿
	Ota_package_status_release int32 = 2 //发布
)

// 升级批次状态
const (
	Ota_batch_status_executing int32 = 1 //执行中
	Ota_batch_status_completed int32 = 2 //完成
	Ota_batch_status_failed    int32 = 3 //失败
	Ota_batch_status_pending   int32 = 4 //待处理
)

// 升级记录状态 1-待确认，2-待推送，3-已推送，4-升级中，5-成功，6-失败，7-已取消
const (
	Ota_record_status_pending      int32 = 1 //待确认
	Ota_record_status_to_be_pushed int32 = 2 //待推送
	Ota_record_status_pushed       int32 = 3 //已推送
	Ota_record_status_upgrading    int32 = 4 //升级中
	Ota_record_status_success      int32 = 5 //成功
	Ota_record_status_failed       int32 = 6 //失败
	Ota_record_status_canceled     int32 = 7 //已取消
)

// 升级范围类型：1-手动勾选设备,2-全选
const (
	Ota_scope_type_manual int32 = 1 //手动勾选设备
	Ota_scope_type_all    int32 = 2 //全选
)

// 升级策略：1-灰度升级，2-批量升级
const (
	Ota_strategy_type_gray int32 = 1 //灰度升级
	Ota_strategy_type_full int32 = 2 //全量升级
)

// 文件上传相关常量
const (
	File_bucket          string = "firmware"
	File_upload_expiry   int64  = 30 * 60 //seconds
	File_download_expiry int64  = 30 * 60 //seconds
)
