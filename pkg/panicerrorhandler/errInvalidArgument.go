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

type ErrInvalidArgument struct {
	system string
	err    error
}

func (e *ErrInvalidArgument) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e ErrInvalidArgument) GetName() string {
	return ErrProcessInvalidArgument
}

func (e ErrInvalidArgument) GetError() error {
	return e.err
}

func (e ErrInvalidArgument) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error())
}

func (e ErrInvalidArgument) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()))
}

func (e ErrInvalidArgument) GinReport(c *gin.Context) {
	c.AbortWithError(http.StatusBadRequest, e.err)
}

func (e ErrInvalidArgument) GRPCReport(errContent *error, prefixMessage string) {
	*errContent = status.Error(codes.InvalidArgument, errors.Wrap(e.err, prefixMessage).Error())
}

func NewErrInvalidArgument(err error) *ErrInvalidArgument {
	return &ErrInvalidArgument{
		err: err,
	}
}
