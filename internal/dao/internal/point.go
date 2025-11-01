// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PointDao is the data access object for the table point.
type PointDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  PointColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// PointColumns defines and stores column names for the table point.
type PointColumns struct {
	Id          string //
	SubdeviceId string //
	PointKey    string //
	PointName   string //
	PointDesc   string //
	DataType    string //
	PointType   string //
	Deleted     string //
}

// pointColumns holds the columns for the table point.
var pointColumns = PointColumns{
	Id:          "id",
	SubdeviceId: "subdevice_id",
	PointKey:    "point_key",
	PointName:   "point_name",
	PointDesc:   "point_desc",
	DataType:    "data_type",
	PointType:   "point_type",
	Deleted:     "deleted",
}

// NewPointDao creates and returns a new DAO object for table data access.
func NewPointDao(handlers ...gdb.ModelHandler) *PointDao {
	return &PointDao{
		group:    "default",
		table:    "point",
		columns:  pointColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PointDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PointDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PointDao) Columns() PointColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PointDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PointDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PointDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
