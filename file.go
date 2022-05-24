package logger

//在文件中输出日志
import (
	"fmt"
	"os"
	"path"
	"time"
)

var (
	//MaxSize日志通道缓冲区大小
	MaxSize = 50000
)

//文件日志结构体
type FileLogger struct {
	Level       LogLevel
	filePath    string //日志文件路径
	fileName    string //日志文件名称
	fileObj     *os.File
	errFileObj  *os.File
	maxFileSize int64 //最大文件大小
	//按时间切割
	// loggerTime  int
	logChan chan *logMsg
}

type logMsg struct {
	level      LogLevel
	lineNumber int
	msg        string
	timeStamp  string
	funcName   string
	fileName   string
}

func FileNewLog(levelStr, filePath, fileName string, maxSize int64) *FileLogger {
	lever, err := parseLogLevel(levelStr)
	if err != nil {
		fmt.Printf("出现错误，%v", err)
	}
	file := &FileLogger{
		Level:       lever,
		filePath:    filePath,
		fileName:    fileName,
		maxFileSize: maxSize,
		logChan:     make(chan *logMsg, MaxSize),
	}
	err = file.initFile() //按照文件路径和文件名将文件打开
	if err != nil {
		panic(err)
	}
	return file
}

//initFile初始化文件
func (f *FileLogger) initFile() error {
	//fullFileName文件全路径
	fullFileName := path.Join(f.filePath, f.fileName)
	fileObj, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("打开文件失败，err:%v\n", err)
		return err
	}
	eerFileObj, err := os.OpenFile(fullFileName+".err", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("打开err文件失败, err:%v\n", err)
		return err
	}
	f.fileObj = fileObj
	f.errFileObj = eerFileObj
	//开启一个后台的goroutine去写日志
	for i := 0; i < 5; i++ {
		go f.writeLogBackground()
	}
	return nil
}

//checkSize获取文件大小
func (f *FileLogger) checkSize(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("获取文件大小出错，err:%v", err)
	}
	return fileInfo.Size() >= f.maxFileSize
}

//checkTime获取文件时间
// func (f *FileLogger) checkTime() bool {
// 	return time.Now().Minute() != f.loggerTime
// }

//cutFile文件切割
func (f *FileLogger) cutFile(file *os.File) (*os.File, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	logname := path.Join(f.filePath, fileInfo.Name())
	nowStr := time.Now().Format("2006010215040500")
	newlogname := fmt.Sprintf("%s.bak%s", logname, nowStr)
	file.Close()
	os.Rename(logname, newlogname)
	fileObj, err := os.OpenFile(logname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return fileObj, nil
}

//enable打印那些级别的日志
func (f *FileLogger) enable(logLevel LogLevel) bool {
	return logLevel >= f.Level
}

// func (f *FileLogger) log(level LogLevel, format string, a ...interface{}) {
// 	if f.enable(level) {
// 		nowTime := time.Now().Format("2006-01-02 15:04:05")
// 		fileName, funcName, lineNo := getInfo(3)
// 		msg := fmt.Sprintf(format, a...)
// 		f.loggerTime = time.Now().Minute()
// 		if f.checkTime() {
// 			newFile, err := f.cutFile(f.fileObj)
// 			if err != nil {
// 				fmt.Printf("文件切割失败，err:%v", err)
// 				return
// 			}
// 			f.fileObj = newFile
// 		}
// 		fmt.Fprintf(f.fileObj, "[%s] [%s] [%s:%s:%d] %s\n", nowTime, getLogString(level), fileName, funcName, lineNo, msg)
// 		if level >= ERROR {
// 			newFile, err := f.cutFile(f.errFileObj)
// 			if err != nil {
// 				fmt.Printf("文件切割失败，err:%v", err)
// 				return
// 			}
// 			f.errFileObj = newFile
// 			fmt.Fprintf(f.errFileObj, "[%s] [%s] [%s:%s:%d] %s\n", nowTime, getLogString(level), fileName, funcName, lineNo, msg)
// 		}
// 	}
// }

//后台写日志
func (f *FileLogger) writeLogBackground() {
	for {
		if f.checkSize(f.fileObj) {
			newFile, err := f.cutFile(f.fileObj)
			if err != nil {
				fmt.Printf("文件切割失败，err:%v", err)
				return
			}
			f.fileObj = newFile
		}
		select {
		case logTmp := <-f.logChan:
			logInfo := fmt.Sprintf("[%s] [%s] [%s:%s:%d] %s\n", logTmp.timeStamp, getLogString(logTmp.level), logTmp.fileName, logTmp.funcName, logTmp.lineNumber, logTmp.msg)
			fmt.Fprint(f.fileObj, logInfo)
			if logTmp.level >= ERROR {
				if f.checkSize(f.errFileObj) {
					newFile, err := f.cutFile(f.errFileObj)
					if err != nil {
						fmt.Printf("文件切割失败，err:%v", err)
						return
					}
					f.errFileObj = newFile
				}
				fmt.Fprint(f.errFileObj, logInfo)
			}
		default:
			//取不到日志休息500毫秒
			time.Sleep(time.Millisecond * 500)
		}
	}
}

// log打印日志
func (f *FileLogger) log(level LogLevel, format string, a ...interface{}) {
	if f.enable(level) {
		nowTime := time.Now().Format("2006-01-02 15:04:05")
		fileName, funcName, lineNumber := getInfo(3)
		msg := fmt.Sprintf(format, a...)
		//logMsg对象
		logTemp := &logMsg{
			level:      level,
			lineNumber: lineNumber,
			msg:        msg,
			timeStamp:  nowTime,
			funcName:   funcName,
			fileName:   fileName,
		}
		select {
		case f.logChan <- logTemp:
		default:
			//如果通道满了，就丢掉日志
		}
	}
}
func (f *FileLogger) Debug(format string, a ...interface{}) {
	f.log(DEBUG, format, a...)
}
func (f *FileLogger) Trace(format string, a ...interface{}) {
	f.log(TRACE, format, a...)
}
func (f *FileLogger) Info(format string, a ...interface{}) {
	f.log(INFO, format, a...)
}
func (f *FileLogger) Warning(format string, a ...interface{}) {
	f.log(WARNING, format, a...)
}
func (f *FileLogger) Error(format string, a ...interface{}) {
	f.log(ERROR, format, a...)
}
func (f *FileLogger) Fatal(format string, a ...interface{}) {
	f.log(FATAL, format, a...)
}
