// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UpgradeBatch is the golang structure for table upgrade_batch.
type UpgradeBatch struct {
	Id                   int         `json:"id"                   orm:"id"                     description:""` //
	Batchstatus          int         `json:"batchstatus"          orm:"batchstatus"            description:""` //
	PackageId            int         `json:"packageId"            orm:"package_id"             description:""` //
	ScopeType            int         `json:"scopeType"            orm:"scope_type"             description:""` //
	ScopeDevices         string      `json:"scopeDevices"         orm:"scope_devices"          description:""` //
	UserConfirm          bool        `json:"userConfirm"          orm:"user_confirm"           description:""` //
	RetryInterval        int         `json:"retryInterval"        orm:"retry_interval"         description:""` //
	MaxRetryCount        int         `json:"maxRetryCount"        orm:"max_retry_count"        description:""` //
	OverrideExistingTask bool        `json:"overrideExistingTask" orm:"override_existing_task" description:""` //
	UpgradeMode          int         `json:"upgradeMode"          orm:"upgrade_mode"           description:""` //
	StrategyType         int         `json:"strategyType"         orm:"strategy_type"          description:""` //
	GrayPercentage       int         `json:"grayPercentage"       orm:"gray_percentage"        description:""` //
	GraySuccessThreshold int         `json:"graySuccessThreshold" orm:"gray_success_threshold" description:""` //
	GrayCompleted        bool        `json:"grayCompleted"        orm:"gray_completed"         description:""` //
	MaxConcurrent        int         `json:"maxConcurrent"        orm:"max_concurrent"         description:""` //
	GroupInterval        int         `json:"groupInterval"        orm:"group_interval"         description:""` //
	RandomDelay          int         `json:"randomDelay"          orm:"random_delay"           description:""` //
	CreateTime           *gtime.Time `json:"createTime"           orm:"create_time"            description:""` //
	UpdateTime           *gtime.Time `json:"updateTime"           orm:"update_time"            description:""` //
	StartTime            *gtime.Time `json:"startTime"            orm:"start_time"             description:""` //
	CompleteTime         *gtime.Time `json:"completeTime"         orm:"complete_time"          description:""` //
	Deleted              bool        `json:"deleted"              orm:"deleted"                description:""` //
	TotalDevices         int         `json:"totalDevices"         orm:"total_devices"          description:""` //
	SuccessDevices       int         `json:"successDevices"       orm:"success_devices"        description:""` //
	FailureDevices       int         `json:"failureDevices"       orm:"failure_devices"        description:""` //
	PendingDevices       int         `json:"pendingDevices"       orm:"pending_devices"        description:""` //
}
