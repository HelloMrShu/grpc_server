package global

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func InitLogger() {
	writer := getLogWriter()
	encoder := getEncoder()

	core := zapcore.NewCore(encoder, writer, zapcore.InfoLevel)
	Logger = zap.New(core, zap.AddCaller())
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	filename := fmt.Sprintf("%s/%s", ServerConfig.Log.Path, ServerConfig.Log.File)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename, //日志文件名称
		MaxSize:    50,       //文件大小限制，默认MB
		MaxBackups: 10,       //备份文件数量
		MaxAge:     30,       //最大天数
		Compress:   false,    //是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

// LoggerMiddleware for request
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		info := make([]string, 0)
		info = append(info, strconv.Itoa(c.Writer.Status()))
		info = append(info, c.Request.Method)
		info = append(info, c.Request.URL.Path)
		info = append(info, c.Request.URL.RawQuery)
		info = append(info, c.ClientIP())

		c.Next()

		timeDiff := time.Since(start)
		info = append(info, timeDiff.String())
		info = append(info, c.Errors.ByType(gin.ErrorTypePrivate).String())
		data := strings.Join(info, "  ")

		Logger.Info(data)
	}
}

func RecoveryMiddleware(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

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
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
