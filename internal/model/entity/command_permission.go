// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CommandPermission is the golang structure for table command_permission.
type CommandPermission struct {
	Id              int         `json:"id"              orm:"id"                description:""` //
	RoleId          int         `json:"roleId"          orm:"role_id"           description:""` //
	CommandConfigId int         `json:"commandConfigId" orm:"command_config_id" description:""` //
	CreateTime      *gtime.Time `json:"createTime"      orm:"create_time"       description:""` //
}
