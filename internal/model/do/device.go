// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Device is the golang structure of table device for DAO operations like Where/Data.
type Device struct {
	g.Meta          `orm:"table:device, do:true"`
	Id              any         //
	Status          any         //
	UpdateStatus    any         //
	ProductId       any         //
	TenantId        any         //
	DeviceKey       any         //
	DeviceDesc      any         //
	DeviceSecret    any         //
	MachineCode     any         //
	HardwareVersion any         //
	SoftwareVersion any         //
	Cpu             any         //
	Mem             any         //
	MemTotal        any         //
	Disk            any         //
	DiskTotal       any         //
	Address         any         //
	CreateTime      *gtime.Time //
	ActivateTime    *gtime.Time //
	LastOnlineTime  *gtime.Time //
	Deleted         any         //
}
