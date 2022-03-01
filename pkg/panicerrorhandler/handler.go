package panicerrorhandler

import (
	"fmt"
	"go.uber.org/zap"
	"sync"
)

type IErrorReport interface {
	error
	SetSystem(system string) IErrorReport
	GetName() string
	GetError() error
	Report(prefix string)
}

func PanicErrorHandler(system, prefixMessage string) {
	if err := recover(); err != nil {
		if errReportInstance, okIErrorReport := err.(IErrorReport); okIErrorReport {
			errReportInstance.SetSystem(system).Report("")
			return
		}
		if e, ok := err.(error); ok {
			logger.Warn(prefixMessage, zap.Error(e))
		} else {
			logger.Sugar().Warn(prefixMessage, err)
		}
	}
}

func PanicWaitGroupErrorHandler(errContent *error, rmu *sync.RWMutex, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := recover(); err != nil {
		rmu.RLock()
		*errContent = fmt.Errorf("%v", err)
		rmu.RUnlock()
	}
}
