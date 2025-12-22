package logger

import (
	"SupperSystem/configs"
	"SupperSystem/pkg/utils"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const LoggerEndStr = "----------------------------------------------------------------------------"
const InfoLevel = zapcore.InfoLevel
const ErrorLevel = zapcore.ErrorLevel

// 异步日志处理器
type logEntry struct {
	level     zapcore.Level
	message   string
	fields    []zap.Field
	caller    string
	timestamp time.Time
}
type asyncLogger struct {
	logCh chan logEntry
	done  chan struct{}
	wg    sync.WaitGroup
}

var (
	asyncLog *asyncLogger
	once     sync.Once
)

// 新建异步日志处理器
func newAsyncLogger(bufferSize int) *asyncLogger {
	logger := &asyncLogger{
		logCh: make(chan logEntry, bufferSize),
		done:  make(chan struct{}),
	}
	logger.wg.Add(1)
	go logger.worker()
	return logger
}

// worker 处理日志的worker
func (a *asyncLogger) worker() {
	defer a.wg.Done()
	for {
		select {
		case entry := <-a.logCh:
			a.safeLog(entry)
		case <-a.done:
			// 处理剩余日志
			for {
				select {
				case entry := <-a.logCh:
					a.safeLog(entry)
				default:
					return
				}
			}
		}
	}
}

func (a *asyncLogger) safeLog(entry logEntry) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		//输出到标准错误
	//		fmt.Fprintf(os.Stderr, "Log panic recovered:\n%v\n\n", r)
	//	}
	//}()
	//switch entry.level {
	//case zapcore.DebugLevel:
	//	zap.L().Debug(entry.message, entry.fields...)
	//case zapcore.InfoLevel:
	//	entry.message = "关机空日志"
	//	zap.L().Info(entry.message, entry.fields...)
	//case zapcore.WarnLevel:
	//	zap.L().Warn(entry.message, entry.fields...)
	//case zapcore.ErrorLevel:
	//	zap.L().Error(entry.message, entry.fields...)
	//default:
	//	zap.L().Info(entry.message, entry.fields...)
	//}
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Log panic recovered:\n%v\n\n", r)
		}
	}()
	// 在日志中添加调用者信息
	allFields := entry.fields
	if entry.caller != "" {
		allFields = append(allFields, zap.String("caller", entry.caller))
	}

	switch entry.level {
	case zapcore.DebugLevel:
		zap.L().Debug(entry.message, allFields...)
	case zapcore.InfoLevel:
		zap.L().Info(entry.message, allFields...)
	case zapcore.WarnLevel:
		zap.L().Warn(entry.message, allFields...)
	case zapcore.ErrorLevel:
		zap.L().Error(entry.message, allFields...)
	default:
		zap.L().Info(entry.message, allFields...)
	}
}

// Close 关闭异步日志处理器
func Close() {
	if asyncLog == nil {
		return
	}
	close(asyncLog.done)
	asyncLog.wg.Wait()
	_ = zap.L().Sync()
}

// InitLog 日志
func InitLog() (err error) {
	// 创建日志目录
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err = os.Mkdir("logs", 0777); err != nil {
			return err
		}
	}
	// 初始化异步日志处理器
	once.Do(func() {
		asyncLog = newAsyncLogger(1000) // 缓冲1000条日志
	})
	// 创建Core三大件，进行初始化
	encoder := getEncoder()
	// 判断日志级别
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // error级别
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // info和debug级别,debug级别是最低的
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})
	infoSyncer := getLogWriter("logs/", InfoLevel)
	errorSyncer := getLogWriter("logs/", ErrorLevel)
	// 创建核心-->如果是debug模式，就在控制台和文件都打印，否则就只写到文件中
	var core zapcore.Core
	if conf.Configs.Server.RunModel == "debug" {
		// 开发模式，日志输出到终端
		// NewTee创建一个核心，将日志条目复制到两个或多个底层核心中。
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), InfoLevel),
		)
	} else {
		coreInfo := zapcore.NewCore(encoder, infoSyncer, lowPriority)
		coreErr := zapcore.NewCore(encoder, errorSyncer, highPriority)
		core = zapcore.NewTee(coreInfo, coreErr)
	}
	// 创建 logger 对象
	// Warn()/Error()等级别的日志会输出堆栈，Debug()/Info()这些级别不会
	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.WarnLevel))
	// 替换全局的 logger, 后续在其他包中只需使用zap.L()调用即可
	zap.ReplaceGlobals(log)
	return nil
}

// 获取调用者信息的辅助函数
func getCallerInfo(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}

	// 获取函数名
	funcName := runtime.FuncForPC(pc).Name()
	// 简化文件路径
	if idx := strings.LastIndex(file, "/"); idx >= 0 {
		file = file[idx+1:]
	}

	return fmt.Sprintf("%s:%d %s", file, line, funcName)
}

// AsyncLog 异步日志
func AsyncLog(logMsg string) {
	if asyncLog == nil {
		// 降级成同步日志
		zap.L().Info(logMsg)
		return
	}
	caller := getCallerInfo(2)
	entry := logEntry{
		level:   zapcore.InfoLevel,
		message: logMsg,
		//caller:    caller,
		timestamp: time.Now(),
	}
	select {
	case asyncLog.logCh <- entry:
	//成功
	default:
		//通道已满,降级为同步
		zap.L().Warn("异步日志通道已满,降级为同步日志")
		zap.L().Info(logMsg, zap.String("caller", caller))
	}
}

// AsyncLogWithFields 结构化日志
func AsyncLogWithFields(level zapcore.Level, msg string, fields ...zap.Field) {
	if asyncLog == nil {
		// 降级为同步日志
		switch level {
		case zapcore.DebugLevel:
			zap.L().Debug(msg, fields...)
		case zapcore.InfoLevel:
			zap.L().Info(msg, fields...)
		case zapcore.WarnLevel:
			zap.L().Warn(msg, fields...)
		case zapcore.ErrorLevel:
			zap.L().Error(msg, fields...)
		}
	}
	formattedMsg := fmt.Sprintf("%s\r\n%s\r\n", msg, LoggerEndStr)
	// 异步记录主日志
	entry := logEntry{
		level:   level,
		message: formattedMsg,
		fields:  fields,
	}
	select {
	case asyncLog.logCh <- entry:
		// 成功写入
	default:
		// 通道已满，降级为同步日志
		zap.L().Warn("异步日志通道已满，降级为同步日志")
		switch level {
		case zapcore.DebugLevel:
			zap.L().Debug(msg, fields...)
		case zapcore.InfoLevel:
			zap.L().Info(msg, fields...)
		case zapcore.WarnLevel:
			zap.L().Warn(msg, fields...)
		case zapcore.ErrorLevel:
			zap.L().Error(msg, fields...)
		}
		zap.L().Info(fmt.Sprintf("\r\n%s\r\n", LoggerEndStr))
		return
	}
}

// 获取Encoder，给初始化logger使用的
func getEncoder() zapcore.Encoder {
	// 使用zap提供的 NewProductionEncoderConfig
	encoderConfig := zap.NewProductionEncoderConfig()
	// 设置时间格式
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	// 时间的key
	encoderConfig.TimeKey = "time"
	// 级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 显示调用者信息
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	// 返回json 格式的 日志编辑器
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 获取切割的问题，给初始化logger使用的
func getLogWriter(filename string, leavel zapcore.Level) zapcore.WriteSyncer {
	var logFileName string
	switch leavel {
	case zapcore.ErrorLevel:
		logFileName = filename + "ERROR_%Y%m%d.log"
	case zapcore.InfoLevel:
		logFileName = filename + "%Y%m%d.log"
	default:
		logFileName = filename + "Other_%Y%m%d.log"
	}
	// 日志轮转前清除 .symlink
	suffix := ".log_symlink"
	if err := utils.RemoveAssignDir(filename, suffix); err != nil {
		zap.L().Error("Error", zap.Error(err))
	}
	// 日志轮转
	hook, err := rotatelogs.New(
		logFileName,
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(0),
		rotatelogs.WithRotationTime(time.Hour*24))
	if err != nil {
		zap.L().Error("ERROR", zap.Error(err))
	}
	return zapcore.AddSync(hook)
}

// GinLogger 用于替换gin框架的Logger中间件，不传参数，直接这样写
func GinLogger(c *gin.Context) {
	start := time.Now()
	var reqData string
	skipPaths := []string{
		"/health",
		"/debug/pprof/",
	}
	for _, skipPath := range skipPaths {
		if strings.HasPrefix(c.Request.URL.Path, skipPath) {
			zap.L().Info(fmt.Sprintf("%s,%s", c.ClientIP(), c.RemoteIP()))
			return
		}
	}
	if shouldLogBody(c) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err == nil {
			reqData = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}
	c.Next()
	// 视图函数执行完成，统计时间，记录日志
	cost := time.Since(start)
	status := c.Writer.Status()
	//sugar.Infof("\r\n事件:接口调用\r\nIP：%s\r\n代理IP:%s\r\nURL：%s\r\nMethod：%s\r\n入参：%s\r\nError：%s\r\nCost：%s\r\n%s\r\n",
	//	c.ClientIP(), c.RemoteIP(), c.Request.URL.Path, c.Request.Method, reqData,
	//	c.Errors.ByType(gin.ErrorTypePrivate).String(), cost, LoggerEndStr)
	//logMsg := fmt.Sprintf("\r\n事件:接口调用\r\nIP：%s\r\n代理IP:%s\r\nURL：%s\r\nMethod：%s\r\n入参：%s\r\nError：%s\r\nCost：%s\r\n%s\r\n",
	//	c.ClientIP(), c.RemoteIP(), c.Request.URL.Path, c.Request.Method, reqData,
	//	c.Errors.ByType(gin.ErrorTypePrivate).String(), cost, LoggerEndStr)
	//AsyncLog(logMsg)
	// 结构化日志
	logFields := []zap.Field{
		zap.String("event", "接口调用"),
		zap.String("IP", c.ClientIP()),
		zap.String("代理IP", c.RemoteIP()),
		zap.String("url", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("user_agent", c.Request.UserAgent()),
		zap.Int("status", status),
		zap.String("cost", cost.String()),
		zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
	}
	if reqData != "" {
		logFields = append(logFields, zap.String("request_body", reqData))
	}
	// 根据状态码选择日志级别
	if status >= http.StatusInternalServerError {
		AsyncLogWithFields(zapcore.ErrorLevel, "HTTP Server Error", logFields...)
	} else if status >= http.StatusBadRequest {
		AsyncLogWithFields(zapcore.WarnLevel, "HTTP Client Error", logFields...)
	} else {
		AsyncLogWithFields(zapcore.InfoLevel, "HTTP Request", logFields...)
	}
}

func shouldLogBody(c *gin.Context) bool {
	// 根据内容类型和路径决定记录 body
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		return false // 不记录文件上传
	}
	// excludedPath 排除一些接口请求日志 /health 健康检查
	excludePaths := []string{"/health", "/metrics", "/debug/pprof"}
	for _, path := range excludePaths {
		if strings.HasPrefix(c.Request.URL.Path, path) {
			return false
		}
	}
	return c.Request.ContentLength > 0 && c.Request.ContentLength < 1024*50 // 只记录小于60KB的body
}

// GinRecovery 用于替换gin框架的Recovery中间件，因为传入参数，再包一层
func GinRecovery(stack bool) gin.HandlerFunc {
	logger := zap.L()
	return func(c *gin.Context) {
		defer func() {
			// defer 延迟调用，出了异常，处理并恢复异常，记录日志
			if err := recover(); err != nil {
				//  这个不必须，检查是否存在断开的连接(broken pipe或者connection reset by peer)---------开始--------
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				// http util包预先准备好的DumpRequest方法
				httpRequest, err := httputil.DumpRequest(c.Request, false)
				if err != nil {
					logger.Error("ERROR", zap.Error(err))
				}
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// 如果连接已断开，我们无法向其写入状态
					c.Error(err.(error))
					c.Abort()
					return
				}
				//  这个不必须，检查是否存在断开的连接(broken pipe或者connection reset by peer)
				// 是否打印堆栈信息，使用的是debug.Stack()，传入false，在日志中就没有堆栈信息
				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				// 有错误，直接返回给前端错误，前端直接报错
				// c.AbortWithStatus(http.StatusInternalServerError)
				// 该方式前端不报错
				c.String(200, "访问出错了")
			}
		}()
		c.Next()
	}
}
