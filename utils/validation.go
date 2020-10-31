package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

// CommonError - error response format
type CommonError struct {
	Errors map[string]interface{} `json:"errors"`
}

func ValidatonError(err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)

	for _, v := range errs {
		field := strings.ToLower(v.Field())

		switch v.Tag() {
		case "required":
			res.Errors[field] = fmt.Sprintf("%v is %v", field, v.Tag())
		case "sinteger":
			res.Errors[field] = fmt.Sprintf("%v is number only", field)
		case "sgte":
			res.Errors[field] = fmt.Sprintf("%v must higher than equal %v", field, v.Param())
		case "slte":
			res.Errors[field] = fmt.Sprintf("%v must less than equal %v", field, v.Param())
		case "max":
			res.Errors[field] = fmt.Sprintf("%v must less than %v character", field, v.Param())
		case "min":
			res.Errors[field] = fmt.Sprintf("%v must higher than %v character", field, v.Param())
		case "email":
			res.Errors[field] = fmt.Sprintf("%v is not a valid email address", v.Value())
		case "username":
			res.Errors[field] = fmt.Sprintf("%v is not a valid username", v.Value())
		}
	}

	return res
}
