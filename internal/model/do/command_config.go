// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CommandConfig is the golang structure of table command_config for DAO operations like Where/Data.
type CommandConfig struct {
	g.Meta         `orm:"table:command_config, do:true"`
	Id             interface{} //
	DeviceId       interface{} //
	CommandName    interface{} //
	CommandDesc    interface{} //
	CommandType    interface{} //
	TimeoutSeconds interface{} //
	CreateTime     *gtime.Time //
	UpdateTime     *gtime.Time //
	Deleted        interface{} //
}
