package zap

import (
	"github.com/justdomepaul/gin-storage/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewLogger(option config.Core) (*zap.Logger, error) {
	return map[string]func(...zap.Option) (*zap.Logger, error){
		"development": zap.NewDevelopment,
		"production":  zap.NewProduction,
		"customized":  NewCustomized(option),
	}[option.LoggerMode]()
}

var (
	Logger *zap.Logger
	sugar  *zap.SugaredLogger
)

func init() {
	systemName := os.Getenv("SYSTEM_NAME")
	if systemName == "" {
		systemName = "anonymous"
	}
	Logger, _ = newZap(zapcore.InfoLevel, systemName)
	sugar = Logger.Sugar()
}

func NewCustomized(option config.Core) func(...zap.Option) (*zap.Logger, error) {
	return func(...zap.Option) (*zap.Logger, error) {
		return newZap(zapcore.InfoLevel, option.SystemName)
	}
}

func newZap(level zapcore.Level, systemName string) (*zap.Logger, error) {
	return zap.New(newCore(level), zap.AddCaller(), zap.Fields(zap.String("system", systemName))), nil
}

func newCore(level zapcore.Level) zapcore.Core {
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		atomicLevel,
	)
}

func SugarInfo(args ...interface{}) {
	sugar.Info(args)
}

func SugarWarn(args ...interface{}) {
	sugar.Warn(args)
}

func SugarError(args ...interface{}) {
	sugar.Error(args)
}
