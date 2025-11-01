// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TaskStatusDao is the data access object for the table task_status.
type TaskStatusDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TaskStatusColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TaskStatusColumns defines and stores column names for the table task_status.
type TaskStatusColumns struct {
	Id             string //
	TaskId         string //
	TaskName       string //
	TaskType       string //
	TaskSource     string //
	ProductId      string //
	DeviceId       string //
	Status         string //
	Progress       string //
	RetryTimes     string //
	MaxRetry       string //
	Message        string //
	Metadata       string //
	LastUpdateFrom string //
	CreateTime     string //
	StartTime      string //
	CompleteTime   string //
	UpdateTime     string //
	Deleted        string //
}

// taskStatusColumns holds the columns for the table task_status.
var taskStatusColumns = TaskStatusColumns{
	Id:             "id",
	TaskId:         "task_id",
	TaskName:       "task_name",
	TaskType:       "task_type",
	TaskSource:     "task_source",
	ProductId:      "product_id",
	DeviceId:       "device_id",
	Status:         "status",
	Progress:       "progress",
	RetryTimes:     "retry_times",
	MaxRetry:       "max_retry",
	Message:        "message",
	Metadata:       "metadata",
	LastUpdateFrom: "last_update_from",
	CreateTime:     "create_time",
	StartTime:      "start_time",
	CompleteTime:   "complete_time",
	UpdateTime:     "update_time",
	Deleted:        "deleted",
}

// NewTaskStatusDao creates and returns a new DAO object for table data access.
func NewTaskStatusDao(handlers ...gdb.ModelHandler) *TaskStatusDao {
	return &TaskStatusDao{
		group:    "default",
		table:    "task_status",
		columns:  taskStatusColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TaskStatusDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TaskStatusDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TaskStatusDao) Columns() TaskStatusColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TaskStatusDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TaskStatusDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TaskStatusDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
