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

type ErrServerExecute struct {
	system string
	err    error
}

func (e *ErrServerExecute) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e ErrServerExecute) GetName() string {
	return ErrProcessServerExecute
}

func (e ErrServerExecute) GetError() error {
	return e.err
}

func (e ErrServerExecute) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error())
}

func (e ErrServerExecute) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()))
}

func (e ErrServerExecute) GinReport(c *gin.Context) {
	c.AbortWithError(http.StatusInternalServerError, e.err)
}

func (e ErrServerExecute) GRPCReport(errContent *error, prefixMessage string) {
	*errContent = status.Error(codes.FailedPrecondition, errors.Wrap(e.err, prefixMessage).Error())
}

func NewErrServerExecute(err error) *ErrServerExecute {
	return &ErrServerExecute{
		err: err,
	}
}
