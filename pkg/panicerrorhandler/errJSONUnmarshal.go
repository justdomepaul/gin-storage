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

type ErrJSONUnmarshal struct {
	system string
	err    error
}

func (e *ErrJSONUnmarshal) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e ErrJSONUnmarshal) GetName() string {
	return ErrJsonUnmarshal
}

func (e ErrJSONUnmarshal) GetError() error {
	return e.err
}

func (e ErrJSONUnmarshal) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error())
}

func (e ErrJSONUnmarshal) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()))
}

func (e ErrJSONUnmarshal) GinReport(c *gin.Context) {
	c.AbortWithError(http.StatusBadRequest, e.err)
}

func (e ErrJSONUnmarshal) GRPCReport(errContent *error, prefixMessage string) {
	*errContent = status.Error(codes.InvalidArgument, errors.Wrap(e.err, prefixMessage).Error())
}

func NewErrJSONUnmarshal(err error) *ErrJSONUnmarshal {
	return &ErrJSONUnmarshal{
		err: err,
	}
}
