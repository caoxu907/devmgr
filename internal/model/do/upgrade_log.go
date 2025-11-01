// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UpgradeLog is the golang structure of table upgrade_log for DAO operations like Where/Data.
type UpgradeLog struct {
	g.Meta    `orm:"table:upgrade_log, do:true"`
	Id        any         //
	RecordId  any         //
	Timestamp *gtime.Time //
	Level     any         //
	Content   any         //
}
