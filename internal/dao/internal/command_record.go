// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CommandRecordDao is the data access object for the table command_record.
type CommandRecordDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  CommandRecordColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// CommandRecordColumns defines and stores column names for the table command_record.
type CommandRecordColumns struct {
	Id              string //
	CommandConfigId string //
	Params          string //
	ClientId        string //
	Status          string //
	ResultCode      string //
	ResultData      string //
	CreateTime      string //
	ExecuteTime     string //
	CompleteTime    string //
	TimeoutSeconds  string //
}

// commandRecordColumns holds the columns for the table command_record.
var commandRecordColumns = CommandRecordColumns{
	Id:              "id",
	CommandConfigId: "command_config_id",
	Params:          "params",
	ClientId:        "client_id",
	Status:          "status",
	ResultCode:      "result_code",
	ResultData:      "result_data",
	CreateTime:      "create_time",
	ExecuteTime:     "execute_time",
	CompleteTime:    "complete_time",
	TimeoutSeconds:  "timeout_seconds",
}

// NewCommandRecordDao creates and returns a new DAO object for table data access.
func NewCommandRecordDao(handlers ...gdb.ModelHandler) *CommandRecordDao {
	return &CommandRecordDao{
		group:    "default",
		table:    "command_record",
		columns:  commandRecordColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CommandRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CommandRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CommandRecordDao) Columns() CommandRecordColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CommandRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CommandRecordDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CommandRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
