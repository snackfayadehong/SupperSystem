package service

import (
	"SupperSystem/internal/controller"
	"SupperSystem/pkg/logger"
	"fmt"
	"time"
)

// DeliveryRetryService 领用出库单重试
func DeliveryRetryService() {
	var logMsg string
	var startDate string
	var endDate string
	d := controller.DeliveryRequestInfo{
		Count: new(int64),
	}
	var now = time.Now()
	s := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()) // 当天0时0点
	e := now.Add(-10 * time.Minute)                                                // 当前时间前推10分钟
	startDate = s.Format("2006-01-02 15:04:05")
	endDate = e.Format("2006-01-02 15:04:05")
	err := d.GetDeliveryNo(startDate, endDate)
	if err != nil {
		return
	}
	if *d.Count == 0 {
		return
	}
	logMsg = fmt.Sprintf("\r\n事件:查询领用出库失败业务数据\r\n查询结果:%+v\r\n%v\r\n%s\r\n", *d.Count, *d.De, logger.LoggerEndStr)
	logger.AsyncLog(logMsg)
	//接口调用
	err = d.DeliveryNoRetryToHis()
	if err != nil {
		logMsg = fmt.Sprintf("\r\n:事件：接口调用错误:%v\r\n%s\r\n", err.Error(), logger.LoggerEndStr)
		return
	}
}
