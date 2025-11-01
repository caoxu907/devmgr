// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TaskStatus is the golang structure for table task_status.
type TaskStatus struct {
	Id             int         `json:"id"             orm:"id"               description:""` //
	TaskId         string      `json:"taskId"         orm:"task_id"          description:""` //
	TaskName       string      `json:"taskName"       orm:"task_name"        description:""` //
	TaskType       string      `json:"taskType"       orm:"task_type"        description:""` //
	TaskSource     string      `json:"taskSource"     orm:"task_source"      description:""` //
	ProductId      int         `json:"productId"      orm:"product_id"       description:""` //
	DeviceId       int         `json:"deviceId"       orm:"device_id"        description:""` //
	Status         int         `json:"status"         orm:"status"           description:""` //
	Progress       int         `json:"progress"       orm:"progress"         description:""` //
	RetryTimes     int         `json:"retryTimes"     orm:"retry_times"      description:""` //
	MaxRetry       int         `json:"maxRetry"       orm:"max_retry"        description:""` //
	Message        string      `json:"message"        orm:"message"          description:""` //
	Metadata       string      `json:"metadata"       orm:"metadata"         description:""` //
	LastUpdateFrom string      `json:"lastUpdateFrom" orm:"last_update_from" description:""` //
	CreateTime     *gtime.Time `json:"createTime"     orm:"create_time"      description:""` //
	StartTime      *gtime.Time `json:"startTime"      orm:"start_time"       description:""` //
	CompleteTime   *gtime.Time `json:"completeTime"   orm:"complete_time"    description:""` //
	UpdateTime     *gtime.Time `json:"updateTime"     orm:"update_time"      description:""` //
	Deleted        bool        `json:"deleted"        orm:"deleted"          description:""` //
}
