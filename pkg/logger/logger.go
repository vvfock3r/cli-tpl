package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	AtomicLevel, _    = NewAtomicLevel("info")
	defaultEncoder, _ = NewEndocer("console")
	defaultSyncer, _  = NewOutput([]string{"stdout"})
)

var lastOutputFiles []*os.File
var currOutputFiles []*os.File

// defaultLogger 设置默认的Logger
var defaultLogger = NewLogger(AtomicLevel, defaultEncoder, defaultSyncer)

// SetDefaultLogger 设置默认的Logger
func SetDefaultLogger(logger *zap.Logger) {
	_ = defaultLogger.Sync()
	defaultLogger = logger
	for _, file := range lastOutputFiles {
		_ = file.Close()
	}
}

// NewLogger 构造函数
func NewLogger(l zap.AtomicLevel, enc zapcore.Encoder, ws zapcore.WriteSyncer) *zap.Logger {
	// 创建 Core
	core := zapcore.NewCore(enc, ws, l)

	// 创建 *Logger
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.DPanicLevel),
	)

	return logger
}

// NewAtomicLevel 日志级别
func NewAtomicLevel(level string) (zap.AtomicLevel, error) {
	l, err := zapcore.ParseLevel(level)
	if err != nil {
		unrecognized := "unrecognized level: " + level
		supported := ", supported values: debug,info,warn,error,dpanic,panic,fatal"
		return zap.NewAtomicLevelAt(zapcore.InvalidLevel), fmt.Errorf(unrecognized + supported)
	}
	return zap.NewAtomicLevelAt(l), nil
}

// NewEndocer 配置选项
func NewEndocer(format string) (zapcore.Encoder, error) {
	// 实例化 encoderConfig
	zap.NewProductionEncoderConfig()
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller",
		FunctionKey:      zapcore.OmitKey,
		MessageKey:       "message",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.LowercaseLevelEncoder,
		EncodeTime:       zapcore.TimeEncoderOfLayout(time.DateTime),
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeName:       func(s string, encoder zapcore.PrimitiveArrayEncoder) { encoder.AppendString(s) },
		ConsoleSeparator: "",
	}

	// 实例化 encoder
	switch format {
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig), nil
	case "console":
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		return zapcore.NewConsoleEncoder(encoderConfig), nil
	default:
		unrecognized := "unrecognized format: " + format
		supported := ", supported values: json,console"
		return zapcore.NewConsoleEncoder(encoderConfig), fmt.Errorf(unrecognized + supported)
	}
}

// NewOutput 输出位置
func NewOutput(outs []string) (zapcore.WriteSyncer, error) {
	var syncerSlice []zapcore.WriteSyncer
	var newOutputFiles []*os.File

	for _, out := range outs {
		switch out {
		case "stdout":
			syncerSlice = append(syncerSlice, zapcore.AddSync(os.Stdout))
		case "stderr":
			syncerSlice = append(syncerSlice, zapcore.AddSync(os.Stderr))
		default:
			file, err := os.OpenFile(out, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return nil, err
			}
			syncerSlice = append(syncerSlice, zapcore.Lock(zapcore.AddSync(file)))
			newOutputFiles = append(newOutputFiles, file)
		}
	}
	lastOutputFiles = currOutputFiles
	currOutputFiles = newOutputFiles
	return zapcore.NewMultiWriteSyncer(syncerSlice...), nil
}

func Debug(msg string, fields ...zap.Field) {
	defaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	defaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	defaultLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	defaultLogger.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	defaultLogger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	defaultLogger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	defaultLogger.Fatal(msg, fields...)
}
