package pkg

import govalidate "github.com/driftappdev/libpackage/resilience/validate"

func Required(name, value string) error        { return govalidate.Required(name, value) }
func MaxLen(name, value string, max int) error { return govalidate.MaxLen(name, value, max) }
func ValidateEmail(value string) error         { return govalidate.Email(value) }
func OneOf(name, value string, allowed ...string) error {
	return govalidate.OneOf(name, value, allowed...)
}
