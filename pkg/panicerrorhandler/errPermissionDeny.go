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

type ErrPermissionDeny struct {
	system string
	err    error
}

func (e *ErrPermissionDeny) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e ErrPermissionDeny) GetName() string {
	return ErrProcessPermissionDeny
}

func (e ErrPermissionDeny) GetError() error {
	return e.err
}

func (e ErrPermissionDeny) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error())
}

func (e ErrPermissionDeny) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()))
}

func (e ErrPermissionDeny) GinReport(c *gin.Context) {
	c.AbortWithError(http.StatusForbidden, e.err)
}

func (e ErrPermissionDeny) GRPCReport(errContent *error, prefixMessage string) {
	*errContent = status.Error(codes.PermissionDenied, errors.Wrap(e.err, prefixMessage).Error())
}

func NewErrPermissionDeny(err error) *ErrPermissionDeny {
	return &ErrPermissionDeny{
		err: err,
	}
}
