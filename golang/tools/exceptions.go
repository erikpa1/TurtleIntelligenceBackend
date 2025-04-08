package tools

import (
	"runtime/debug"
	"sync"
	"turtle/lg"
)

func Recover(message string) {
	if r := recover(); r != nil {
		lg.LogStackTraceErr(message, r)
		debug.PrintStack()
	}
}

func RecoverWithMutUnlock(mutex sync.Mutex, message string) {
	defer mutex.Unlock()
	r := recover()
	if r != nil {
		lg.LogStackTraceErr(message, r)
		debug.PrintStack()
	}

}
