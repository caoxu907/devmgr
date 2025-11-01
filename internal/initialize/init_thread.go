package initialize

import (
	"context"
	"devmgr/internal/consts"
	"devmgr/internal/service"
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
)

func InitThread(ctx context.Context) error {
	initDeviceCheck(ctx)
	err := initMqClient(ctx)
	if err != nil {
		return err
	}

	return nil
}

func StopThread() {
	closeDeviceCheck()
	closeMqClient()
}

// ----MQ客户端----
func initMqClient(ctx context.Context) error {

	client := service.MQClient()
	if client == nil {
		g.Log().Line().Error(ctx, "MQClient初始化失败")
		return fmt.Errorf("MQClient初始化失败")
	}

	//处理设备通用消息
	handleCommonMsg := func(topic string, key []byte, payload []byte) {
		g.Log().Line().Infof(ctx, "handleCommonMsg收到消息: topic=%s, key=%s, value=%s", topic, string(key), string(payload))
		if len(key) != 0 && len(payload) != 0 {
			_ = service.DevApi().ParseDeviceMessage(ctx, string(key), string(payload))

		}
	}
	client.RegisterSubscription(consts.Kafka_topic_notify_up, handleCommonMsg)

	// 启动Kafka客户端
	if err := client.Start(); err != nil {
		g.Log().Fatal(ctx, err)
		return err
	}

	return nil
}

func closeMqClient() {
	client := service.MQClient()
	if client != nil {
		client.Close()
	}
}

// ----设备检查 ----
var (
	deviceCheckCancel  context.CancelFunc
	deviceCheckCtx     context.Context
	deviceCheckMutex   sync.Mutex            // 添加互斥锁
	deviceCheckLockKey = "device_check_lock" // Redis锁的key
	deviceCheckLockTTL = 40                  // 锁的过期时间(秒)
	devStatusInterval  = 10                  //设备状态检查间隔（秒）
	otaStatusInterval  = 10                  //// OTA升级状态检查间隔（秒）
)

func initDeviceCheck(ctx context.Context) {
	initDeviceCheckWithLock(ctx)
	//initDeviceCheckWithoutLock(ctx)
}

func initDeviceCheckWithLock(ctx context.Context) {
	deviceCheckMutex.Lock()
	defer deviceCheckMutex.Unlock()

	deviceCheckCtx, deviceCheckCancel = context.WithCancel(ctx)

	// 启动单独的goroutine进行检查
	go func(ctx context.Context) {
		g.Log().Line().Info(ctx, "设备检查任务启动")
		defer g.Log().Line().Info(ctx, "设备检查任务退出")

		// 记录上次执行时间
		lastDevStatusCheck := time.Now()
		lastOtaStatusCheck := time.Now()
		lastLockRenew := time.Now()

		// 锁续约间隔（秒），应小于锁的TTL
		lockRenewInterval := deviceCheckLockTTL / 2

		// 尝试获取分布式锁
		redis := g.Redis()
		lockValue := fmt.Sprintf("%d", time.Now().UnixNano())
		hasLock := false

		// 首次尝试获取锁
		if acquireLock(ctx, redis, lockValue) {
			hasLock = true
			g.Log().Line().Info(ctx, "成功获取设备检查锁")
		} else {
			g.Log().Line().Info(ctx, "未能获取设备检查锁，进入观察模式")
		}

		for {
			select {
			case <-ctx.Done():
				g.Log().Line().Info(ctx, "设备检查任务已停止")
				// 释放锁
				if hasLock {
					releaseLock(ctx, redis, lockValue)
				}
				return
			case <-time.After(5 * time.Second): // 每5秒检查一次
				now := time.Now()

				// 如果没有锁，尝试获取锁
				if !hasLock {
					if acquireLock(ctx, redis, lockValue) {
						hasLock = true
						g.Log().Line().Info(ctx, "成功获取设备检查锁")
					}
				}

				// 没有锁则跳过后续检查
				if !hasLock {
					continue
				}

				// 锁续约
				if hasLock && now.Sub(lastLockRenew).Seconds() >= float64(lockRenewInterval) {
					if !renewLock(ctx, redis, lockValue) {
						g.Log().Line().Warning(ctx, "设备检查锁续约失败，可能已被其他实例获取")
						hasLock = false
					} else {
						g.Log().Debug(ctx, "设备检查锁续约成功")
						lastLockRenew = now
					}
				}

				// 只有获得锁的实例才执行检查
				if hasLock {
					// 检查设备状态
					if now.Sub(lastDevStatusCheck).Seconds() >= float64(devStatusInterval) {
						g.Log().Debug(ctx, "执行设备状态检查")
						service.DevApi().DevStatusCheck(ctx)
						service.DevApi().CommandTimeoutCheck(ctx)
						lastDevStatusCheck = now
					}

					// 检查OTA升级状态
					if now.Sub(lastOtaStatusCheck).Seconds() >= float64(otaStatusInterval) {
						//g.Log().Debug(ctx, "执行OTA升级状态检查")
						service.Ota().CheckUpgradeTaskStatus(ctx)
						lastOtaStatusCheck = now
					}
				}
			}
		}
	}(deviceCheckCtx)

	return
}

//func initDeviceCheckWithoutLock(ctx context.Context) {
//	deviceCheckMutex.Lock()
//	defer deviceCheckMutex.Unlock()
//
//	deviceCheckCtx, deviceCheckCancel = context.WithCancel(ctx)
//
//	// 启动单独的goroutine进行检查
//	go func(ctx context.Context) {
//		g.Log().Line().Info(ctx, "设备检查任务启动")
//		defer g.Log().Line().Info(ctx, "设备检查任务退出")
//
//		// 记录上次执行时间
//		lastDevStatusCheck := time.Now()
//		lastOtaStatusCheck := time.Now()
//
//		for {
//			select {
//			case <-ctx.Done():
//				g.Log().Line().Info(ctx, "设备检查任务已停止")
//				return
//			case <-time.After(5 * time.Second): // 每5秒检查一次
//				now := time.Now()
//
//				// 检查设备状态
//				if now.Sub(lastDevStatusCheck).Seconds() >= float64(devStatusInterval) {
//					g.Log().Debug(ctx, "执行设备状态检查")
//					service.DevApi().DevStatusCheck(ctx)
//					lastDevStatusCheck = now
//				}
//
//				// 检查OTA升级状态
//				if now.Sub(lastOtaStatusCheck).Seconds() >= float64(otaStatusInterval) {
//					//g.Log().Debug(ctx, "执行OTA升级状态检查")
//					//service.Ota().CheckUpgradeTaskStatus(ctx)
//					lastOtaStatusCheck = now
//				}
//			}
//		}
//	}(deviceCheckCtx)
//
//	return
//}

/*
"SET" - Redis 的 SET 命令，用于设置键值对
deviceCheckLockKey - 锁的键名，在代码中定义为 "device_check_lock"
value - 锁的值，一般是唯一标识符（代码中使用的是纳秒级时间戳）
"NX" - 表示 "Not eXists"，仅当键不存在时才设置
"EX" - 表示 "EXpire"，设置键的过期时间（单位：秒）
deviceCheckLockTTL - 锁的过期时间，代码中设为 15 秒

原子性操作：Redis 的 SET 命令带 NX 和 EX 选项是原子性的，整个检查键是否存在并设置值的过程不会被打断。
互斥性：由于 NX 选项，只有当键不存在时才能设置成功。这确保了在任何时刻只有一个客户端能成功获取锁。
防止死锁：通过 EX 设置过期时间，即使持有锁的客户端崩溃，锁也会在一定时间后自动释放。
锁的标识：通过设置唯一的 value，确保只有获取锁的客户端才能释放锁，防止误释放
*/

// acquireLock 获取Redis分布式锁
func acquireLock(ctx context.Context, redis *gredis.Redis, value string) bool {
	// 使用SET命令搭配NX和EX选项
	result, err := redis.Do(ctx, "SET", deviceCheckLockKey, value, "NX", "EX", deviceCheckLockTTL)
	if err != nil {
		g.Log().Line().Info(ctx, "获取分布式锁出错:", err)
		return false
	}
	return result.String() == "OK"
}

// renewLock 续约Redis分布式锁
func renewLock(ctx context.Context, redis *gredis.Redis, value string) bool {
	// 使用Lua脚本确保原子性操作：只有当前值匹配时才更新过期时间
	script := `
        if redis.call('get', KEYS[1]) == ARGV[1] then
            return redis.call('expire', KEYS[1], ARGV[2])
        else
            return 0
        end
    `
	// 正确使用 Eval 方法，numKeys=1 表示使用1个键
	result, err := redis.Eval(ctx, script, 1, []string{deviceCheckLockKey}, []interface{}{value, deviceCheckLockTTL})
	if err != nil {
		g.Log().Line().Error(ctx, "续约分布式锁出错:", err)
		return false
	}
	return result.Int() == 1
}

// releaseLock 释放Redis分布式锁
func releaseLock(ctx context.Context, redis *gredis.Redis, value string) bool {
	// 使用Lua脚本确保原子性操作：只删除自己的锁
	script := `
        if redis.call('get', KEYS[1]) == ARGV[1] then
            return redis.call('del', KEYS[1])
        else
            return 0
        end
    `
	// 正确使用 Eval 方法，numKeys=1 表示使用1个键
	result, err := redis.Eval(ctx, script, 1, []string{deviceCheckLockKey}, []interface{}{value})
	if err != nil {
		g.Log().Line().Error(ctx, "释放分布式锁出错:", err)
		return false
	}
	return result.Int() == 1
}

func closeDeviceCheck() {
	deviceCheckMutex.Lock()
	defer deviceCheckMutex.Unlock()

	if deviceCheckCancel != nil {
		deviceCheckCancel()
	}
}
