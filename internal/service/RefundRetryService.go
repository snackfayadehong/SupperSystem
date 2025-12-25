package service

import (
	"SupperSystem/internal/controller"
	"SupperSystem/pkg/logger"
	"fmt"
	"time"
)

// RefundRetryService 退库单重试
func RefundRetryService() {
	var logMsg string
	r := controller.RefundRequestInfo{
		Count: new(int64),
	}
	var now = time.Now()
	s := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()) // 当天0时0点
	e := now.Add(-10 * time.Minute)                                                // 当前时间前推10分钟
	startDate := s.Format("2006-01-02 15:04:05")
	endDate := e.Format("2006-01-02 15:04:05")
	err := r.GetRefundNo(startDate, endDate)
	if err != nil {
		return
	}
	if *r.Count == 0 {
		return
	}
	logMsg = fmt.Sprintf("\r\n事件:查询科室退库失败业务数据\r\n查询结果:%+v\r\n%v\r\n%s\r\n", *r.Count, *r.Re, logger.LoggerEndStr)
	logger.AsyncLog(logMsg)
	//接口调用
	err = r.RetryRefundToHis()
	if err != nil {
		logMsg = fmt.Sprintf("\r\n:事件：接口调用错误:%v\r\n%s\r\n", err.Error(), logger.LoggerEndStr)
		return
	}
}
