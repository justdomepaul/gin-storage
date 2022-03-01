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

type ErrNotFound struct {
	system string
	err    error
}

func (e *ErrNotFound) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e ErrNotFound) GetName() string {
	return ErrDataNotFound
}

func (e ErrNotFound) GetError() error {
	return e.err
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error())
}

func (e ErrNotFound) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()))
}

func (e ErrNotFound) GinReport(c *gin.Context) {
	c.AbortWithError(http.StatusNotFound, e.err)
}

func (e ErrNotFound) GRPCReport(errContent *error, prefixMessage string) {
	*errContent = status.Error(codes.NotFound, errors.Wrap(e.err, prefixMessage).Error())
}

func NewErrNotFound(err error) *ErrNotFound {
	return &ErrNotFound{
		err: err,
	}
}
