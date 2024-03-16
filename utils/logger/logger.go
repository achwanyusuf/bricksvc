package logger

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

const (
	LevelDebug string = "debug"
	LevelInfo  string = "info"
	LevelWarn  string = "warn"
	LevelError string = "error"
	LevelFatal string = "fatal"
	LevelPanic string = "panic"
)

var (
	defaultConf Config = Config{
		Level: LevelDebug,
		CustomFields: map[string]interface{}{
			"namespace":  "bricksvc",
			"version":    "v1.0.0",
			"build_time": time.Now(),
			"pid":        os.Getpid(),
		},
	}
)

type Config struct {
	IsFile       bool
	FilePath     string
	Level        string
	CustomFields map[string]interface{}
}

type Logger struct {
	log      *logrus.Logger
	logEntry *logrus.Entry
}

type LoggerInterface interface {
	DebugWithContext(ctx *fiber.Ctx, v ...interface{})
	InfoWithContext(ctx *fiber.Ctx, v ...interface{})
	WarnWithContext(ctx *fiber.Ctx, v ...interface{})
	ErrorWithContext(ctx *fiber.Ctx, v ...interface{})
	FatalWithContext(ctx *fiber.Ctx, v ...interface{})
	PanicWithContext(ctx *fiber.Ctx, v ...interface{})
}

var Log *Logger

func New(c *Config) LoggerInterface {
	if c == nil {
		c = &defaultConf
	}
	log := logrus.New()
	if c.IsFile {
		f, err := os.OpenFile(c.FilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o755)
		if err != nil {
			panic(err)
		}
		log.SetOutput(f)
	} else {
		log.SetOutput(os.Stdout)
	}
	logEntry := log.WithFields(c.CustomFields)
	log.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	setLevel(log, c.Level)
	Log = &Logger{
		log:      log,
		logEntry: logEntry,
	}
	return Log
}

func setLevel(log *logrus.Logger, level string) {
	switch level {
	case LevelDebug:
		log.SetLevel(logrus.DebugLevel)
	case LevelInfo:
		log.SetLevel(logrus.InfoLevel)
	case LevelWarn:
		log.SetLevel(logrus.WarnLevel)
	case LevelError:
		log.SetLevel(logrus.ErrorLevel)
	case LevelFatal:
		log.SetLevel(logrus.FatalLevel)
	case LevelPanic:
		log.SetLevel(logrus.PanicLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
}

func (d *Logger) buildContextField(ctx *fiber.Ctx) *logrus.Entry {
	header := ctx.GetReqHeaders()
	d.logEntry = d.logEntry.WithFields(map[string]interface{}{
		"request_id":  header["X-Custom-Header"],
		"method":      ctx.Method(),
		"scheme":      header["x-request-scheme"],
		"client_ip":   header["x-forwarded-for"],
		"path":        ctx.Path(),
		"user_agent":  string(ctx.Request().Header.UserAgent()),
		"remote_addr": ctx.Context().RemoteAddr(),
		"remote_ip":   ctx.Context().RemoteIP(),
	})
	return d.logEntry
}

func (l *Logger) DebugWithContext(ctx *fiber.Ctx, v ...interface{}) {
	l.buildContextField(ctx).Debug(v...)
}

func (l *Logger) InfoWithContext(ctx *fiber.Ctx, v ...interface{}) {
	l.buildContextField(ctx).Info(v...)
}

func (l *Logger) WarnWithContext(ctx *fiber.Ctx, v ...interface{}) {
	l.buildContextField(ctx).Warn(v...)
}

func (l *Logger) ErrorWithContext(ctx *fiber.Ctx, v ...interface{}) {
	l.buildContextField(ctx).Error(v...)
}

func (l *Logger) FatalWithContext(ctx *fiber.Ctx, v ...interface{}) {
	l.buildContextField(ctx).Fatal(v...)
}

func (l *Logger) PanicWithContext(ctx *fiber.Ctx, v ...interface{}) {
	l.buildContextField(ctx).Panic(v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.logEntry.Debug(v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.logEntry.Info(v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.logEntry.Warn(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.logEntry.Error(v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.logEntry.Fatal(v...)
}

func (l *Logger) Panic(v ...interface{}) {
	l.logEntry.Panic(v...)
}
