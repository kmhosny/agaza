package logger

import (
	"io"
	"log"
	"os"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		if err = os.Mkdir("log", os.ModePerm); err != nil {
			panic("Cannot create log folder, kindly manually create it")
		}

	}
	logFile, err := os.OpenFile("log/log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		Trace = log.New(io.MultiWriter(os.Stdout),
			"TRACE: ",
			log.Ldate|log.Ltime|log.Lshortfile)

		Info = log.New(io.MultiWriter(os.Stdout),
			"INFO: ",
			log.Ldate|log.Ltime|log.Lshortfile)

		Warning = log.New(io.MultiWriter(os.Stdout),
			"WARNING: ",
			log.Ldate|log.Ltime|log.Lshortfile)

		Error = log.New(io.MultiWriter(os.Stdout),
			"ERROR: ",
			log.Ldate|log.Ltime|log.Lshortfile)
		Error.Println(err)
	} else {
		Trace = log.New(io.MultiWriter(os.Stdout, logFile),
			"TRACE: ",
			log.Ldate|log.Ltime|log.Lshortfile)

		Info = log.New(io.MultiWriter(os.Stdout, logFile),
			"INFO: ",
			log.Ldate|log.Ltime|log.Lshortfile)

		Warning = log.New(io.MultiWriter(os.Stdout, logFile),
			"WARNING: ",
			log.Ldate|log.Ltime|log.Lshortfile)

		Error = log.New(io.MultiWriter(os.Stdout, logFile),
			"ERROR: ",
			log.Ldate|log.Ltime|log.Lshortfile)
	}
}
