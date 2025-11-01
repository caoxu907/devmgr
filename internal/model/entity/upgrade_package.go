// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UpgradePackage is the golang structure for table upgrade_package.
type UpgradePackage struct {
	Id                 int         `json:"id"                 orm:"id"                   description:""` //
	PackageStatus      int         `json:"packageStatus"      orm:"package_status"       description:""` //
	PackageName        string      `json:"packageName"        orm:"package_name"         description:""` //
	ProductId          int         `json:"productId"          orm:"product_id"           description:""` //
	Version            string      `json:"version"            orm:"version"              description:""` //
	SupportAllVersions bool        `json:"supportAllVersions" orm:"support_all_versions" description:""` //
	AvailableVersions  string      `json:"availableVersions"  orm:"available_versions"   description:""` //
	PackagePath        string      `json:"packagePath"        orm:"package_path"         description:""` //
	PackageSize        int         `json:"packageSize"        orm:"package_size"         description:""` //
	PackageHash        string      `json:"packageHash"        orm:"package_hash"         description:""` //
	PackageDesc        string      `json:"packageDesc"        orm:"package_desc"         description:""` //
	CreateTime         *gtime.Time `json:"createTime"         orm:"create_time"          description:""` //
	UpdateTime         *gtime.Time `json:"updateTime"         orm:"update_time"          description:""` //
	Deleted            bool        `json:"deleted"            orm:"deleted"              description:""` //
}
