package task

import (
	"SupperSystem/pkg/logger"
	"SupperSystem/internal/service"
	"SupperSystem/pkg/utils"
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
	// 停止cron
	ctx := t.cron.Stop()
	//等待
	select {
	case <-ctx.Done():
	// logger
	case <-time.After(30 * time.Second):
		// logger
	}
	t.running = false
	t.entries = make(map[string]cron.EntryID)
}

func (t *TaskManager) registerTasks() error {
	// 出库任务 -- 10分钟
	entryID, err := t.cron.AddFunc("*/10 * * * *", t.DeliveryWrappedTask)
	if err != nil {
		return fmt.Errorf("添加出库任务失败:%w", err)
	}
	t.entries["delivery"] = entryID

	// 科室退库任务 -- 30分钟
	entryID, err = t.cron.AddFunc("*/30 * * * *", t.RefundWrappedTask)
	if err != nil {
		return fmt.Errorf("添加科室退库任务失败:%w", err)
	}
	t.entries["refund"] = entryID

	//// 记录已注册任务
	//for name ,id := range t.entries {
	//	// 写日志
	//}
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
	now := time.Now()
	if utils.IsWithinWorkingTime() {
		logger.AsyncLog(fmt.Sprintf("\r\n事件:定时任务开始执行\r\n时间:%v\r\n%s\r\n",
			now.Format("2006-01-02 15:04:05"), logger.LoggerEndStr))
		service.DeliveryRetryService()
		logger.AsyncLog(fmt.Sprintf("\r\n事件:定时任务执行完成\r\n时间:%v\r\n%s\r\n",
			time.Now().Format("2006-01-02 15:04:05"), logger.LoggerEndStr))
	} else {
		logger.AsyncLog(fmt.Sprintf("\r\n事件:定时任务跳过\r\n时间:%v\r\n原因:非执行时间\r\n%s\r\n",
			now.Format("2006-01-02 15:04:05"), logger.LoggerEndStr))
	}
}

// RefundWrappedTask 退库定时执行
func (t *TaskManager) RefundWrappedTask() {
	now := time.Now()
	if utils.IsWithinWorkingTime() {
		logger.AsyncLog(fmt.Sprintf("\r\n事件:定时任务开始执行\r\n时间:%v\r\n%s\r\n",
			now.Format("2006-01-02 15:04:05"), logger.LoggerEndStr))
		service.RefundRetryService()
		logger.AsyncLog(fmt.Sprintf("\r\n事件:定时任务执行完成\r\n时间:%v\r\n%s\r\n",
			time.Now().Format("2006-01-02 15:04:05"), logger.LoggerEndStr))
	} else {
		logger.AsyncLog(fmt.Sprintf("\r\n事件:定时任务跳过\r\n时间:%v\r\n原因:非执行时间\r\n%s\r\n",
			now.Format("2006-01-02 15:04:05"), logger.LoggerEndStr))
	}
}
