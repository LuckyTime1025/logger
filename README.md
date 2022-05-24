## Logger
#
### 异步日志
### 按照文件大小分割
### 日志级别
#
```go
log := logger.ConsoleNewLog("Warning")
	log := logger.FileNewLog("Debug", "filePath", "fileName", "fileSize")
	for {
		log.Debug("Debug")
		log.Trace("Trace")
		log.Info("Info")
		log.Warning("Warning")
		log.Error("Error")
		log.Fatal("Fatal")
	}
```
