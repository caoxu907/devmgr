// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DeviceDao is the data access object for the table device.
type DeviceDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DeviceColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DeviceColumns defines and stores column names for the table device.
type DeviceColumns struct {
	Id              string //
	Status          string //
	UpdateStatus    string //
	ProductId       string //
	TenantId        string //
	DeviceKey       string //
	DeviceDesc      string //
	DeviceSecret    string //
	MachineCode     string //
	HardwareVersion string //
	SoftwareVersion string //
	Cpu             string //
	Mem             string //
	MemTotal        string //
	Disk            string //
	DiskTotal       string //
	Address         string //
	CreateTime      string //
	ActivateTime    string //
	LastOnlineTime  string //
	Deleted         string //
}

// deviceColumns holds the columns for the table device.
var deviceColumns = DeviceColumns{
	Id:              "id",
	Status:          "status",
	UpdateStatus:    "update_status",
	ProductId:       "product_id",
	TenantId:        "tenant_id",
	DeviceKey:       "device_key",
	DeviceDesc:      "device_desc",
	DeviceSecret:    "device_secret",
	MachineCode:     "machine_code",
	HardwareVersion: "hardware_version",
	SoftwareVersion: "software_version",
	Cpu:             "cpu",
	Mem:             "mem",
	MemTotal:        "mem_total",
	Disk:            "disk",
	DiskTotal:       "disk_total",
	Address:         "address",
	CreateTime:      "create_time",
	ActivateTime:    "activate_time",
	LastOnlineTime:  "last_online_time",
	Deleted:         "deleted",
}

// NewDeviceDao creates and returns a new DAO object for table data access.
func NewDeviceDao(handlers ...gdb.ModelHandler) *DeviceDao {
	return &DeviceDao{
		group:    "default",
		table:    "device",
		columns:  deviceColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DeviceDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DeviceDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DeviceDao) Columns() DeviceColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DeviceDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DeviceDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DeviceDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
