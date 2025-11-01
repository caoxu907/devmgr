// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UpgradeRecord is the golang structure for table upgrade_record.
type UpgradeRecord struct {
	Id            int         `json:"id"            orm:"id"             description:""` //
	UpgradeStatus int         `json:"upgradeStatus" orm:"upgrade_status" description:""` //
	BatchId       int         `json:"batchId"       orm:"batch_id"       description:""` //
	DeviceId      int         `json:"deviceId"      orm:"device_id"      description:""` //
	IsActive      bool        `json:"isActive"      orm:"is_active"      description:""` //
	RetryTimes    int         `json:"retryTimes"    orm:"retry_times"    description:""` //
	Progress      int         `json:"progress"      orm:"progress"       description:""` //
	IsGrayBatch   bool        `json:"isGrayBatch"   orm:"is_gray_batch"  description:""` //
	FromVersion   string      `json:"fromVersion"   orm:"from_version"   description:""` //
	ToVersion     string      `json:"toVersion"     orm:"to_version"     description:""` //
	Duration      int         `json:"duration"      orm:"duration"       description:""` //
	IsRollback    bool        `json:"isRollback"    orm:"is_rollback"    description:""` //
	CreateTime    *gtime.Time `json:"createTime"    orm:"create_time"    description:""` //
	PlanTime      *gtime.Time `json:"planTime"      orm:"plan_time"      description:""` //
	UpdateTime    *gtime.Time `json:"updateTime"    orm:"update_time"    description:""` //
	StartTime     *gtime.Time `json:"startTime"     orm:"start_time"     description:""` //
	CompleteTime  *gtime.Time `json:"completeTime"  orm:"complete_time"  description:""` //
}
