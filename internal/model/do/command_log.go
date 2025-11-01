// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CommandLog is the golang structure of table command_log for DAO operations like Where/Data.
type CommandLog struct {
	g.Meta    `orm:"table:command_log, do:true"`
	Id        interface{} //
	CommandId interface{} //
	Timestamp *gtime.Time //
	Level     interface{} //
	Content   interface{} //
}
