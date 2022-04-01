package log

import "go.uber.org/zap"

var (
	globalLogger *Logger
)

func InitGlobalLogger(cfg *LogConfig, opts ...zap.Option) error {
	log, err := NewLogger(cfg, opts...)
	globalLogger = log
	return err
}

func GlobalLogger() *Logger {
	return globalLogger
}

func FromZap(logger *zap.Logger) *Logger {
	return &Logger{
		SugaredLogger: logger.Sugar(),
	}
}

func DPanic(args ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.DPanicf(template, args...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.DPanicw(msg, keysAndValues...)
}

func Debug(args ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Debugw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Errorw(msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Fatalf(template, args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Fatalw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Infow(msg, keysAndValues...)
}

func Panic(args ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Panicf(template, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Panicw(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	if globalLogger == nil {
		InitGlobalLogger(&LogConfig{})
	}
	globalLogger.Warnw(msg, keysAndValues...)
}
