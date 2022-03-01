package zap

import (
	"github.com/justdomepaul/gin-storage/pkg/config"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"reflect"
	"testing"
)

type LoggerSuite struct {
	suite.Suite
}

func (suite *LoggerSuite) TestNewLogger() {
	option := config.Core{
		LoggerMode:  "customized",
		ReleaseMode: true,
		SystemName:  "system",
	}
	result, err := NewLogger(option)
	suite.NoError(err)
	suite.Equal("*zap.Logger", reflect.TypeOf(result).String())
}

func (suite *LoggerSuite) TestNewCustomized() {
	op := config.Core{
		LoggerMode:  "customized",
		ReleaseMode: true,
		SystemName:  "system",
	}
	logger, err := NewCustomized(op)()
	suite.NoError(err)
	suite.Equal("*zap.Logger", reflect.TypeOf(logger).String())
	suite.T().Log(logger)
}

func (suite *LoggerSuite) TestLogger() {
	suite.NotNil(Logger)
	suite.Equal("*zap.Logger", reflect.TypeOf(Logger).String())
}

func (suite *LoggerSuite) TestSugarInfo() {
	withSugar(suite.T(), zap.DebugLevel, nil, func(log *zap.SugaredLogger, logs *observer.ObservedLogs) {
		sugar = log
		SugarInfo(struct {
			Name string
		}{
			Name: "Max",
		})
		require.Equal(suite.T(), 1, logs.Len(), "Expected only one log entry to be written.")
		suite.T().Log(logs.AllUntimed()[0])
		suite.Equal("[{Max}]", logs.All()[0].Message)
	})
}

func (suite *LoggerSuite) TestSugarWarn() {
	withSugar(suite.T(), zap.DebugLevel, nil, func(log *zap.SugaredLogger, logs *observer.ObservedLogs) {
		sugar = log
		SugarWarn(struct {
			Name string
		}{
			Name: "Max",
		})
		require.Equal(suite.T(), 1, logs.Len(), "Expected only one log entry to be written.")
		suite.T().Log(logs.AllUntimed()[0])
		suite.Equal("[{Max}]", logs.All()[0].Message)
	})
}

func (suite *LoggerSuite) TestSugarError() {
	withSugar(suite.T(), zap.DebugLevel, nil, func(log *zap.SugaredLogger, logs *observer.ObservedLogs) {
		sugar = log
		SugarError(struct {
			Name string
		}{
			Name: "Max",
		})
		require.Equal(suite.T(), 1, logs.Len(), "Expected only one log entry to be written.")
		suite.T().Log(logs.AllUntimed()[0])
		suite.Equal("[{Max}]", logs.All()[0].Message)
	})
}

func TestLoggerSuite(t *testing.T) {
	suite.Run(t, new(LoggerSuite))
}

// Here specifically to introduce an easily-identifiable filename for testing
// stacktraces and caller skips.
func withLogger(t testing.TB, e zapcore.LevelEnabler, opts []zap.Option, f func(*zap.Logger, *observer.ObservedLogs)) {
	fac, logs := observer.New(e)
	log := zap.New(fac, opts...)
	f(log, logs)
}

func withSugar(t testing.TB, e zapcore.LevelEnabler, opts []zap.Option, f func(*zap.SugaredLogger, *observer.ObservedLogs)) {
	withLogger(t, e, opts, func(logger *zap.Logger, logs *observer.ObservedLogs) { f(logger.Sugar(), logs) })
}
