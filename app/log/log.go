package log

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	Dir          string `json:"dir"`
	Prefix       string `json:"prefix"`
	Suffix       string `json:"suffix"`
	Level        string `json:"level"` // log level
	MaxAge       int    `json:"maxAge"`
	RotationTime int    `json:"rotationTime"`
	Development  bool   `json:"development"`
}

const (
	defaultLogDir         = "/var/log/"
	defaultLogPrefix      = "sense-api"
	defaultLogSuffix      = ".log"
	defaultMaxAge         = 7 // day
	defaultRotationTime   = 2 // hour
	defaultLogLevelString = "DEBUG"
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger(cfg *LogConfig, opts ...zap.Option) (*Logger, error) {
	logger := &Logger{}

	dir := cfg.Dir
	prefix := cfg.Prefix
	suffix := cfg.Suffix
	maxAge := cfg.MaxAge
	rotationTime := cfg.RotationTime
	level := cfg.Level
	development := cfg.Development

	if dir == "" {
		dir = defaultLogDir
	}
	if prefix == "" {
		prefix = defaultLogPrefix
	}
	if suffix == "" {
		suffix = defaultLogSuffix
	}
	if maxAge <= 0 {
		maxAge = defaultMaxAge
	}
	if rotationTime <= 0 {
		rotationTime = defaultRotationTime
	}
	if level == "" {
		level = defaultLogLevelString
	}

	writer, err := newRotateLogger(dir, prefix, suffix, maxAge, rotationTime, development)
	if err != nil {
		return nil, err
	}

	var zl zapcore.Level
	switch level {
	case "DEBUG":
		zl = zap.DebugLevel
	case "INFO":
		zl = zap.InfoLevel
	case "WARN":
		zl = zap.WarnLevel
	case "ERROR":
		zl = zap.ErrorLevel
	case "PANIC":
		zl = zap.PanicLevel
	case "FATAL":
		zl = zap.FatalLevel
	default:
		zl = zap.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writer,
		zl,
	)
	l := zap.New(core, zap.AddCaller())

	if development {
		l = l.WithOptions(zap.Development())
	}
	if 0 != len(opts) {
		l = l.WithOptions(opts...)
	}

	logger.SugaredLogger = l.Sugar()
	return logger, nil
}

func (l *Logger) L() *zap.Logger {
	if nil == l.SugaredLogger {
		return zap.L()
	}
	return l.SugaredLogger.Desugar()
}

func (l *Logger) S() *zap.SugaredLogger {
	if nil == l.SugaredLogger {
		return zap.S()
	}
	return l.SugaredLogger
}

func (l *Logger) Named(label string) *Logger {
	return &Logger{
		SugaredLogger: l.SugaredLogger.Named(label),
	}
}

func newRotateLogger(dir, prefix, suffix string, maxAge, rotationTime int, development bool) (zapcore.WriteSyncer, error) {
	if development {
		return zapcore.Lock(os.Stdout), nil
	}
	if dir[0:1] != "/" {
		dir = getCurrPath() + "/" + dir
	}
	if dir[len(dir)-1:] == "/" {
		dir = dir[:len(dir)-1]
	}
	// 校验路径是否存在。如果不存在，则创建之
	if ok, err := pathExist(dir); err != nil {
		return nil, err
	} else if !ok {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	baseLogName := prefix + suffix
	baseLogPath := path.Join(dir, baseLogName)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),                               // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(maxAge)*24*time.Hour),          // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Duration(rotationTime)*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		return nil, err
	}

	w := zapcore.AddSync(writer)
	return w, nil
}

func getCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	return ret
}

func pathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
