// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UpgradePackage is the golang structure of table upgrade_package for DAO operations like Where/Data.
type UpgradePackage struct {
	g.Meta             `orm:"table:upgrade_package, do:true"`
	Id                 any         //
	PackageStatus      any         //
	PackageName        any         //
	ProductId          any         //
	Version            any         //
	SupportAllVersions any         //
	AvailableVersions  any         //
	PackagePath        any         //
	PackageSize        any         //
	PackageHash        any         //
	PackageDesc        any         //
	CreateTime         *gtime.Time //
	UpdateTime         *gtime.Time //
	Deleted            any         //
}
