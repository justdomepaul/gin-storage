package panicerrorhandler

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type ErrDBExecute struct {
	system string
	err    error
}

func (e *ErrDBExecute) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e ErrDBExecute) GetName() string {
	return ErrDbExecute
}

func (e ErrDBExecute) GetError() error {
	return e.err
}

func (e ErrDBExecute) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error())
}

func (e ErrDBExecute) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()))
}

func (e ErrDBExecute) GinReport(c *gin.Context) {
	code := http.StatusConflict
	switch e.err {
	case sql.ErrNoRows, gocql.ErrNotFound, iterator.Done:
		code = http.StatusNotFound
	case driver.ErrBadConn, gocql.ErrNoConnections:
		code = http.StatusServiceUnavailable
	}
	c.AbortWithError(code, e.err)
}

func (e ErrDBExecute) GRPCReport(errContent *error, prefixMessage string) {
	code := codes.FailedPrecondition
	switch e.err {
	case ErrUpdateNoEffect:
		code = codes.FailedPrecondition
	case sql.ErrNoRows, gocql.ErrNotFound, iterator.Done, ErrNoRows:
		code = codes.NotFound
	case driver.ErrBadConn, gocql.ErrNoConnections:
		code = codes.Unavailable
	}
	*errContent = status.Error(code, errors.Wrap(e.err, prefixMessage).Error())
}

func NewErrDBExecute(err error) *ErrDBExecute {
	return &ErrDBExecute{
		err: err,
	}
}
