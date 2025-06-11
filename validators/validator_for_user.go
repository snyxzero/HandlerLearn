package validators

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
)

type ValidatorForUser struct {
	validator *validator.Validate
}

func NewValidatorForUser() *ValidatorForUser {
	validator := validator.New()
	validator.RegisterValidation("russian_name", validateRussianName)
	validator.RegisterValidation("gmail_email", validateGmailEmail)
	return &ValidatorForUser{
		validator: validator,
	}
}

func validateRussianName(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[а-яА-ЯёЁ]+$`)
	return re.MatchString(fl.Field().String())
}
func validateGmailEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	return strings.HasSuffix(email, "@gmail.com")
}

func (obj *ValidatorForUser) Struct(i interface{}) error {
	return obj.validator.Struct(i)
}
