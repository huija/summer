package summer

import (
	"errors"
	"github.com/huija/summer/conf"
	"github.com/huija/summer/container/pipeline"
	"github.com/huija/summer/logs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type zapSetup struct {
	enc   zapcore.Encoder
	level zapcore.Level
	out   zapcore.WriteSyncer
}

var setups []*zapSetup       // core's setup
var ws []zapcore.WriteSyncer // core's root writers
var cores []zapcore.Core     // cores

// ZapLogger fast logger
var ZapLogger *zap.Logger

// https://github.com/uber-go/zap
// https://github.com/natefinch/lumberjack
// 使用教程: https://zhuanlan.zhihu.com/p/88856378
// zap源码分析: https://mp.weixin.qq.com/s/i0bMh_gLLrdnhAEWlF-xDw
func logger() (err error) {
	if conf.Config.Logs == nil {
		return nil
	}

	for _, c := range conf.Config.Logs {
		// skip DPanic level
		if c.Level >= logs.PanicLevel {
			c.Level++
		}
		// transform into cores
		switch c.Type {
		case logs.Console:
			setupConsoleLog(c)
		case logs.File:
			setupFileLog(c)
		default:
			return errors.New("config: log type invalid")
		}
	}

	// setup cores
	for _, s := range setups {
		cores = append(cores, zapcore.NewCore(s.enc, s.out, s.level))
	}

	// final logger
	ZapLogger = zap.New(zapcore.NewTee(cores...), zap.AddCaller())

	// setup global variables
	logs.SugaredLogger = ZapLogger.Sugar()
	logs.Writer = zapcore.NewMultiWriteSyncer(ws...)

	logs.SugaredLogger.Debug("loggers init successfully!")
	return
}

func setupConsoleLog(c *logs.Log) {
	ws = append(ws, zapcore.AddSync(os.Stdout))
	// encoderConfig
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05.000")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	setups = append(setups, &zapSetup{
		enc:   zapcore.NewConsoleEncoder(encoderConfig),
		level: zapcore.Level(c.Level),
		out:   ws[len(ws)-1],
	})
}

func setupFileLog(c *logs.Log) {
	sliceLogger := &lumberjack.Logger{
		Filename:   c.File.Path,
		MaxSize:    c.File.Maxsize,
		MaxBackups: c.File.MaxBackups,
		MaxAge:     c.File.MaxAge,
		Compress:   c.File.Compress,
		LocalTime:  c.File.LocalZone,
	}

	RegisterClose(Log, func() error {
		logs.SugaredLogger.Info("slice logger close...")
		return sliceLogger.Close()
	}, pipeline.ClosePipeLevel(pipeline.Debugger))

	ws = append(ws, zapcore.AddSync(zapcore.AddSync(sliceLogger)))
	// encoderConfig
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	setups = append(setups, &zapSetup{
		enc:   zapcore.NewJSONEncoder(encoderConfig),
		level: zapcore.Level(c.Level),
		out:   ws[len(ws)-1],
	})
}
