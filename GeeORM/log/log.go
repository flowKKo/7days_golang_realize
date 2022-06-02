package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	//  \033[color m content \033[ctl code  is used to control color and blink of the content
	// [info ] is blue and [error] is red
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

// support the level of log
const(
	InfoLevel = iota
	ErrorLevel
	Disabled
)

// SetLevel controls log level
func SetLevel(level int){
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers{
		logger.SetOutput(os.Stdout)
	}

	if ErrorLevel < level{
		errorLog.SetOutput(ioutil.Discard)
	}

	if InfoLevel < level{
		infoLog.SetOutput(ioutil.Discard)
	}
}