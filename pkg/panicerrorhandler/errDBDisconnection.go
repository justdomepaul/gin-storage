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

type ErrDBDisconnection struct {
	system string
	err    error
}

func (e *ErrDBDisconnection) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e ErrDBDisconnection) GetName() string {
	return ErrDbDisconnection
}

func (e ErrDBDisconnection) GetError() error {
	return e.err
}

func (e ErrDBDisconnection) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error())
}

func (e ErrDBDisconnection) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()))
}

func (e ErrDBDisconnection) GinReport(c *gin.Context) {
	c.AbortWithError(http.StatusServiceUnavailable, e.err)
}

func (e ErrDBDisconnection) GRPCReport(errContent *error, prefixMessage string) {
	*errContent = status.Error(codes.Unavailable, errors.Wrap(e.err, prefixMessage).Error())
}

func NewErrDBDisconnection(err error) *ErrDBDisconnection {
	return &ErrDBDisconnection{
		err: err,
	}
}
