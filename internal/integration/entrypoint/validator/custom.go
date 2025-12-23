package validator

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	errors2 "github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/entrypoint/errors"

	"github.com/go-playground/validator/v10"
)

type Validate interface {
	Struct(code string, s interface{}) error
}

type customValidator struct {
	validator *validator.Validate
}

func CustomValidator() Validate {
	validate := validator.New()
	_ = validate.RegisterValidation("not_nil", notNil)
	_ = validate.RegisterValidation("oneof_ignore_case", oneOfIgnoreCase)
	_ = validate.RegisterValidation("datetime_format", dateTimeFormat)
	_ = validate.RegisterValidation("hour_minute", hourMinute)
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &customValidator{
		validator: validate,
	}
}

func notNil(_ validator.FieldLevel) bool {
	return true
}

func oneOfIgnoreCase(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()
	allowedValues := strings.Split(fl.Param(), " ")

	for _, allowedValue := range allowedValues {
		if strings.EqualFold(fieldValue, allowedValue) {
			return true
		}
	}

	return false
}

func dateTimeFormat(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()

	if fieldValue == "" {
		return true
	}

	format := fl.Param()

	_, err := time.Parse(format, fieldValue)

	return err == nil
}

func hourMinute(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()

	hourMinuteString := strings.Split(fieldValue, ":")

	if len(hourMinuteString) != 2 {
		return false
	}

	hour, err := strconv.Atoi(hourMinuteString[0])
	if err != nil || hour < 0 || hour > 23 {
		return false
	}

	minute, err := strconv.Atoi(hourMinuteString[1])
	if err != nil || minute < 0 || minute > 59 {
		return false
	}

	return true
}

func (c *customValidator) Struct(code string, s interface{}) error {
	err := c.validator.Struct(s)
	if err != nil {
		return errors2.ValidationError(err, code)
	}
	return nil
}
