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

type ErrGRPCConnectionSuite struct {
	suite.Suite
	obLog *observer.ObservedLogs
}

func (suite *ErrGRPCConnectionSuite) SetupTest() {
	observedZapCore, observedLogs := observer.New(zap.WarnLevel)
	logger = zap.New(observedZapCore, zap.Fields(zap.String("system", "Mock system")))
	suite.obLog = observedLogs
}

func (suite *ErrGRPCConnectionSuite) TestNewErrGRPCConnection() {
	t := suite.T()
	assert.Equal(t, "*panicerrorhandler.ErrGRPCConnection", reflect.TypeOf(NewErrGRPCConnection(
		errors.New("got error"),
		struct{}{},
	)).String())
}

func (suite *ErrGRPCConnectionSuite) TestNewErrGRPCConnectionGetNameMethod() {
	t := suite.T()
	assert.Equal(t, ErrGrpcConnection, NewErrGRPCConnection(
		errors.New("got error"),
		struct{}{},
	).GetName())
}

func (suite *ErrGRPCConnectionSuite) TestNewErrGRPCConnectionGetErrorMethod() {
	t := suite.T()
	assert.Equal(t, errors.New("got error").Error(), NewErrGRPCConnection(
		errors.New("got error"),
		struct{}{},
	).GetError().Error())
}

func (suite *ErrGRPCConnectionSuite) TestNewErrGRPCConnectionImplementError() {
	t := suite.T()
	assert.Implements(t, (*error)(nil), NewErrGRPCConnection(
		errors.New("got error"),
		struct{}{},
	))
}

func (suite *ErrGRPCConnectionSuite) TestNewErrGRPCConnectionErrorMethod() {
	t := suite.T()
	assert.Equal(t, "[ERROR]: got error , response data: {}\n", NewErrGRPCConnection(
		errors.New("got error"),
		struct{}{},
	).Error())
}

func (suite *ErrGRPCConnectionSuite) TestNewErrGRPCConnectionReportMethod() {
	NewErrGRPCConnection(errors.New("got error"), struct{}{}).SetSystem("Mock system").Report("Error GRPC")
	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("Error GRPC", firstLog.Message)
	suite.Equal("Mock system", firstLog.Context[0].String)
	suite.Equal("got error", errors.Cause(firstLog.Context[1].Interface.(error)).Error())
}

func (suite *ErrGRPCConnectionSuite) TestNewErrGRPCConnectionGinReportMethod() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler("Mock Gin", "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		panic(NewErrGRPCConnection(errors.New("got error"), struct{}{}))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusServiceUnavailable, result.StatusCode)

	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("Mock system", firstLog.Context[0].String)
	suite.Equal("got error", errors.Cause(firstLog.Context[1].Interface.(error)).Error())
}

func (suite *ErrGRPCConnectionSuite) TestPanicGRPCErrorHandlerNewErrGRPCConnection() {
	t := suite.T()
	var errContent error
	func() {
		defer PanicGRPCErrorHandler(&errContent, "MockGRPCHandler", "Test error handler")
		panic(NewErrGRPCConnection(errors.New("database disconnect"), struct{}{}))
	}()
	assert.Error(t, errContent)
	if s, ok := status.FromError(errContent); ok {
		assert.Equal(t, "Unavailable", s.Code().String())
		assert.Equal(t, "Test error handler: database disconnect: response data: {}", s.Message())
		assert.Equal(t, "rpc error: code = Unavailable desc = Test error handler: database disconnect: response data: {}", s.Err().Error())
	}
}

func TestErrGRPCConnectionSuite(t *testing.T) {
	suite.Run(t, new(ErrGRPCConnectionSuite))
}
