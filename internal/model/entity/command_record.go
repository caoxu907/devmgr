// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CommandRecord is the golang structure for table command_record.
type CommandRecord struct {
	Id              int         `json:"id"              orm:"id"                description:""` //
	CommandConfigId int         `json:"commandConfigId" orm:"command_config_id" description:""` //
	Params          string      `json:"params"          orm:"params"            description:""` //
	ClientId        string      `json:"clientId"        orm:"client_id"         description:""` //
	Status          int         `json:"status"          orm:"status"            description:""` //
	ResultCode      int         `json:"resultCode"      orm:"result_code"       description:""` //
	ResultData      string      `json:"resultData"      orm:"result_data"       description:""` //
	CreateTime      *gtime.Time `json:"createTime"      orm:"create_time"       description:""` //
	ExecuteTime     *gtime.Time `json:"executeTime"     orm:"execute_time"      description:""` //
	CompleteTime    *gtime.Time `json:"completeTime"    orm:"complete_time"     description:""` //
	TimeoutSeconds  int         `json:"timeoutSeconds"  orm:"timeout_seconds"   description:""` //
}
