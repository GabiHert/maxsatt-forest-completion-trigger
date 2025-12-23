package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/lsgndln/dd-trace-go/ddtrace/tracer"
)

type loggerAdapter interface {
	Debug(ctx context.Context, message string, metadata ...any)
	Error(ctx context.Context, err error, message string, metadata ...any)
}

type Client interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

type customHttp struct {
	client *http.Client
	logger loggerAdapter
}

func HttpClient(timeout time.Duration, logger loggerAdapter) Client {
	client := &http.Client{
		Timeout: timeout,
	}

	return &customHttp{
		client: client,
		logger: logger,
	}
}

func (c *customHttp) withContext(ctx context.Context, req *http.Request) *http.Client {
	if span, found := tracer.SpanFromContext(ctx); found && os.Getenv("DD_TRACE_ENABLED") == "true" {
		_ = tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(req.Header))
	}
	return c.client
}

func (c *customHttp) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	filteredHeaders := make(map[string][]string)
	for k, v := range req.Header {
		if k == "Authorization" {
			filteredHeaders[k] = []string{"Bearer ******"}
		} else {
			filteredHeaders[k] = v
		}
	}

	requestBody := readRequestBody(req)
	c.logger.Debug(ctx, "Starting", map[string]any{
		"url":     req.URL.String(),
		"method":  req.Method,
		"headers": filteredHeaders,
		"request": filterPayloadForLog(requestBody),
	})

	resp, err := c.withContext(ctx, req).Do(req)

	var errMessage *string
	if err != nil {
		value := err.Error()
		errMessage = &value
	}
	var status *int
	if resp != nil {
		value := resp.StatusCode
		status = &value
	}
	responseBody := readResponseBody(resp)

	if err != nil || (status != nil && *status >= 400) {
		c.logger.Error(ctx, err, "Finished", map[string]any{
			"url":      req.URL.String(),
			"method":   req.Method,
			"headers":  filteredHeaders,
			"request":  filterPayloadForLog(requestBody),
			"response": filterPayloadForLog(responseBody),
			"status":   status,
			"err":      errMessage,
		})
	} else {
		c.logger.Debug(ctx, "Finished", map[string]any{
			"url":      req.URL.String(),
			"method":   req.Method,
			"headers":  filteredHeaders,
			"request":  filterPayloadForLog(requestBody),
			"response": filterPayloadForLog(responseBody),
			"status":   status,
			"err":      errMessage,
		})
	}
	return resp, err
}

func readRequestBody(req *http.Request) map[string]any {
	if req == nil || req.Body == nil {
		return nil
	}

	defer req.Body.Close()
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil
	}

	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var responseBody map[string]any
	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		return nil
	}

	return responseBody
}

func readResponseBody(resp *http.Response) map[string]any {
	if resp == nil || resp.Body == nil {
		return nil
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var responseBody map[string]any
	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		return nil
	}

	return responseBody
}

func filterPayloadForLog(payload map[string]any) map[string]any {
	if payload == nil {
		return nil
	}

	filteredPayload := make(map[string]any)
	for k, v := range payload {
		switch k {
		case "Authorization":
			filteredPayload[k] = "Bearer ******"
		case "access_token":
			filteredPayload[k] = "******"
		default:
			filteredPayload[k] = v
		}
	}

	return filteredPayload
}
