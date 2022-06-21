package common

import (
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"runtime/debug"
)

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}
	if err := recover(); err != nil {
		logger.Logger().Errorf("%v", err)
		logger.Logger().Errorf("Stacktrace:\n%s\n", debug.Stack())
	}
}

func GoSafe(fn func()) {
	go RunSafe(fn)
}

func BatchGoSafe(fns ...func()) {
	for _, fn := range fns {
		go RunSafe(fn)
	}
}

func RunSafe(fn func()) {
	defer Recover()

	fn()
}
