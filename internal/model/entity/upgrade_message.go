// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UpgradeMessage is the golang structure for table upgrade_message.
type UpgradeMessage struct {
	Id        int64       `json:"id"        orm:"id"        description:""` //
	RecordId  int64       `json:"recordId"  orm:"record_id" description:""` //
	Timestamp *gtime.Time `json:"timestamp" orm:"timestamp" description:""` //
	Level     string      `json:"level"     orm:"level"     description:""` //
	Content   string      `json:"content"   orm:"content"   description:""` //
}
