package task

import (
	conf "SupperSystem/configs"
	"SupperSystem/internal/service"
	"SupperSystem/pkg/logger"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type TaskManager struct {
	cron    *cron.Cron
	mu      sync.RWMutex
	running bool
	entries map[string]cron.EntryID
}

func NewTaskManager() *TaskManager {
	c := cron.New(
		cron.WithLogger(cron.PrintfLogger(log.New(io.Discard, "", 0))),
	)
	return &TaskManager{
		cron:    c,
		running: false,
		entries: make(map[string]cron.EntryID),
	}
}

func (t *TaskManager) Start() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.running {
		return fmt.Errorf("定时任务已经在运行")
	}
	// 注册所有定时任务
	if err := t.registerTasks(); err != nil {
		return fmt.Errorf("注册定时任务失败:%w", err)
	}
	// 启动cron
	t.cron.Start()
	t.running = true
	// 日志
	return nil
}

func (t *TaskManager) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.running {
		return
	}
	// 停止 cron
	ctx := t.cron.Stop()
	//等待
	select {
	case <-ctx.Done():
		logger.AsyncLog(fmt.Sprintf("定时任务已停止\r\n%s", logger.LoggerEndStr))
	case <-time.After(30 * time.Second):
		logger.AsyncLog(fmt.Sprintf("定时任务停止超时\r\n%s", logger.LoggerEndStr))
	}
	t.running = false
	t.entries = make(map[string]cron.EntryID)
}

func (t *TaskManager) registerTasks() error {
	// 获取Cron配置
	deliverySpec := conf.Configs.Cron.DeliverySchedule
	// 出库任务
	if deliverySpec == "" {
		deliverySpec = "*/10 8-20 * * *" // 默认：9-19点每10分钟
	}
	// 是否禁用
	if deliverySpec == "_" {
		logger.AsyncLog(fmt.Sprintf("定时任务[出库]已禁止(配置为-)\r\n%s", logger.LoggerEndStr))
	} else {
		entryID, err := t.cron.AddFunc(deliverySpec, t.DeliveryWrappedTask)
		if err != nil {
			return fmt.Errorf("添加出库任务失败:%w", err)
		}
		t.entries["delivery"] = entryID
		logger.AsyncLog(fmt.Sprintf("定时任务[出库]已注册，周期: %s\r\n%s", deliverySpec, logger.LoggerEndStr))
	}

	// 科室退库任务 -- 30分钟
	refundSpec := conf.Configs.Cron.RefundSchedule
	if refundSpec == "" {
		refundSpec = "*/30 9-19 * * * *" // 9-19点每30分钟
	}
	// 是否禁用
	if refundSpec == "_" {
		logger.AsyncLog(fmt.Sprintf("定时任务[退库]已禁止(配置为-)\r\n%s", logger.LoggerEndStr))
	} else {
		entryID, err := t.cron.AddFunc(refundSpec, t.RefundWrappedTask)
		if err != nil {
			return fmt.Errorf("添加科室退库任务失败:%w", err)
		}
		t.entries["refund"] = entryID
		logger.AsyncLog(fmt.Sprintf("定时任务[退库]已注册，周期: %s\r\n%s", refundSpec, logger.LoggerEndStr))
	}
	return nil
}

func (t *TaskManager) IsRunning() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.running
}

func (t *TaskManager) GetTaskEntries() map[string]cron.EntryID {
	t.mu.RLock()
	defer t.mu.RUnlock()
	result := make(map[string]cron.EntryID)
	for k, v := range t.entries {
		result[k] = v
	}
	return result
}

// DeliveryWrappedTask 出库定时执行
func (t *TaskManager) DeliveryWrappedTask() {
	// 移除 IsWithinWorkingTime 判断 依赖 Cron 表达式控制时间范围
	now := time.Now()
	logger.AsyncLog(fmt.Sprintf("\r\n事件:定时任务[出库]开始\r\n时间:%v\r\n%s\r\n",
		now.Format("2006-01-02 15:04:05"), logger.LoggerEndStr))
	service.DeliveryRetryService()
	logger.AsyncLog(fmt.Sprintf("\r\n事件:定时任务[出库]完成\r\n时间:%v\r\n%s\r\n",
		time.Now().Format("2006-01-02 15:04:05"), logger.LoggerEndStr))
}

// RefundWrappedTask 退库定时执行
func (t *TaskManager) RefundWrappedTask() {
	now := time.Now()
	logger.AsyncLog(fmt.Sprintf("\r\n事件:定时任务[退库]开始\r\n时间:%v\r\n%s\r\n",
		now.Format("2006-01-02 15:04:05"), logger.LoggerEndStr))
	service.RefundRetryService()
	logger.AsyncLog(fmt.Sprintf("\r\n事件:定时任务[退库]完成\r\n时间:%v\r\n%s\r\n",
		time.Now().Format("2006-01-02 15:04:05"), logger.LoggerEndStr))
}
