package model

type DeviceMessage struct {
	MessageType string      `json:"type"`
	Content     interface{} `json:"content"`
}

// DeviceMessageCommandSend 下发命令
type DeviceMessageCommandSend struct {
	CommandId int32  `json:"exec_id"`
	Command   string `json:"command"`
	Param     string `json:"param"`
}

// DeviceMessageCommandResponse 命令执行结果
type DeviceMessageCommandResponse struct {
	CommandId int32  `json:"exec_id"`
	Level     int32  `json:"level"`
	Content   string `json:"content"`
	Status    int32  `json:"status"`
}

// DeviceMessageOtaSend 下发OTA升级
type DeviceMessageOtaSend struct {
	MessageType    string `json:"type"` //Ota_Command_Type_Cancel
	Id             int32  `json:"id"`
	Url            string `json:"url"`
	PackageVersion string `json:"package_version"`
	PackageHash    string `json:"package_hash"`
	PackageSize    int32  `json:"package_size"`
}

type DeviceMessageOtaResponse struct {
	MessageType string `json:"type"`
	Id          int32  `json:"id"`
	Level       int32  `json:"level"`
	Content     string `json:"content"`
	Status      int32  `json:"status"`
}

// DeviceMessageHeartbeat 设备上报心跳
type DeviceMessageHeartbeat struct {
	Version   string  `json:"version"`    // 版本
	Cpu       float32 `json:"cpu"`        // CPU占用率
	Mem       float32 `json:"mem"`        // 内存占用大小（MB）
	MemTotal  float32 `json:"mem_total"`  // 内存总大小（MB）
	Disk      float32 `json:"disk"`       // 磁盘占用大小（MB）
	DiskTotal float32 `json:"disk_total"` // 磁盘总大小（MB）
}
