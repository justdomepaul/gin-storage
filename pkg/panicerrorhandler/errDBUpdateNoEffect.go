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

type ErrDBUpdateNoEffect struct {
	system string
	err    error
}

func (e *ErrDBUpdateNoEffect) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e ErrDBUpdateNoEffect) GetName() string {
	return ErrDbUpdateNoEffect
}

func (e ErrDBUpdateNoEffect) GetError() error {
	return e.err
}

func (e ErrDBUpdateNoEffect) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error())
}

func (e ErrDBUpdateNoEffect) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()))
}

func (e ErrDBUpdateNoEffect) GinReport(c *gin.Context) {
	c.AbortWithError(http.StatusNotFound, e.err)
}

func (e ErrDBUpdateNoEffect) GRPCReport(errContent *error, prefixMessage string) {
	*errContent = status.Error(codes.NotFound, errors.Wrap(e.err, prefixMessage).Error())
}

func NewErrDBUpdateNoEffect(err error) *ErrDBUpdateNoEffect {
	return &ErrDBUpdateNoEffect{
		err: err,
	}
}
