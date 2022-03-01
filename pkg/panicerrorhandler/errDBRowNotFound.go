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

type ErrDBRowNotFound struct {
	system string
	err    error
}

func (e *ErrDBRowNotFound) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e ErrDBRowNotFound) GetName() string {
	return ErrDbRowNotFound
}

func (e ErrDBRowNotFound) GetError() error {
	return e.err
}

func (e ErrDBRowNotFound) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error())
}

func (e ErrDBRowNotFound) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()))
}

func (e ErrDBRowNotFound) GinReport(c *gin.Context) {
	c.AbortWithError(http.StatusNotFound, e.err)
}

func (e ErrDBRowNotFound) GRPCReport(errContent *error, prefixMessage string) {
	*errContent = status.Error(codes.NotFound, errors.Wrap(e.err, prefixMessage).Error())
}

func NewErrDBRowNotFound(err error) *ErrDBRowNotFound {
	return &ErrDBRowNotFound{
		err: err,
	}
}
