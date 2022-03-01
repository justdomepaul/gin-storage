package panicerrorhandler

import (
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc/status"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type ErrVariableSuite struct {
	suite.Suite
	obLog *observer.ObservedLogs
}

func (suite *ErrVariableSuite) SetupTest() {
	observedZapCore, observedLogs := observer.New(zap.WarnLevel)
	logger = zap.New(observedZapCore, zap.Fields(zap.String("system", "Mock system")))
	suite.obLog = observedLogs
}

func (suite *ErrVariableSuite) TestNewErrVariable() {
	t := suite.T()
	assert.Equal(t, "*panicerrorhandler.ErrVariable", reflect.TypeOf(NewErrVariable(errors.New("got error"))).String())
}

func (suite *ErrVariableSuite) TestNewErrVariableGetNameMethod() {
	t := suite.T()
	assert.Equal(t, ErrProcessVariable, NewErrVariable(errors.New("got error")).GetName())
}

func (suite *ErrVariableSuite) TestNewErrVariableGetErrorMethod() {
	t := suite.T()
	assert.Equal(t, errors.New("got error"), NewErrVariable(errors.New("got error")).GetError())
}

func (suite *ErrVariableSuite) TestNewErrVariableImplementError() {
	t := suite.T()
	assert.Implements(t, (*error)(nil), NewErrVariable(errors.New("got error")))
}

func (suite *ErrVariableSuite) TestNewErrVariableErrorMethod() {
	t := suite.T()
	assert.Equal(t, "[ERROR]: got error\n", NewErrVariable(errors.New("got error")).Error())
}

func (suite *ErrVariableSuite) TestNewErrVariableReportMethod() {
	NewErrVariable(errors.New("got error")).SetSystem("Mock system").Report("")
	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("", firstLog.Message)
	suite.Equal("Mock system", firstLog.Context[0].String)
	suite.Equal("got error", errors.Cause(firstLog.Context[1].Interface.(error)).Error())
}

func (suite *ErrVariableSuite) TestNewErrVariableGinReportMethod() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler("Mock Gin", "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		panic(NewErrVariable(errors.New("got error")))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusBadRequest, result.StatusCode)

	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("Mock system", firstLog.Context[0].String)
	suite.Equal("got error", errors.Cause(firstLog.Context[1].Interface.(error)).Error())
}

func (suite *ErrVariableSuite) TestPanicGRPCErrorHandlerNewErrVariable() {
	t := suite.T()
	var errContent error
	func() {
		defer PanicGRPCErrorHandler(&errContent, "MockGRPCHandler", "Test error handler")
		panic(NewErrVariable(errors.New("database disconnect")))
	}()
	assert.Error(t, errContent)
	if s, ok := status.FromError(errContent); ok {
		assert.Equal(t, "InvalidArgument", s.Code().String())
		assert.Equal(t, "Test error handler: database disconnect", s.Message())
		assert.Equal(t, "rpc error: code = InvalidArgument desc = Test error handler: database disconnect", s.Err().Error())
	}
}

func TestErrVariableSuite(t *testing.T) {
	suite.Run(t, new(ErrVariableSuite))
}
