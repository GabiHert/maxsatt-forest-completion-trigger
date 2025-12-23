package logger

import (
	"context"
)

type Logger interface {
	SetCorrelationID(ctx context.Context, correlationId string)
	GetTransactionID(ctx context.Context) string
	Trace(ctx context.Context, message string, metadata ...any)
	Info(ctx context.Context, message string, metadata ...any)
	Debug(ctx context.Context, message string, metadata ...any)
	Warn(ctx context.Context, err error, message string, metadata ...any)
	Error(ctx context.Context, err error, message string, metadata ...any)
}

type loggerWrapper struct {
	logger Logger
}

func LoggerWrapper() Logger {
	return &loggerWrapper{}
}

func (l *loggerWrapper) GetCorrelationID(ctx context.Context) string {
	logCtx := GetContext(ctx)
	return logCtx.GetCorrelationId()
}

func (l *loggerWrapper) GetTransactionID(ctx context.Context) string {
	logCtx := GetContext(ctx)
	return logCtx.getTransactionId()
}

func (l *loggerWrapper) SetCorrelationID(ctx context.Context, correlationId string) {
	logCtx := GetContext(ctx)
	logCtx.SetCorrelationId(correlationId)
}

func (l *loggerWrapper) Trace(ctx context.Context, message string, metadata ...any) {
	Trace(ctx, message, metadata...)
}

func (l *loggerWrapper) Info(ctx context.Context, message string, metadata ...any) {
	Info(ctx, message, metadata...)
}

func (l *loggerWrapper) Debug(ctx context.Context, message string, metadata ...any) {
	Debug(ctx, message, metadata...)
}

func (l *loggerWrapper) Warn(ctx context.Context, err error, message string, metadata ...any) {
	Warn(ctx, err, message, metadata...)
}

func (l *loggerWrapper) Error(ctx context.Context, err error, message string, metadata ...any) {
	Error(ctx, err, message, metadata...)
}
