// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UpgradeMessage is the golang structure of table upgrade_message for DAO operations like Where/Data.
type UpgradeMessage struct {
	g.Meta    `orm:"table:upgrade_message, do:true"`
	Id        interface{} //
	RecordId  interface{} //
	Timestamp *gtime.Time //
	Level     interface{} //
	Content   interface{} //
}
