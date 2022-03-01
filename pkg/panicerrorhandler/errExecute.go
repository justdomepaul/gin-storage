package panicerrorhandler

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type ErrExecute struct {
	system string
	err    error
}

func (e *ErrExecute) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e ErrExecute) GetName() string {
	return ErrProcessExecute
}

func (e ErrExecute) GetError() error {
	return e.err
}

func (e ErrExecute) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error())
}

func (e ErrExecute) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()))
}

func (e ErrExecute) GinReport(c *gin.Context) {
	c.AbortWithError(http.StatusUnprocessableEntity, e.err)
}

func (e ErrExecute) GRPCReport(errContent *error, prefixMessage string) {
	*errContent = status.Error(codes.FailedPrecondition, errors.Wrap(e.err, prefixMessage).Error())
}

func NewErrExecute(err error) *ErrExecute {
	return &ErrExecute{
		err: err,
	}
}
