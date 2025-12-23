package tracer

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/lsgndln/dd-trace-go/ddtrace/ext"
	"github.com/lsgndln/dd-trace-go/ddtrace/tracer"
)

type fiberTracer struct {
}

func FiberTracer() *fiberTracer {
	return &fiberTracer{}
}

func (f *fiberTracer) Handle(ctx *fiber.Ctx) error {
	var headers = make(tracer.HTTPHeadersCarrier)
	for key, values := range ctx.GetReqHeaders() {
		for _, value := range values {
			headers.Set(key, value)
		}
	}

	opts := []tracer.StartSpanOption{
		tracer.SpanType(ext.SpanTypeWeb),
		tracer.ServiceName(os.Getenv("DD_SERVICE")),
		tracer.Tag(ext.HTTPMethod, ctx.Method()),
		tracer.Tag(ext.HTTPURL, string(ctx.Request().URI().PathOriginal())),
		tracer.Measured(),
	}

	spanCtx, err := tracer.Extract(headers)
	if err == nil {
		opts = append(opts, tracer.ChildOf(spanCtx))
	}

	span, extractedCtx := tracer.StartSpanFromContext(ctx.Context(), ctx.Method()+" "+string(ctx.Request().URI().PathOriginal()), opts...)
	defer span.Finish()

	ctx.SetUserContext(extractedCtx)

	err = ctx.Next()

	if span != nil {
		status := ctx.Response().StatusCode()
		if status == 0 {
			status = http.StatusOK
		}
		span.SetTag(ext.HTTPCode, strconv.Itoa(status))
	}

	return err
}
