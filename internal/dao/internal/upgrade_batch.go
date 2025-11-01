// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UpgradeBatchDao is the data access object for the table upgrade_batch.
type UpgradeBatchDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  UpgradeBatchColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// UpgradeBatchColumns defines and stores column names for the table upgrade_batch.
type UpgradeBatchColumns struct {
	Id                   string //
	Batchstatus          string //
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
	GrayCompleted        string //
	MaxConcurrent        string //
	GroupInterval        string //
	RandomDelay          string //
	CreateTime           string //
	UpdateTime           string //
	StartTime            string //
	CompleteTime         string //
	Deleted              string //
	TotalDevices         string //
	SuccessDevices       string //
	FailureDevices       string //
	PendingDevices       string //
}

// upgradeBatchColumns holds the columns for the table upgrade_batch.
var upgradeBatchColumns = UpgradeBatchColumns{
	Id:                   "id",
	Batchstatus:          "batchstatus",
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
	GrayCompleted:        "gray_completed",
	MaxConcurrent:        "max_concurrent",
	GroupInterval:        "group_interval",
	RandomDelay:          "random_delay",
	CreateTime:           "create_time",
	UpdateTime:           "update_time",
	StartTime:            "start_time",
	CompleteTime:         "complete_time",
	Deleted:              "deleted",
	TotalDevices:         "total_devices",
	SuccessDevices:       "success_devices",
	FailureDevices:       "failure_devices",
	PendingDevices:       "pending_devices",
}

// NewUpgradeBatchDao creates and returns a new DAO object for table data access.
func NewUpgradeBatchDao(handlers ...gdb.ModelHandler) *UpgradeBatchDao {
	return &UpgradeBatchDao{
		group:    "default",
		table:    "upgrade_batch",
		columns:  upgradeBatchColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UpgradeBatchDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UpgradeBatchDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UpgradeBatchDao) Columns() UpgradeBatchColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UpgradeBatchDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UpgradeBatchDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UpgradeBatchDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
