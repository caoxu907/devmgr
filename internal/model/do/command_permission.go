// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CommandPermission is the golang structure of table command_permission for DAO operations like Where/Data.
type CommandPermission struct {
	g.Meta          `orm:"table:command_permission, do:true"`
	Id              interface{} //
	RoleId          interface{} //
	CommandConfigId interface{} //
	CreateTime      *gtime.Time //
}
