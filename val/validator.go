package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile("^[a-z0-9_]+$").MatchString
	isValidFullname = regexp.MustCompile("^[a-zA-Z\\s]+$").MatchString
)

func ValidatorString(value string, minLeght int, maxLenght int) error {
	n := len(value)
	if n < minLeght || n > maxLenght {
		return fmt.Errorf("must contain form %d-%d characters: %s", minLeght, maxLenght, value)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidatorString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("must contain only lower case letter characters, or underscore")
	}
	return nil
}

func ValidatePassword(value string) error {
	return ValidatorString(value, 6, 100)
}

func ValidateEmail(value string) error {
	if err := ValidatorString(value, 3, 100); err != nil {
		return err
	}
	_, err := mail.ParseAddress(value)
	if err != nil {
		return fmt.Errorf("%s is not a valid email address", value)
	}
	return nil
}

func ValidateFullname(value string) error {
	if err := ValidatorString(value, 3, 100); err != nil {
		return err
	}
	if !isValidFullname(value) {
		return fmt.Errorf("must contain only letter, or spaces")
	}
	return nil
}
