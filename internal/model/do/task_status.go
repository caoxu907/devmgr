// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TaskStatus is the golang structure of table task_status for DAO operations like Where/Data.
type TaskStatus struct {
	g.Meta         `orm:"table:task_status, do:true"`
	Id             any         //
	TaskId         any         //
	TaskName       any         //
	TaskType       any         //
	TaskSource     any         //
	ProductId      any         //
	DeviceId       any         //
	Status         any         //
	Progress       any         //
	RetryTimes     any         //
	MaxRetry       any         //
	Message        any         //
	Metadata       any         //
	LastUpdateFrom any         //
	CreateTime     *gtime.Time //
	StartTime      *gtime.Time //
	CompleteTime   *gtime.Time //
	UpdateTime     *gtime.Time //
	Deleted        any         //
}
