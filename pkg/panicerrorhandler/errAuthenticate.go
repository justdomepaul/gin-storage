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

type ErrAuthenticate struct {
	system string
	err    error
}

func (e *ErrAuthenticate) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e ErrAuthenticate) GetName() string {
	return ErrProcessAuthenticate
}

func (e ErrAuthenticate) GetError() error {
	return e.err
}

func (e ErrAuthenticate) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error())
}

func (e ErrAuthenticate) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()))
}

func (e ErrAuthenticate) GinReport(c *gin.Context) {
	c.AbortWithError(http.StatusUnauthorized, e.err)
}

func (e ErrAuthenticate) GRPCReport(errContent *error, prefixMessage string) {
	*errContent = status.Error(codes.Unauthenticated, errors.Wrap(e.err, prefixMessage).Error())
}

func NewErrAuthenticate(err error) *ErrAuthenticate {
	return &ErrAuthenticate{
		err: err,
	}
}
