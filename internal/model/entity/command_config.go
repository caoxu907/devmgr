// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CommandConfig is the golang structure for table command_config.
type CommandConfig struct {
	Id             int         `json:"id"             orm:"id"              description:""` //
	DeviceId       int         `json:"deviceId"       orm:"device_id"       description:""` //
	CommandName    string      `json:"commandName"    orm:"command_name"    description:""` //
	CommandDesc    string      `json:"commandDesc"    orm:"command_desc"    description:""` //
	CommandType    int         `json:"commandType"    orm:"command_type"    description:""` //
	TimeoutSeconds int         `json:"timeoutSeconds" orm:"timeout_seconds" description:""` //
	CreateTime     *gtime.Time `json:"createTime"     orm:"create_time"     description:""` //
	UpdateTime     *gtime.Time `json:"updateTime"     orm:"update_time"     description:""` //
	Deleted        bool        `json:"deleted"        orm:"deleted"         description:""` //
}
