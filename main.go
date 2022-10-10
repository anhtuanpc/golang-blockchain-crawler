package main

import (
	"cbridgewrapper/logger"
	"cbridgewrapper/service"
	"runtime"
)

func main() {
	logger.Logger.Infof("Service initialize...")
	go service.ProcessTransaction()
	logger.Logger.Infof("Service started")
	runtime.Goexit()
}
