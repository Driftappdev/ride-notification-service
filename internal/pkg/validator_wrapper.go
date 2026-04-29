package pkg

import govalidator "github.com/driftappdev/libpackage/resilience/validator"

type Rule = govalidator.Rule
type Validator = govalidator.Validator
type FieldRules = govalidator.FieldRules
type ValidationErrors = govalidator.ValidationErrors

func NewValidator() *Validator { return govalidator.New() }
func ValidateMap(fields FieldRules, data map[string]interface{}) error {
	return govalidator.Validate(fields, data)
}
func ValidationErrorsOf(err error) (ValidationErrors, bool) {
	return govalidator.AsValidationErrors(err)
}
