package service

import (
	"WorkloadQuery/controller"
	"WorkloadQuery/logger"
	"fmt"
)

// RefundRetryService 退库单重试
func RefundRetryService() {
	var logMsg string
	r := controller.RefundRequestInfo{
		Count: new(int64),
	}
	err := r.GetRefundNo()
	if err != nil {
		return
	}
	if *r.Count == 0 {
		return
	}
	logMsg = fmt.Sprintf("\r\n事件:查询领用出库失败业务数据\r\n查询结果:%+v\r\n%v\r\n%s\r\n", *r.Count, *r.Re, logger.LoggerEndStr)
	logger.AsyncLog(logMsg)
	//接口调用
	err = r.RetryRefundToHis()
	if err != nil {
		logMsg = fmt.Sprintf("\r\n:事件：接口调用错误:%v\r\n%s\r\n", err.Error(), logger.LoggerEndStr)
		return
	}
}
