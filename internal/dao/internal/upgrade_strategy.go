// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UpgradeStrategyDao is the data access object for the table upgrade_strategy.
type UpgradeStrategyDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  UpgradeStrategyColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// UpgradeStrategyColumns defines and stores column names for the table upgrade_strategy.
type UpgradeStrategyColumns struct {
	Id                   string //
	Name                 string //
	PackageId            string //
	ScopeType            string //
	ScopeDevices         string //
	UserConfirm          string //
	RetryInterval        string //
	MaxRetryCount        string //
	OverrideExistingTask string //
	UpgradeMode          string //
	StrategyType         string //
	GrayPercentage       string //
	GraySuccessThreshold string //
	Concurrency          string //
	Status               string //
	TotalDevices         string //
	SuccessCount         string //
	FailureCount         string //
	PendingCount         string //
	CreateTime           string //
	UpdateTime           string //
	StartTime            string //
	CompleteTime         string //
	Deleted              string //
}

// upgradeStrategyColumns holds the columns for the table upgrade_strategy.
var upgradeStrategyColumns = UpgradeStrategyColumns{
	Id:                   "id",
	Name:                 "name",
	PackageId:            "package_id",
	ScopeType:            "scope_type",
	ScopeDevices:         "scope_devices",
	UserConfirm:          "user_confirm",
	RetryInterval:        "retry_interval",
	MaxRetryCount:        "max_retry_count",
	OverrideExistingTask: "override_existing_task",
	UpgradeMode:          "upgrade_mode",
	StrategyType:         "strategy_type",
	GrayPercentage:       "gray_percentage",
	GraySuccessThreshold: "gray_success_threshold",
	Concurrency:          "concurrency",
	Status:               "status",
	TotalDevices:         "total_devices",
	SuccessCount:         "success_count",
	FailureCount:         "failure_count",
	PendingCount:         "pending_count",
	CreateTime:           "create_time",
	UpdateTime:           "update_time",
	StartTime:            "start_time",
	CompleteTime:         "complete_time",
	Deleted:              "deleted",
}

// NewUpgradeStrategyDao creates and returns a new DAO object for table data access.
func NewUpgradeStrategyDao(handlers ...gdb.ModelHandler) *UpgradeStrategyDao {
	return &UpgradeStrategyDao{
		group:    "default",
		table:    "upgrade_strategy",
		columns:  upgradeStrategyColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UpgradeStrategyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UpgradeStrategyDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UpgradeStrategyDao) Columns() UpgradeStrategyColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UpgradeStrategyDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UpgradeStrategyDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UpgradeStrategyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
