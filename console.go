package logger

//在控制台输出日志
import (
	"fmt"
	"time"
)

//ConsoleLogger日志
type ConsoleLogger struct {
	Level LogLevel
}

//ConsoleNewLog构造函数
func ConsoleNewLog(levelStr string) ConsoleLogger {
	lever, err := parseLogLevel(levelStr)
	if err != nil {
		fmt.Printf("出现错误，%v", err)
	}
	return ConsoleLogger{
		Level: lever,
	}
}

//enable打印那些级别的日志
func (c *ConsoleLogger) enable(logLevel LogLevel) bool {
	return logLevel >= c.Level
}

//log打印日志
func (c *ConsoleLogger) log(level LogLevel, format string, a ...interface{}) {
	if c.enable(level) {
		nowTime := time.Now().Format("2006-01-02 15:04:05")
		fileName, funcName, lineNo := getInfo(3)
		msg := fmt.Sprintf(format, a...)
		fmt.Printf("[%s] [%s] [%s:%s:%d] %s\n", nowTime, getLogString(level), fileName, funcName, lineNo, msg)
	}

}
func (c *ConsoleLogger) Debug(format string, a ...interface{}) {
	c.log(DEBUG, format, a...)
}
func (c *ConsoleLogger) Trace(format string, a ...interface{}) {
	c.log(TRACE, format, a...)
}
func (c *ConsoleLogger) Info(format string, a ...interface{}) {
	c.log(INFO, format, a...)
}
func (c *ConsoleLogger) Warning(format string, a ...interface{}) {
	c.log(WARNING, format, a...)
}
func (c *ConsoleLogger) Error(format string, a ...interface{}) {
	c.log(ERROR, format, a...)
}
func (c *ConsoleLogger) Fatal(format string, a ...interface{}) {
	c.log(FATAL, format, a...)
}
