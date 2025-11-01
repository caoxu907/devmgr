// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CommandRecord is the golang structure of table command_record for DAO operations like Where/Data.
type CommandRecord struct {
	g.Meta          `orm:"table:command_record, do:true"`
	Id              interface{} //
	CommandConfigId interface{} //
	Params          interface{} //
	ClientId        interface{} //
	Status          interface{} //
	ResultCode      interface{} //
	ResultData      interface{} //
	CreateTime      *gtime.Time //
	ExecuteTime     *gtime.Time //
	CompleteTime    *gtime.Time //
	TimeoutSeconds  interface{} //
}
