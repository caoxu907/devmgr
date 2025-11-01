// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Product is the golang structure for table product.
type Product struct {
	Id           int         `json:"id"           orm:"id"            description:""` //
	ProductName  string      `json:"productName"  orm:"product_name"  description:""` //
	ProductAlias string      `json:"productAlias" orm:"product_alias" description:""` //
	ProductDesc  string      `json:"productDesc"  orm:"product_desc"  description:""` //
	PackageId    string      `json:"packageId"    orm:"package_id"    description:""` //
	Deleted      bool        `json:"deleted"      orm:"deleted"       description:""` //
	CreateTime   *gtime.Time `json:"createTime"   orm:"create_time"   description:""` //
	UpdateTime   *gtime.Time `json:"updateTime"   orm:"update_time"   description:""` //
}
