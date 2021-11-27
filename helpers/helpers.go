package helpers

import (
	"strings"

	"github.com/go-playground/validator"
)

type ReqLogin struct {
	Username string `validate:"required,min=6"`
	Password string `validate:"required,min=6"`
}
type ErrorResponseLogin struct {
	FailedField string
	Message       string
}

func ValidateStruct(req ReqLogin) []*ErrorResponseLogin {
	var errors []*ErrorResponseLogin
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
	
			var data = []string{"This",err.Field(), "must be", err.Tag(), err.Param()}
			var element ErrorResponseLogin
			element.FailedField = err.Field()	
			element.Message = strings.Join(data, " ")
			errors = append(errors, &element)
		}
	}
	
	return errors
}