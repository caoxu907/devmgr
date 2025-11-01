// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UpgradeBatch is the golang structure of table upgrade_batch for DAO operations like Where/Data.
type UpgradeBatch struct {
	g.Meta               `orm:"table:upgrade_batch, do:true"`
	Id                   any         //
	Batchstatus          any         //
	PackageId            any         //
	ScopeType            any         //
	ScopeDevices         any         //
	UserConfirm          any         //
	RetryInterval        any         //
	MaxRetryCount        any         //
	OverrideExistingTask any         //
	UpgradeMode          any         //
	StrategyType         any         //
	GrayPercentage       any         //
	GraySuccessThreshold any         //
	GrayCompleted        any         //
	MaxConcurrent        any         //
	GroupInterval        any         //
	RandomDelay          any         //
	CreateTime           *gtime.Time //
	UpdateTime           *gtime.Time //
	StartTime            *gtime.Time //
	CompleteTime         *gtime.Time //
	Deleted              any         //
	TotalDevices         any         //
	SuccessDevices       any         //
	FailureDevices       any         //
	PendingDevices       any         //
}
