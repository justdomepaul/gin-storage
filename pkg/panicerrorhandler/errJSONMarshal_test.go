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

type ErrJSONMarshalSuite struct {
	suite.Suite
	obLog *observer.ObservedLogs
}

func (suite *ErrJSONMarshalSuite) SetupTest() {
	observedZapCore, observedLogs := observer.New(zap.WarnLevel)
	logger = zap.New(observedZapCore, zap.Fields(zap.String("system", "Mock system")))
	suite.obLog = observedLogs
}

func (suite *ErrJSONMarshalSuite) TestNewErrJSONMarshal() {
	t := suite.T()
	assert.Equal(t, "*panicerrorhandler.ErrJSONMarshal", reflect.TypeOf(NewErrJSONMarshal(errors.New("got error"))).String())
}

func (suite *ErrJSONMarshalSuite) TestNewErrJSONMarshalGetNameMethod() {
	t := suite.T()
	assert.Equal(t, ErrJsonMarshal, NewErrJSONMarshal(errors.New("got error")).GetName())
}

func (suite *ErrJSONMarshalSuite) TestNewErrJSONMarshalGetErrorMethod() {
	t := suite.T()
	assert.Equal(t, errors.New("got error"), NewErrJSONMarshal(errors.New("got error")).GetError())
}

func (suite *ErrJSONMarshalSuite) TestNewErrJSONMarshalImplementError() {
	t := suite.T()
	assert.Implements(t, (*error)(nil), NewErrJSONMarshal(errors.New("got error")))
}

func (suite *ErrJSONMarshalSuite) TestNewErrJSONMarshalErrorMethod() {
	t := suite.T()
	assert.Equal(t, "[ERROR]: got error\n", NewErrJSONMarshal(errors.New("got error")).Error())
}

func (suite *ErrJSONMarshalSuite) TestNewErrJSONMarshalReportMethod() {
	NewErrJSONMarshal(errors.New("got error")).SetSystem("Mock system").Report("Error JSON")
	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("Error JSON", firstLog.Message)
	suite.Equal("Mock system", firstLog.Context[0].String)
	suite.Equal("got error", errors.Cause(firstLog.Context[1].Interface.(error)).Error())
}

func (suite *ErrJSONMarshalSuite) TestNewErrJSONMarshalGinReportMethod() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler("Mock Gin", "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		panic(NewErrJSONMarshal(errors.New("got error")))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusInternalServerError, result.StatusCode)

	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("Mock system", firstLog.Context[0].String)
	suite.Equal("got error", errors.Cause(firstLog.Context[1].Interface.(error)).Error())
}

func (suite *ErrJSONMarshalSuite) TestPanicGRPCErrorHandlerNewErrJSONMarshal() {
	t := suite.T()
	var errContent error
	func() {
		defer PanicGRPCErrorHandler(&errContent, "MockGRPCHandler", "Test error handler")
		panic(NewErrJSONMarshal(errors.New("json marshal disconnect")))
	}()
	assert.Error(t, errContent)
	if s, ok := status.FromError(errContent); ok {
		assert.Equal(t, "InvalidArgument", s.Code().String())
		assert.Equal(t, "Test error handler: json marshal disconnect", s.Message())
		assert.Equal(t, "rpc error: code = InvalidArgument desc = Test error handler: json marshal disconnect", s.Err().Error())
	}
}

func TestErrJSONMarshalSuite(t *testing.T) {
	suite.Run(t, new(ErrJSONMarshalSuite))
}
