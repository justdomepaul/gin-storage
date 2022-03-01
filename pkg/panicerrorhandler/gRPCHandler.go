package panicerrorhandler

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
)

type IGRPCErrorReport interface {
	IErrorReport
	GRPCReport(errContent *error, prefixMessage string)
}

func PanicGRPCErrorHandler(errContent *error, system, prefixMessage string) {
	if err := recover(); err != nil {
		if result, ok := err.(error); ok {
			if errReportInstance, okIErrorReport := result.(IGRPCErrorReport); okIErrorReport {
				errReportInstance.SetSystem(system).Report(prefixMessage)
				errReportInstance.GRPCReport(errContent, prefixMessage)
				return
			}
			fromStatusError(errContent, result, prefixMessage)
			return
		}
		*errContent = status.Error(codes.Unknown, errors.Wrap(fmt.Errorf("%v", err), prefixMessage).Error())
	}
}

func fromStatusError(errContent *error, input error, prefixMessage string) {
	if statusErr, ok := status.FromError(input); ok {
		*errContent = status.Error(statusErr.Code(), errors.Wrap(errors.New(statusErr.Message()), prefixMessage).Error())
	} else {
		*errContent = status.Error(codes.Unknown, errors.Wrap(input.(error), prefixMessage).Error())
	}
}
