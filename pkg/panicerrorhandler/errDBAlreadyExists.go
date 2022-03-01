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

type ErrDBAlreadyExists struct {
	system string
	err    error
}

func (e *ErrDBAlreadyExists) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e ErrDBAlreadyExists) GetName() string {
	return ErrDbAlreadyExists
}

func (e ErrDBAlreadyExists) GetError() error {
	return e.err
}

func (e ErrDBAlreadyExists) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error())
}

func (e ErrDBAlreadyExists) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()))
}

func (e ErrDBAlreadyExists) GinReport(c *gin.Context) {
	c.AbortWithError(http.StatusConflict, e.err)
}

func (e ErrDBAlreadyExists) GRPCReport(errContent *error, prefixMessage string) {
	*errContent = status.Error(codes.AlreadyExists, errors.Wrap(e.err, prefixMessage).Error())
}

func NewErrDBAlreadyExists(err error) *ErrDBAlreadyExists {
	return &ErrDBAlreadyExists{
		err: err,
	}
}
