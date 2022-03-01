package panicerrorhandler

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"net/http"
	"testing"
)

const (
	ERR_MOCK = "errMock"
)

type MockIError struct {
	system string
	name   string
	err    error
}

func (e *MockIError) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e MockIError) GetName() string {
	return e.name
}

func (e MockIError) GetError() error {
	return e.err
}

func (e MockIError) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error())
}

func (e MockIError) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()))
}

func (e MockIError) GinReport(c *gin.Context) {
	c.AbortWithError(http.StatusInternalServerError, e.err)
}

func NewMockIError(err error) IErrorReport {
	return &MockIError{
		name: ERR_MOCK,
		err:  err,
	}
}

type HandlerSuite struct {
	suite.Suite
	obLog *observer.ObservedLogs
}

func (suite *HandlerSuite) SetupTest() {
	observedZapCore, observedLogs := observer.New(zap.WarnLevel)
	logger = zap.New(observedZapCore)
	suite.obLog = observedLogs
}

func (suite *HandlerSuite) TestPanicErrorHandler() {
	func() {
		defer PanicErrorHandler("Mock Root System", "Mock test: ")
		panic(NewMockIError(errors.New("got error")))
	}()
	require.Equal(suite.T(), 1, suite.obLog.Len())
}

func (suite *HandlerSuite) TestPanicErrorHandlerPanicNormalError() {
	func() {
		defer PanicErrorHandler("Mock Root System", "Mock test: ")
		panic(errors.New("got error"))
	}()

	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("Mock test: ", firstLog.Message)
	suite.Equal("got error", errors.Cause(firstLog.Context[0].Interface.(error)).Error())
}

func (suite *HandlerSuite) TestPanicErrorHandlerPanicNormalStringError() {
	func() {
		defer PanicErrorHandler("Mock Root System", "Mock test: ")
		panic("got error")
	}()

	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("Mock test: got error", firstLog.Message)
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
