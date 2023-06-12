package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

// InitLog 日志
func InitLog() (logFile *os.File, logConfig *gin.LoggerConfig, err error) {
	_, err = os.Stat("logs")
	if os.IsNotExist(err) {
		err = os.Mkdir("logs", 0700)
	}
	logTime := time.Now().Format("2006-01-02")
	_, err = os.Stat(fmt.Sprintf("logs/%s.log", logTime))
	if os.IsNotExist(err) {
		logFile, err = os.OpenFile(fmt.Sprintf("logs/%s.log", logTime), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return
		}
	}
	// gin日志配置
	logConfig = &gin.LoggerConfig{
		Formatter: func(params gin.LogFormatterParams) string {
			return fmt.Sprintf("客户端IP:%s,请求时间:[%s],请求方式:%s,请求地址:%s,http协议版本:%s,请求状态码:%d,响应时间:%s,客户端:%s,错误信息:%s\r\n",
				params.ClientIP,
				params.TimeStamp.Format("2006-01-02 15:04:05"),
				params.Method,
				params.Path,
				params.Request.Proto,
				params.StatusCode,
				params.Latency,
				params.Request.UserAgent(),
				params.ErrorMessage,
			)
		},
		Output: logFile,
	}
	return
}
