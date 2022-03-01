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

type ErrJWTExecute struct {
	system string
	err    error
}

func (e *ErrJWTExecute) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e ErrJWTExecute) GetName() string {
	return ErrJwtExecute
}

func (e ErrJWTExecute) GetError() error {
	return e.err
}

func (e ErrJWTExecute) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error())
}

func (e ErrJWTExecute) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()))
}

func (e ErrJWTExecute) GinReport(c *gin.Context) {
	c.AbortWithError(http.StatusForbidden, e.err)
}

func (e ErrJWTExecute) GRPCReport(errContent *error, prefixMessage string) {
	*errContent = status.Error(codes.PermissionDenied, errors.Wrap(e.err, prefixMessage).Error())
}

func NewErrJWTExecute(err error) *ErrJWTExecute {
	return &ErrJWTExecute{
		err: err,
	}
}