package panicerrorhandler

import (
	"github.com/cockroachdb/errors"
	"github.com/gogo/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"testing"
)

type GRPCHandlerSuite struct {
	suite.Suite
}

func (suite *GRPCHandlerSuite) TestPanicGRPCErrorHandlerNormalStringError() {
	t := suite.T()
	var errContent error
	func() {
		defer PanicGRPCErrorHandler(&errContent, "MockGRPCHandler", "Test error handler")
		panic("got error")
	}()
	assert.Error(t, errContent)
	if s, ok := status.FromError(errContent); ok {
		assert.Equal(t, "Unknown", s.Code().String())
		assert.Equal(t, "Test error handler: got error", s.Message())
		assert.Equal(t, "rpc error: code = Unknown desc = Test error handler: got error", s.Err().Error())
	}
}

func (suite *GRPCHandlerSuite) TestPanicGRPCErrorHandlerUnknown() {
	t := suite.T()
	var errContent error
	func() {
		defer PanicGRPCErrorHandler(&errContent, "MockGRPCHandler", "Test error handler")
		panic(errors.New("got error"))
	}()
	assert.Error(t, errContent)
	if s, ok := status.FromError(errContent); ok {
		assert.Equal(t, "Unknown", s.Code().String())
		assert.Equal(t, "Test error handler: got error", s.Message())
		assert.Equal(t, "rpc error: code = Unknown desc = Test error handler: got error", s.Err().Error())
	}
}

func (suite *GRPCHandlerSuite) TestPanicGRPCErrorHandlerStatusError() {
	t := suite.T()
	var errContent error
	func() {
		defer PanicGRPCErrorHandler(&errContent, "MockGRPCHandler", "Test error handler")
		panic(status.Errorf(codes.AlreadyExists, "got error"))
	}()
	assert.Error(t, errContent)
	if s, ok := status.FromError(errContent); ok {
		assert.Equal(t, "AlreadyExists", s.Code().String())
		assert.Equal(t, "Test error handler: got error", s.Message())
		assert.Equal(t, "rpc error: code = AlreadyExists desc = Test error handler: got error", s.Err().Error())
	}
}

func TestGRPCHandlerSuite(t *testing.T) {
	suite.Run(t, new(GRPCHandlerSuite))
}
