package helper

import (
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Meta `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIRespose(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

func FormatValidationError(err error) interface{} {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}
	return errors
}
/* Own Custom Made Format
func FormatValidationError(err error) interface{} {
	var errors []string
	ok := reflect.TypeOf(reflect.Array) == reflect.TypeOf(err)
	if !ok {
		errors = append(errors, err.Error())
	} else {
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, e.Error())
		}
	}
	errorMessage := gin.H{"errors": errors}
	return errorMessage
}*/