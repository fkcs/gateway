package logger

import (
	"github.com/fkcs/gateway/internal/interfaces/dto"
	errord "github.com/fkcs/gateway/internal/utils/error"
	"github.com/fkcs/gateway/internal/utils/types"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

var logger *ZapLogger

type ZapLogger struct {
	*zap.SugaredLogger
	logLevel zap.AtomicLevel
}

func LogInit(logLevel, filename string, jsonFormat, logInConsole bool) {
	var encoder zapcore.Encoder
	var level zapcore.Level

	hook := lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    128,      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,       // 日志文件最多保存多少个备份
		MaxAge:     7,        // 文件最多保存多少天
		Compress:   true,     // 是否压缩
	}
	/*
		var syncer zapcore.WriteSyncer
		if !logInConsole && "" != filename {
			syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook))
		} else {
			syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
		}*/

	formatEncodeTime := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     formatEncodeTime,
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	if jsonFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 设置日志级别,debug可以打印出info,debug,warn；info级别可以打印warn，info；warn只能打印warn ,debug->info->warn->error
	if level.UnmarshalText([]byte(logLevel)) != nil {
		level = zapcore.DebugLevel
	}
	atomicLevel := zap.NewAtomicLevelAt(level)
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)),
		atomicLevel,
	)
	logger = &ZapLogger{
		SugaredLogger: zap.New(core).WithOptions(zap.AddCaller()).Sugar(),
		logLevel:      atomicLevel,
	}
}

func Logger() *ZapLogger {
	if logger == nil {
		logger.Errorf("logger is nil")
		return nil
	}
	return logger
}

func SetLogLevel(level string) dto.ErrorCode {
	switch strings.ToLower(level) {
	case "debug":
		logger.logLevel.SetLevel(zap.DebugLevel)
	case "info":
		logger.logLevel.SetLevel(zap.InfoLevel)
	case "warn":
		logger.logLevel.SetLevel(zap.WarnLevel)
	case "error":
		logger.logLevel.SetLevel(zap.ErrorLevel)
	case "panic":
		logger.logLevel.SetLevel(zap.PanicLevel)
	case "fatal":
		logger.logLevel.SetLevel(zap.FatalLevel)
	default:
		logger.Error("Unknown log level: %s", level)
		return errord.MakeBadRequest(types.InvalidLogLevel)
	}
	return errord.MakeOkRsp(level)
}

func LoadLogLevel() dto.ErrorCode {
	level := logger.logLevel.Level().String()
	Logger().Infof("level is %v", level)
	return errord.MakeOkRsp(level)
}
