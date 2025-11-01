// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Product is the golang structure of table product for DAO operations like Where/Data.
type Product struct {
	g.Meta       `orm:"table:product, do:true"`
	Id           any         //
	ProductName  any         //
	ProductAlias any         //
	ProductDesc  any         //
	PackageId    any         //
	Deleted      any         //
	CreateTime   *gtime.Time //
	UpdateTime   *gtime.Time //
}
