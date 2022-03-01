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

type ErrJWTExecuteSuite struct {
	suite.Suite
	obLog *observer.ObservedLogs
}

func (suite *ErrJWTExecuteSuite) SetupTest() {
	observedZapCore, observedLogs := observer.New(zap.WarnLevel)
	logger = zap.New(observedZapCore, zap.Fields(zap.String("system", "Mock system")))
	suite.obLog = observedLogs
}

func (suite *ErrJWTExecuteSuite) TestNewErrJWTExecute() {
	t := suite.T()
	assert.Equal(t, "*panicerrorhandler.ErrJWTExecute", reflect.TypeOf(NewErrJWTExecute(errors.New("got error"))).String())
}

func (suite *ErrJWTExecuteSuite) TestNewErrJWTExecuteGetNameMethod() {
	t := suite.T()
	assert.Equal(t, ErrJwtExecute, NewErrJWTExecute(errors.New("got error")).GetName())
}

func (suite *ErrJWTExecuteSuite) TestNewErrJWTExecuteGetErrorMethod() {
	t := suite.T()
	assert.Equal(t, errors.New("got error"), NewErrJWTExecute(errors.New("got error")).GetError())
}

func (suite *ErrJWTExecuteSuite) TestNewErrJWTExecuteImplementError() {
	t := suite.T()
	assert.Implements(t, (*error)(nil), NewErrJWTExecute(errors.New("got error")))
}

func (suite *ErrJWTExecuteSuite) TestNewErrJWTExecuteErrorMethod() {
	t := suite.T()
	assert.Equal(t, "[ERROR]: got error\n", NewErrJWTExecute(errors.New("got error")).Error())
}

func (suite *ErrJWTExecuteSuite) TestNewErrJWTExecuteReportMethod() {
	NewErrJWTExecute(errors.New("got error")).SetSystem("Mock system").Report("")
	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("", firstLog.Message)
	suite.Equal("Mock system", firstLog.Context[0].String)
	suite.Equal("got error", errors.Cause(firstLog.Context[1].Interface.(error)).Error())
}

func (suite *ErrJWTExecuteSuite) TestNewErrJWTExecuteGinReportMethod() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler("Mock Gin", "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		panic(NewErrJWTExecute(errors.New("got error")))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusForbidden, result.StatusCode)

	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("Mock system", firstLog.Context[0].String)
	suite.Equal("got error", errors.Cause(firstLog.Context[1].Interface.(error)).Error())
}

func (suite *ErrJWTExecuteSuite) TestPanicGRPCErrorHandlerNewErrJWTExecute() {
	t := suite.T()
	var errContent error
	func() {
		defer PanicGRPCErrorHandler(&errContent, "MockGRPCHandler", "Test error handler")
		panic(NewErrJWTExecute(errors.New("database disconnect")))
	}()
	assert.Error(t, errContent)
	if s, ok := status.FromError(errContent); ok {
		assert.Equal(t, "PermissionDenied", s.Code().String())
		assert.Equal(t, "Test error handler: database disconnect", s.Message())
		assert.Equal(t, "rpc error: code = PermissionDenied desc = Test error handler: database disconnect", s.Err().Error())
	}
}

func TestErrJWTExecuteSuite(t *testing.T) {
	suite.Run(t, new(ErrJWTExecuteSuite))
}
