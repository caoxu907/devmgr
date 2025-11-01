// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UpgradePackageDao is the data access object for the table upgrade_package.
type UpgradePackageDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  UpgradePackageColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// UpgradePackageColumns defines and stores column names for the table upgrade_package.
type UpgradePackageColumns struct {
	Id                 string //
	PackageStatus      string //
	PackageName        string //
	ProductId          string //
	Version            string //
	SupportAllVersions string //
	AvailableVersions  string //
	PackagePath        string //
	PackageSize        string //
	PackageHash        string //
	PackageDesc        string //
	CreateTime         string //
	UpdateTime         string //
	Deleted            string //
}

// upgradePackageColumns holds the columns for the table upgrade_package.
var upgradePackageColumns = UpgradePackageColumns{
	Id:                 "id",
	PackageStatus:      "package_status",
	PackageName:        "package_name",
	ProductId:          "product_id",
	Version:            "version",
	SupportAllVersions: "support_all_versions",
	AvailableVersions:  "available_versions",
	PackagePath:        "package_path",
	PackageSize:        "package_size",
	PackageHash:        "package_hash",
	PackageDesc:        "package_desc",
	CreateTime:         "create_time",
	UpdateTime:         "update_time",
	Deleted:            "deleted",
}

// NewUpgradePackageDao creates and returns a new DAO object for table data access.
func NewUpgradePackageDao(handlers ...gdb.ModelHandler) *UpgradePackageDao {
	return &UpgradePackageDao{
		group:    "default",
		table:    "upgrade_package",
		columns:  upgradePackageColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UpgradePackageDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UpgradePackageDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UpgradePackageDao) Columns() UpgradePackageColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UpgradePackageDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UpgradePackageDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *UpgradePackageDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
