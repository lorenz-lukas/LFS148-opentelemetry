package logging

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	name   string
	module string
	base   *zap.Logger
}

func New(name, module string) *Logger {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey: "message",
		LineEnding: zapcore.DefaultLineEnding,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)

	return &Logger{
		name:   name,
		module: module,
		base:   zap.New(core),
	}
}

func (l *Logger) WithModule(module string) *Logger {
	return &Logger{
		name:   l.name,
		module: module,
		base:   l.base,
	}
}

func (l *Logger) Info(format string, args ...any) {
	l.base.Info(l.message("INFO", format, args...))
}

func (l *Logger) Warn(format string, args ...any) {
	l.base.Warn(l.message("WARNING", format, args...))
}

func (l *Logger) Error(format string, args ...any) {
	l.base.Error(l.message("ERROR", format, args...))
}

func (l *Logger) message(level, format string, args ...any) string {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	msg = strings.ReplaceAll(msg, "\n", " ")
	return fmt.Sprintf("%s:%s:%s:%s", level, l.name, l.module, msg)
}
