// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UpgradeLog is the golang structure for table upgrade_log.
type UpgradeLog struct {
	Id        int         `json:"id"        orm:"id"        description:""` //
	RecordId  int         `json:"recordId"  orm:"record_id" description:""` //
	Timestamp *gtime.Time `json:"timestamp" orm:"timestamp" description:""` //
	Level     int         `json:"level"     orm:"level"     description:""` //
	Content   string      `json:"content"   orm:"content"   description:""` //
}
