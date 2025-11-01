// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CmdRecordDao is the data access object for the table cmd_record.
type CmdRecordDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  CmdRecordColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// CmdRecordColumns defines and stores column names for the table cmd_record.
type CmdRecordColumns struct {
	Id              string //
	CommandConfigId string //
	DeviceId        string //
	Params          string //
	ClientId        string //
	Status          string //
	ResultCode      string //
	ResultData      string //
	ErrorMessage    string //
	CreateTime      string //
	ExecuteTime     string //
	CompleteTime    string //
	TimeoutSeconds  string //
}

// cmdRecordColumns holds the columns for the table cmd_record.
var cmdRecordColumns = CmdRecordColumns{
	Id:              "id",
	CommandConfigId: "command_config_id",
	DeviceId:        "device_id",
	Params:          "params",
	ClientId:        "client_id",
	Status:          "status",
	ResultCode:      "result_code",
	ResultData:      "result_data",
	ErrorMessage:    "error_message",
	CreateTime:      "create_time",
	ExecuteTime:     "execute_time",
	CompleteTime:    "complete_time",
	TimeoutSeconds:  "timeout_seconds",
}

// NewCmdRecordDao creates and returns a new DAO object for table data access.
func NewCmdRecordDao(handlers ...gdb.ModelHandler) *CmdRecordDao {
	return &CmdRecordDao{
		group:    "default",
		table:    "cmd_record",
		columns:  cmdRecordColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CmdRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CmdRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CmdRecordDao) Columns() CmdRecordColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CmdRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CmdRecordDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CmdRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
