package lib

import (
	"log"
	"os"
)

var (
	loggerInstance *log.Logger = nil
	logFile        *os.File
)

func LoggerInit(filename string) (error) {
	var err error = nil
	logFile, err = os.Create(filename)
	if err != nil {
		return err
	}

	loggerInstance = log.New(logFile, "[error]", log.Llongfile)
	return nil

}

func LoggerClose() {
	logFile.Close()
}

func Logger() (*log.Logger) {
	return loggerInstance
}
