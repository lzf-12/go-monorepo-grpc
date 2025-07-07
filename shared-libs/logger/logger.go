package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// interface for logging operations
type Logger interface {
	Debug(msg string, fields ...map[string]interface{})
	Info(msg string, fields ...map[string]interface{})
	Error(msg string, fields ...map[string]interface{})
	Warn(msg string, fields ...map[string]interface{})
	Fatal(msg string, fields ...map[string]interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	WithFields(fields map[string]interface{}) Logger
	WithContext(ctx context.Context) Logger
	WithComponent(component string) Logger
}

// wraps zerolog.Logger to implement our Logger interface
type ZerologWrapper struct {
	logger zerolog.Logger
}

type Config struct {
	Level       string // debug, info, warn, error, fatal
	Format      string // json, console
	Output      io.Writer
	TimeFormat  string
	Caller      bool
	ServiceName string
	Version     string
}

func DefaultConfig() *Config {
	return &Config{
		Level:       "info",
		Format:      "json",
		Output:      os.Stdout,
		TimeFormat:  time.RFC3339,
		Caller:      true,
		ServiceName: "",
		Version:     "",
	}
}

// creates a new logger with the given configuration
func New(config *Config) Logger {
	if config == nil {
		config = DefaultConfig()
	}

	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	if config.TimeFormat != "" {
		zerolog.TimeFieldFormat = config.TimeFormat
	}

	var zlog zerolog.Logger
	if config.Format == "console" {
		zlog = zerolog.New(zerolog.ConsoleWriter{
			Out:        config.Output,
			TimeFormat: "15:04:05",
			NoColor:    false,
		})
	} else {
		zlog = zerolog.New(config.Output)
	}

	zlog = zlog.With().Timestamp().Logger()

	if config.Caller {
		zlog = zlog.With().Caller().Logger()
	}

	if config.ServiceName != "" {
		zlog = zlog.With().Str("service", config.ServiceName).Logger()
	}

	if config.Version != "" {
		zlog = zlog.With().Str("version", config.Version).Logger()
	}

	return &ZerologWrapper{logger: zlog}
}

// NewFromZerolog creates a wrapper from an existing zerolog.Logger
func NewFromZerolog(zlog zerolog.Logger) Logger {
	return &ZerologWrapper{logger: zlog}
}

// Debug logs a debug message
func (z *ZerologWrapper) Debug(msg string, fields ...map[string]interface{}) {
	event := z.logger.Debug()
	z.addFields(event, fields...)
	event.Msg(msg)
}

// Info logs an info message
func (z *ZerologWrapper) Info(msg string, fields ...map[string]interface{}) {
	event := z.logger.Info()
	z.addFields(event, fields...)
	event.Msg(msg)
}

// Error logs an error message
func (z *ZerologWrapper) Error(msg string, fields ...map[string]interface{}) {
	event := z.logger.Error()
	z.addFields(event, fields...)
	event.Msg(msg)
}

// Warn logs a warning message
func (z *ZerologWrapper) Warn(msg string, fields ...map[string]interface{}) {
	event := z.logger.Warn()
	z.addFields(event, fields...)
	event.Msg(msg)
}

// Fatal logs a fatal message and exits
func (z *ZerologWrapper) Fatal(msg string, fields ...map[string]interface{}) {
	event := z.logger.Fatal()
	z.addFields(event, fields...)
	event.Msg(msg)
}

// Debugf logs a debug message with format
func (z *ZerologWrapper) Debugf(format string, args ...interface{}) {
	z.logger.Debug().Msgf(format, args...)
}

// Infof logs an info message with format
func (z *ZerologWrapper) Infof(format string, args ...interface{}) {
	z.logger.Info().Msgf(format, args...)
}

// Errorf logs an error message with format
func (z *ZerologWrapper) Errorf(format string, args ...interface{}) {
	z.logger.Error().Msgf(format, args...)
}

// Warnf logs a warning message with format
func (z *ZerologWrapper) Warnf(format string, args ...interface{}) {
	z.logger.Warn().Msgf(format, args...)
}

// Fatalf logs a fatal message with format and exits
func (z *ZerologWrapper) Fatalf(format string, args ...interface{}) {
	z.logger.Fatal().Msgf(format, args...)
}

// WithFields creates a new logger with additional fields
func (z *ZerologWrapper) WithFields(fields map[string]interface{}) Logger {
	ctx := z.logger.With()
	for k, v := range fields {
		ctx = z.addFieldToContext(ctx, k, v)
	}
	return &ZerologWrapper{logger: ctx.Logger()}
}

// WithContext creates a new logger with context
func (z *ZerologWrapper) WithContext(ctx context.Context) Logger {
	return &ZerologWrapper{logger: z.logger.With().Ctx(ctx).Logger()}
}

// WithComponent creates a new logger with a component field
func (z *ZerologWrapper) WithComponent(component string) Logger {
	return &ZerologWrapper{logger: z.logger.With().Str("component", component).Logger()}
}

func (z *ZerologWrapper) addFields(event *zerolog.Event, fields ...map[string]interface{}) {
	if len(fields) == 0 {
		return
	}

	for _, fieldMap := range fields {
		for k, v := range fieldMap {
			z.addFieldToEvent(event, k, v)
		}
	}
}

func (z *ZerologWrapper) addFieldToEvent(event *zerolog.Event, key string, value interface{}) {
	switch v := value.(type) {
	case string:
		event.Str(key, v)
	case int:
		event.Int(key, v)
	case int64:
		event.Int64(key, v)
	case float64:
		event.Float64(key, v)
	case bool:
		event.Bool(key, v)
	case time.Time:
		event.Time(key, v)
	case time.Duration:
		event.Dur(key, v)
	case error:
		event.Err(value.(error))
	default:
		event.Interface(key, v)
	}
}

func (z *ZerologWrapper) addFieldToContext(ctx zerolog.Context, key string, value interface{}) zerolog.Context {
	switch v := value.(type) {
	case string:
		return ctx.Str(key, v)
	case int:
		return ctx.Int(key, v)
	case int64:
		return ctx.Int64(key, v)
	case float64:
		return ctx.Float64(key, v)
	case bool:
		return ctx.Bool(key, v)
	case time.Time:
		return ctx.Time(key, v)
	case time.Duration:
		return ctx.Dur(key, v)
	default:
		return ctx.Interface(key, v)
	}
}

var globalLogger Logger

func InitGlobal(config *Config) {
	globalLogger = New(config)
}

func GetGlobal() Logger {
	if globalLogger == nil {
		globalLogger = New(DefaultConfig())
	}
	return globalLogger
}

func Debug(msg string, fields ...map[string]interface{}) {
	GetGlobal().Debug(msg, fields...)
}

func Info(msg string, fields ...map[string]interface{}) {
	GetGlobal().Info(msg, fields...)
}

func Error(msg string, fields ...map[string]interface{}) {
	GetGlobal().Error(msg, fields...)
}

func Warn(msg string, fields ...map[string]interface{}) {
	GetGlobal().Warn(msg, fields...)
}

func Fatal(msg string, fields ...map[string]interface{}) {
	GetGlobal().Fatal(msg, fields...)
}

func Debugf(format string, args ...interface{}) {
	GetGlobal().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	GetGlobal().Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	GetGlobal().Errorf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	GetGlobal().Warnf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	GetGlobal().Fatalf(format, args...)
}

func WithFields(fields map[string]interface{}) Logger {
	return GetGlobal().WithFields(fields)
}

func WithComponent(component string) Logger {
	return GetGlobal().WithComponent(component)
}

func WithContext(ctx context.Context) Logger {
	return GetGlobal().WithContext(ctx)
}
