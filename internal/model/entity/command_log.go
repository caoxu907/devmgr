// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CommandLog is the golang structure for table command_log.
type CommandLog struct {
	Id        int         `json:"id"        orm:"id"         description:""` //
	CommandId int         `json:"commandId" orm:"command_id" description:""` //
	Timestamp *gtime.Time `json:"timestamp" orm:"timestamp"  description:""` //
	Level     int         `json:"level"     orm:"level"      description:""` //
	Content   string      `json:"content"   orm:"content"    description:""` //
}
