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

type ErrGRPCConnection struct {
	system   string
	err      error
	response interface{}
}

func (e *ErrGRPCConnection) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e ErrGRPCConnection) GetName() string {
	return ErrGrpcConnection
}

func (e ErrGRPCConnection) GetError() error {
	return e.err
}

func (e ErrGRPCConnection) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error(), ", response data:", e.response)
}

func (e ErrGRPCConnection) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()), zap.Error(fmt.Errorf("response data: %v \n", e.response)))
}

func (e ErrGRPCConnection) GinReport(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusServiceUnavailable, e.response)
}

func (e ErrGRPCConnection) GRPCReport(errContent *error, prefixMessage string) {
	*errContent = status.Error(codes.Unavailable, errors.Wrap(fmt.Errorf("response data: %v", e.response), errors.Wrap(e.err, prefixMessage).Error()).Error())
}

func NewErrGRPCConnection(err error, response interface{}) *ErrGRPCConnection {
	return &ErrGRPCConnection{
		err:      err,
		response: response,
	}
}
