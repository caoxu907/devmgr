// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CommandConfigDao is the data access object for the table command_config.
type CommandConfigDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  CommandConfigColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// CommandConfigColumns defines and stores column names for the table command_config.
type CommandConfigColumns struct {
	Id             string //
	DeviceId       string //
	CommandName    string //
	CommandDesc    string //
	CommandType    string //
	TimeoutSeconds string //
	CreateTime     string //
	UpdateTime     string //
	Deleted        string //
}

// commandConfigColumns holds the columns for the table command_config.
var commandConfigColumns = CommandConfigColumns{
	Id:             "id",
	DeviceId:       "device_id",
	CommandName:    "command_name",
	CommandDesc:    "command_desc",
	CommandType:    "command_type",
	TimeoutSeconds: "timeout_seconds",
	CreateTime:     "create_time",
	UpdateTime:     "update_time",
	Deleted:        "deleted",
}

// NewCommandConfigDao creates and returns a new DAO object for table data access.
func NewCommandConfigDao(handlers ...gdb.ModelHandler) *CommandConfigDao {
	return &CommandConfigDao{
		group:    "default",
		table:    "command_config",
		columns:  commandConfigColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CommandConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CommandConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CommandConfigDao) Columns() CommandConfigColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CommandConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CommandConfigDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CommandConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
