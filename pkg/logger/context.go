package logger

import (
	"context"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lsgndln/dd-trace-go/ddtrace/tracer"
)

func StartDatadogAgent() {
	if parseBool(os.Getenv("DD_TRACE_ENABLED")) {
		tracer.Start(
			tracer.WithAgentAddr(os.Getenv("DD_AGENT_HOST")+":8126"),
			tracer.WithGlobalTag("env", os.Getenv("ENV")),
		)
	}
}

type Ctx interface {
	Stop()
	Context() context.Context
	getContextByKey(key string) context.Context
	GetRequestId() *string
	GetLogGroup() *string
	GetLogStream() *string
	GetCorrelationId() string
	getTransactionId() string
	SetCorrelationId(correlationId string)
	SetReceiveCount(receiveCount int)
	GetReceiveCount() int
	SetOperationUniqueId(id *string)
	GetOperationUniqueId() *string
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
}

type loggerCtx struct {
	operationUniqueId *string
	logGroup          *string
	logStream         *string
	requestId         *string
	ctxMap            map[string]*ctxCarrier
	uniqueId          string
	correlationId     string
	transactionId     string
	receiveCount      int
}

type ctxCarrier struct {
	ctx      context.Context
	parent   *string
	name     string
	children []string
}

func (c *loggerCtx) Deadline() (deadline time.Time, ok bool) {
	return c.Context().Deadline()
}

func (c *loggerCtx) Done() <-chan struct{} {
	return c.Context().Done()
}

func (c *loggerCtx) Err() error {
	return c.Context().Err()
}

func (c *loggerCtx) Value(key any) any {
	return c.Context().Value(key)
}

func GetContext(carrier ...any) Ctx {
	var ctx context.Context
	if carrier != nil && len(carrier) > 0 {
		switch carrierValue := carrier[0].(type) {
		case *fiber.Ctx:
			ctx = carrierValue.UserContext()
		case fiber.Ctx:
			ctx = carrierValue.UserContext()
		case Ctx:
			return carrierValue
		case context.Context:
			ctx = carrierValue
		default:
			ctx = context.Background()
		}
	} else {
		ctx = context.Background()
	}

	var requestId *string
	var logGroup *string
	var logStream *string
	if ctx.Value("requestId") != nil {
		requestIdVal := ctx.Value("requestId").(string)
		requestId = &requestIdVal
	}
	if lambdaContext, ok := lambdacontext.FromContext(ctx); ok {
		requestId = &lambdaContext.AwsRequestID
		logGroupName := os.Getenv("AWS_LAMBDA_LOG_GROUP_NAME")
		if logGroupName != "" {
			logGroup = &logGroupName
		}
		logStreamName := os.Getenv("AWS_LAMBDA_LOG_STREAM_NAME")
		if logStreamName != "" {
			logStream = &logStreamName
		}
	}

	var correlationId string
	if ctx.Value("correlationId") != nil {
		correlationId = ctx.Value("correlationId").(string)
	} else {
		correlationId = uuid.New().String()
	}

	var transactionId string
	if ctx.Value("transactionId") != nil {
		transactionId = ctx.Value("transactionId").(string)
	} else {
		transactionId = uuid.New().String()
	}

	caller := getCaller(2)
	if _, found := tracer.SpanFromContext(ctx); !found {
		_, ctx = tracer.StartSpanFromContext(ctx, *caller)
	}

	uniqueId := uuid.New().String()
	return &loggerCtx{
		uniqueId:      uniqueId,
		logGroup:      logGroup,
		logStream:     logStream,
		requestId:     requestId,
		correlationId: correlationId,
		transactionId: transactionId,
		ctxMap: map[string]*ctxCarrier{
			uniqueId + *caller: {
				ctx:      ctx,
				name:     *getCaller(2),
				children: make([]string, 0),
			},
		},
	}
}

func (c *loggerCtx) Stop() {
	for key, carrier := range c.ctxMap {
		if span, found := tracer.SpanFromContext(carrier.ctx); found && strings.HasPrefix(key, c.uniqueId) {
			span.Finish()
		}
	}
}

func (c *loggerCtx) Context() context.Context {
	offset := 0
	for {
		if caller := getCaller(offset); caller != nil {
			if carrier, ok := c.ctxMap[c.uniqueId+*caller]; ok {
				return carrier.ctx
			} else {
				offset++
			}
		} else {
			return context.TODO()
		}
	}
}

func (c *loggerCtx) getContextByKey(key string) context.Context {
	if carrier, found := c.ctxMap[c.uniqueId+key]; found {
		return carrier.ctx
	}
	return context.TODO()
}

func (c *loggerCtx) GetRequestId() *string {
	return c.requestId
}

func (c *loggerCtx) GetLogGroup() *string {
	return c.logGroup
}

func (c *loggerCtx) GetLogStream() *string {
	return c.logStream
}

func (c *loggerCtx) GetCorrelationId() string {
	return c.correlationId
}

func (c *loggerCtx) getTransactionId() string {
	return c.transactionId
}

func (c *loggerCtx) SetCorrelationId(correlationId string) {
	c.correlationId = correlationId
}

func (c *loggerCtx) SetReceiveCount(receiveCount int) {
	c.receiveCount = receiveCount
}

func (c *loggerCtx) GetReceiveCount() int {
	return c.receiveCount
}

func (c *loggerCtx) SetOperationUniqueId(id *string) {
	if c.operationUniqueId != nil && *c.operationUniqueId != "" {
		c.operationUniqueId = id
	}
}

func (c *loggerCtx) GetOperationUniqueId() *string {
	return c.operationUniqueId
}

func getCaller(offset int) *string {
	if pc, _, _, ok := runtime.Caller(offset); ok {
		names := strings.Split(runtime.FuncForPC(pc).Name(), "/")
		name := strings.TrimSuffix(names[len(names)-1], ".func1")
		if strings.Contains(name, "logger.(*loggerWrapper)") {
			return getCaller(offset + 2)
		}
		return &name
	}
	return nil
}
