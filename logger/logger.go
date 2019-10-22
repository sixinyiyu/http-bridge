package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

func  newZapEncoderConfig()  zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "lineNum",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    CustomLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func init() {
	// 日志文件切割
	lumberjackLogger := lumberjack.Logger{
		Filename: "./logs/output.log",
		MaxSize: 128,
		MaxBackups: 3,
		MaxAge: 7,
		Compress: true,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(newZapEncoderConfig()), // zapcore.NewJSONEncoder()
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&lumberjackLogger)),
		zap.InfoLevel,
	)

	Logger = zap.New(core, zap.AddCaller(), zap.Development())
	Sugar = Logger.Sugar()
}