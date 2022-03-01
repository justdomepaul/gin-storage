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

type ErrVariable struct {
	system string
	err    error
}

func (e *ErrVariable) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e ErrVariable) GetName() string {
	return ErrProcessVariable
}

func (e ErrVariable) GetError() error {
	return e.err
}

func (e ErrVariable) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error())
}

func (e ErrVariable) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()))
}

func (e ErrVariable) GinReport(c *gin.Context) {
	c.AbortWithError(http.StatusBadRequest, e.err)
}

func (e ErrVariable) GRPCReport(errContent *error, prefixMessage string) {
	*errContent = status.Error(codes.InvalidArgument, errors.Wrap(e.err, prefixMessage).Error())
}

func NewErrVariable(err error) *ErrVariable {
	return &ErrVariable{
		err: err,
	}
}
