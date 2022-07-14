/* Package exposing method to retrieving an aggregated logging to log on all required levels */

package logging

import (
	"log"
	"os"
	"errors"
	"io"
)

type Logger struct {
	infoLogger *log.Logger
	warnLogger *log.Logger
	errorLogger *log.Logger
}

func (l *Logger) Info(s string, a ...interface{}) {
	l.infoLogger.Printf(s, a...)
}

func (l *Logger) Warn(s string, a ...interface{}) {
	l.warnLogger.Printf(s, a...)
}

func (l *Logger) Error(s string, a ...interface{}) {
	l.errorLogger.Panicf(s, a...)
}

// Public method to instantiate base loggers, return aggregated logger object
func GetLogger(logPath string) (Logger, *os.File) {
	var file *os.File
	if exists(logPath){
		file, _ = os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		file, _ = os.Create(logPath)
	}
	mw := io.MultiWriter(os.Stdout, file) 
	// common logging flags
	flags := log.LstdFlags | log.Lshortfile
	infoLogger := log.New(mw, "INFO ", flags)
	warnLogger := log.New(mw, "WARN ", flags)
	errorLogger := log.New(mw, "ERROR ", flags)

	return Logger{infoLogger: infoLogger, warnLogger: warnLogger, errorLogger: errorLogger}, file
}

// helper method to check if given log file already exists
func exists(name string) bool {
    _, err := os.Stat(name)
    if err == nil {
        return true
    }
    if errors.Is(err, os.ErrNotExist) {
        return false
    }
    return false
}
