// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UpgradeRecord is the golang structure of table upgrade_record for DAO operations like Where/Data.
type UpgradeRecord struct {
	g.Meta        `orm:"table:upgrade_record, do:true"`
	Id            any         //
	UpgradeStatus any         //
	BatchId       any         //
	DeviceId      any         //
	IsActive      any         //
	RetryTimes    any         //
	Progress      any         //
	IsGrayBatch   any         //
	FromVersion   any         //
	ToVersion     any         //
	Duration      any         //
	IsRollback    any         //
	CreateTime    *gtime.Time //
	PlanTime      *gtime.Time //
	UpdateTime    *gtime.Time //
	StartTime     *gtime.Time //
	CompleteTime  *gtime.Time //
}
