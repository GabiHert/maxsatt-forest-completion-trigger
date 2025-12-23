package logger

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lsgndln/dd-trace-go/ddtrace/tracer"
)

func init() {
	log.SetFlags(0)
}

type logPattern struct {
	Metadata      any     `json:"metadata,omitempty"`
	TraceId       any     `json:"dd.trace_id,omitempty"`
	SpanId        any     `json:"dd.span_id,omitempty"`
	RequestId     *string `json:"requestId,omitempty"`
	Timestamp     string  `json:"timestamp"`
	Level         string  `json:"level"`
	Service       string  `json:"service,omitempty"`
	Operation     string  `json:"operation,omitempty"`
	ErrorMsg      string  `json:"exceptions,omitempty"`
	Message       string  `json:"detail"`
	CorrelationId string  `json:"correlationId"`
	TransactionId string  `json:"transactionId"`
}

type logLevel string

const (
	errorLevel logLevel = "ERROR"
	warnLevel  logLevel = "WARN"
	infoLevel  logLevel = "INFO"
	debugLevel logLevel = "DEBUG"
	traceLevel logLevel = "TRACE"
	offLevel   logLevel = "OFF"
)

func Error(ctx context.Context, err error, message string, metadata ...any) {
	lCtx := GetContext(ctx)
	logMessage := generateLogMessage(lCtx, errorLevel, message, err, metadata)
	if shouldPrintLog(errorLevel) {
		log.Println(logMessage)
	}
}

func Warn(ctx context.Context, err error, message string, metadata ...any) {
	lCtx := GetContext(ctx)
	logMessage := generateLogMessage(lCtx, warnLevel, message, err, metadata)
	if shouldPrintLog(warnLevel) {
		log.Println(logMessage)
	}
}

func Info(ctx context.Context, message string, metadata ...any) {
	lCtx := GetContext(ctx)
	logMessage := generateLogMessage(lCtx, infoLevel, message, nil, metadata)
	if shouldPrintLog(infoLevel) {
		log.Println(logMessage)
	}
}

func Debug(ctx context.Context, message string, metadata ...any) {
	lCtx := GetContext(ctx)
	logMessage := generateLogMessage(lCtx, debugLevel, message, nil, metadata)
	if shouldPrintLog(debugLevel) {
		log.Println(logMessage)
	}
}

func Trace(ctx context.Context, message string, metadata ...any) {
	lCtx := GetContext(ctx)
	logMessage := generateLogMessage(lCtx, traceLevel, message, nil, metadata)
	if shouldPrintLog(traceLevel) {
		log.Println(logMessage)
	}
}

func shouldPrintLog(level logLevel) bool {
	switch strings.Replace(strings.ToUpper(os.Getenv("LOG_LEVEL")), " ", "", -1) {
	case string(errorLevel):
		return level == errorLevel
	case string(warnLevel):
		return level == errorLevel || level == warnLevel
	case string(infoLevel):
		return level == infoLevel || level == errorLevel || level == warnLevel
	case string(debugLevel):
		return level == infoLevel || level == errorLevel || level == debugLevel || level == warnLevel
	case string(traceLevel):
		return true
	case string(offLevel):
		return false
	default:
		return true
	}
}

func generateLogMessage(ctx Ctx, level logLevel, message string, err error, metadata ...[]any) string {
	logg := &logPattern{}
	logg.Level = string(level)
	logg.Message = message
	logg.Operation = *getCaller(3)
	logg.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	logg.RequestId = getRequestId(ctx)
	logg.CorrelationId = getCorrelationId(ctx)
	logg.TransactionId = getTransactionId(ctx)

	if span, ok := tracer.SpanFromContext(ctx.Context()); ok && parseBool(os.Getenv("DD_TRACE_ENABLED")) {
		logg.SpanId = strconv.FormatUint(span.Context().SpanID(), 10)
		logg.TraceId = strconv.FormatUint(span.Context().TraceID(), 10)
	}

	if len(metadata[0]) > 0 {
		metadataMap := make(map[string]any)
		for _, value := range metadata {
			metadataMap[entityName(value)] = value
		}
		logg.Metadata = metadataMap
	}

	if err != nil {
		errorMsg, _ := json.Marshal(errorModel(ctx, err))
		logg.ErrorMsg = string(errorMsg)
	}

	logMessage, _ := json.Marshal(logg)
	return string(logMessage)
}

func getRequestId(ctx Ctx) *string {
	return ctx.GetRequestId()
}

func getCorrelationId(ctx Ctx) string {
	return ctx.GetCorrelationId()
}

func getTransactionId(ctx Ctx) string {
	return ctx.getTransactionId()
}
