package myLog

import (
	"fmt"
	"io"
	"log"
	"os"
)

var Trace, Info, Warning, Fatal *log.Logger

func logInit(traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {
	Trace = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle, "Warning: ", log.Ldate|log.Ltime|log.Lshortfile)
	Fatal = log.New(errorHandle, "Fatal: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func init() {

	logFile, logErr := os.OpenFile("test.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("open file err", logErr)
		os.Exit(1)
	}
	logInit(logFile, logFile, logFile, logFile)
}
