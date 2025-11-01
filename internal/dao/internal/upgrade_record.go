// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UpgradeRecordDao is the data access object for the table upgrade_record.
type UpgradeRecordDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  UpgradeRecordColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// UpgradeRecordColumns defines and stores column names for the table upgrade_record.
type UpgradeRecordColumns struct {
	Id            string //
	UpgradeStatus string //
	BatchId       string //
	DeviceId      string //
	IsActive      string //
	RetryTimes    string //
	Progress      string //
	IsGrayBatch   string //
	FromVersion   string //
	ToVersion     string //
	Duration      string //
	IsRollback    string //
	CreateTime    string //
	PlanTime      string //
	UpdateTime    string //
	StartTime     string //
	CompleteTime  string //
}

// upgradeRecordColumns holds the columns for the table upgrade_record.
var upgradeRecordColumns = UpgradeRecordColumns{
	Id:            "id",
	UpgradeStatus: "upgrade_status",
	BatchId:       "batch_id",
	DeviceId:      "device_id",
	IsActive:      "is_active",
	RetryTimes:    "retry_times",
	Progress:      "progress",
	IsGrayBatch:   "is_gray_batch",
	FromVersion:   "from_version",
	ToVersion:     "to_version",
	Duration:      "duration",
	IsRollback:    "is_rollback",
	CreateTime:    "create_time",
	PlanTime:      "plan_time",
	UpdateTime:    "update_time",
	StartTime:     "start_time",
	CompleteTime:  "complete_time",
}

// NewUpgradeRecordDao creates and returns a new DAO object for table data access.
func NewUpgradeRecordDao(handlers ...gdb.ModelHandler) *UpgradeRecordDao {
	return &UpgradeRecordDao{
		group:    "default",
		table:    "upgrade_record",
		columns:  upgradeRecordColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UpgradeRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UpgradeRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UpgradeRecordDao) Columns() UpgradeRecordColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UpgradeRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UpgradeRecordDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UpgradeRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
