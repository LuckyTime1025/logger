package logger

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
)

//Logger接口
type Logger interface {
	Debug(format string, a ...interface{})
	Trace(format string, a ...interface{})
	Info(format string, a ...interface{})
	Warning(format string, a ...interface{})
	Error(format string, a ...interface{})
	Fatal(format string, a ...interface{})
}

//LogLevel日志级别
type LogLevel uint16

const (
	//定义日志级别
	UNKNOWN LogLevel = iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

//parseLogLevel解析日志级别
func parseLogLevel(levelStr string) (LogLevel, error) {
	levelStr = strings.ToUpper(levelStr)
	switch levelStr {
	case "DEBUG":
		return DEBUG, nil
	case "TRACE":
		return TRACE, nil
	case "INFO":
		return INFO, nil
	case "WARNING":
		return WARNING, nil
	case "ERROR":
		return ERROR, nil
	case "FATAL":
		return FATAL, nil
	}
	return UNKNOWN, errors.New("日志级别错误")
}

//获取日志级别字符串
func getLogString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}
	return "DEBUG"
}

//getInfo获取文件名，方法名，行号
func getInfo(skip int) (fileName, funcName string, lineNumber int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		fmt.Printf("runtime.Caller failed\n")
		return
	}
	funcName = strings.Split(runtime.FuncForPC(pc).Name(), ".")[1]
	fileName = path.Base(file)
	return funcName, fileName, line
}
