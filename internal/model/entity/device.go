// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Device is the golang structure for table device.
type Device struct {
	Id              int         `json:"id"              orm:"id"               description:""` //
	Status          int         `json:"status"          orm:"status"           description:""` //
	UpdateStatus    int         `json:"updateStatus"    orm:"update_status"    description:""` //
	ProductId       int         `json:"productId"       orm:"product_id"       description:""` //
	TenantId        int         `json:"tenantId"        orm:"tenant_id"        description:""` //
	DeviceKey       string      `json:"deviceKey"       orm:"device_key"       description:""` //
	DeviceDesc      string      `json:"deviceDesc"      orm:"device_desc"      description:""` //
	DeviceSecret    string      `json:"deviceSecret"    orm:"device_secret"    description:""` //
	MachineCode     string      `json:"machineCode"     orm:"machine_code"     description:""` //
	HardwareVersion string      `json:"hardwareVersion" orm:"hardware_version" description:""` //
	SoftwareVersion string      `json:"softwareVersion" orm:"software_version" description:""` //
	Cpu             float64     `json:"cpu"             orm:"cpu"              description:""` //
	Mem             float64     `json:"mem"             orm:"mem"              description:""` //
	MemTotal        float64     `json:"memTotal"        orm:"mem_total"        description:""` //
	Disk            float64     `json:"disk"            orm:"disk"             description:""` //
	DiskTotal       float64     `json:"diskTotal"       orm:"disk_total"       description:""` //
	Address         string      `json:"address"         orm:"address"          description:""` //
	CreateTime      *gtime.Time `json:"createTime"      orm:"create_time"      description:""` //
	ActivateTime    *gtime.Time `json:"activateTime"    orm:"activate_time"    description:""` //
	LastOnlineTime  *gtime.Time `json:"lastOnlineTime"  orm:"last_online_time" description:""` //
	Deleted         bool        `json:"deleted"         orm:"deleted"          description:""` //
}
