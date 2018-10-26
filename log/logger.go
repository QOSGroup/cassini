package log

import (
	"github.com/cihub/seelog"
)

var dfaultConfig = `
<seelog minlevel="info">
	<outputs formatid="formater"><console />
	</outputs>
	<formats>
		<format id="formater" format="[%Date(2006-01-02 15:04:05.000000000)][%LEV] %Msg%n"/>
	</formats>
</seelog>`

func init() {
	logger, _ := seelog.LoggerFromConfigAsBytes([]byte(dfaultConfig))
	Replace(logger)
}

// Replace logger
func Replace(logger seelog.LoggerInterface) {
	seelog.ReplaceLogger(logger)
}

// Flush immediately processes all currently queued logs.
func Flush() {
	seelog.Flush()
}

// Trace logs 详细运行跟踪日志，可能影响程序性能，所以生产环境不配置输出，仅在开发测试环境使用
func Trace(v ...interface{}) {
	seelog.Trace(v...)
}

// Debug logs
func Debug(v ...interface{}) {
	seelog.Debug(v...)
}

// Info logs
func Info(v ...interface{}) {
	seelog.Info(v...)
}

// Error logs
func Error(v ...interface{}) {
	seelog.Error(v...)
}

// Tracef logs 详细运行跟踪日志，可能影响程序性能，所以生产环境不配置输出，仅在开发测试环境使用
func Tracef(format string, params ...interface{}) {
	seelog.Tracef(format, params...)
}

// Debugf formats logs
func Debugf(format string, params ...interface{}) {
	seelog.Debugf(format, params...)
}

// Infof formats logs
func Infof(format string, params ...interface{}) {
	seelog.Infof(format, params...)
}

// Errorf formats logs
func Errorf(format string, params ...interface{}) {
	seelog.Errorf(format, params...)

}
