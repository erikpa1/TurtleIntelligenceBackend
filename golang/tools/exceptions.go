package tools

import (
	"github.com/erikpa1/turtle/lg"
	"runtime/debug"
	"sync"
)

func Recover(message string, onErrr ...func(any)) {
	if r := recover(); r != nil {
		lg.LogStackTraceErr(message, r)
		lg.LogW(r)
		debug.PrintStack()

		for _, receiver := range onErrr {
			receiver(r)
		}

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
