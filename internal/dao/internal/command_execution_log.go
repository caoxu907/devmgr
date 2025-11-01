// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CommandExecutionLogDao is the data access object for the table command_execution_log.
type CommandExecutionLogDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  CommandExecutionLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// CommandExecutionLogColumns defines and stores column names for the table command_execution_log.
type CommandExecutionLogColumns struct {
	Id        string //
	CommandId string //
	Status    string //
	Message   string //
	LogTime   string //
}

// commandExecutionLogColumns holds the columns for the table command_execution_log.
var commandExecutionLogColumns = CommandExecutionLogColumns{
	Id:        "id",
	CommandId: "command_id",
	Status:    "status",
	Message:   "message",
	LogTime:   "log_time",
}

// NewCommandExecutionLogDao creates and returns a new DAO object for table data access.
func NewCommandExecutionLogDao(handlers ...gdb.ModelHandler) *CommandExecutionLogDao {
	return &CommandExecutionLogDao{
		group:    "default",
		table:    "command_execution_log",
		columns:  commandExecutionLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CommandExecutionLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CommandExecutionLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CommandExecutionLogDao) Columns() CommandExecutionLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CommandExecutionLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CommandExecutionLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CommandExecutionLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
