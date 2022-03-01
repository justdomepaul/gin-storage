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

type ErrDBConnection struct {
	system string
	err    error
}

func (e *ErrDBConnection) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e ErrDBConnection) GetName() string {
	return ErrDbConnection
}

func (e ErrDBConnection) GetError() error {
	return e.err
}

func (e ErrDBConnection) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error())
}

func (e ErrDBConnection) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()))
}

func (e ErrDBConnection) GinReport(c *gin.Context) {
	c.AbortWithError(http.StatusServiceUnavailable, e.err)
}

func (e ErrDBConnection) GRPCReport(errContent *error, prefixMessage string) {
	*errContent = status.Error(codes.Unavailable, errors.Wrap(e.err, prefixMessage).Error())
}

func NewErrDBConnection(err error) *ErrDBConnection {
	return &ErrDBConnection{
		err: err,
	}
}
