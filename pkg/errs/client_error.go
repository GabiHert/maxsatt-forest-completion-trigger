package errs

import (
	"encoding/json"
	"errors"
	"io"
)

type clientError struct {
	message      string
	description  string
	code         string
	errorDetails []ErrorDetails
	statusCode   int
}

func ClientError(description, code string, status int, responseBody ...io.ReadCloser) BaseError {
	err := clientError{
		statusCode:   status,
		message:      "client error occurred",
		description:  description,
		code:         code,
		errorDetails: nil,
	}

	err.readResponseBody(responseBody)

	return err
}

func (c clientError) ToDto(transactionId string) Error {
	return Error{
		Error: ErrorData{
			Description:  c.description,
			Code:         c.code,
			Id:           transactionId,
			ErrorDetails: c.errorDetails,
		},
	}
}

func (c clientError) Error() string {
	return c.message
}

func (c clientError) Description() string {
	return c.description
}

func (c clientError) StatusCode() int {
	return c.statusCode
}

func (c clientError) SetStatusCode(statusCode int, description string) {
	c.statusCode = statusCode
	c.description = description
}

func (c clientError) Is(target error) bool {
	return errors.As(target, &c)
}

func (c clientError) Code() string {
	return c.code
}

func (c clientError) ErrorDetails() map[string]any {
	var details = make(map[string]any)
	for _, detail := range c.errorDetails {
		details["attribute"] = detail.Attribute
		details["messages"] = detail.Messages
	}
	return details
}

func (clientError) Retries() int {
	return 0
}

func (clientError) InternalNotify() bool {
	return true
}

func (clientError) Abort() bool {
	return true
}

func (c *clientError) readResponseBody(responseBody []io.ReadCloser) {
	if len(responseBody) == 0 {
		return
	}
	var response map[string]any
	err := json.NewDecoder(responseBody[0]).Decode(&response)
	if err != nil {
		return
	}

	if err, ok := response["error"].(map[string]any); ok {
		if desc, exists := err["description"]; exists {
			c.description = desc.(string)
		}
		if code, exists := err["code"]; exists {
			c.code = code.(string)
		}
	}
}
