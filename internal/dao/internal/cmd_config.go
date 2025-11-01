// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CmdConfigDao is the data access object for the table cmd_config.
type CmdConfigDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  CmdConfigColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// CmdConfigColumns defines and stores column names for the table cmd_config.
type CmdConfigColumns struct {
	Id             string //
	Name           string //
	CommandKey     string //
	DeviceId       string //
	Description    string //
	CommandType    string //
	CommandLine    string //
	TimeoutSeconds string //
	CreateTime     string //
	UpdateTime     string //
	Deleted        string //
}

// cmdConfigColumns holds the columns for the table cmd_config.
var cmdConfigColumns = CmdConfigColumns{
	Id:             "id",
	Name:           "name",
	CommandKey:     "command_key",
	DeviceId:       "device_id",
	Description:    "description",
	CommandType:    "command_type",
	CommandLine:    "command_line",
	TimeoutSeconds: "timeout_seconds",
	CreateTime:     "create_time",
	UpdateTime:     "update_time",
	Deleted:        "deleted",
}

// NewCmdConfigDao creates and returns a new DAO object for table data access.
func NewCmdConfigDao(handlers ...gdb.ModelHandler) *CmdConfigDao {
	return &CmdConfigDao{
		group:    "default",
		table:    "cmd_config",
		columns:  cmdConfigColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CmdConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CmdConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CmdConfigDao) Columns() CmdConfigColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CmdConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CmdConfigDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CmdConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
