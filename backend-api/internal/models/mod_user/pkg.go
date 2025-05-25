package mod_user

import (
	"fmt"
	"net/mail"
	"regexp"
	"unicode/utf8"
)

var (
	//_ACCOUNT_Phone = regexp.MustCompile(`^1[3456789]\d{9}$`)
	_ACCOUNT_Phone = regexp.MustCompile(`^\d{6,20}$`)

	_ACCOUNT_Password = regexp.MustCompile(`^[a-zA-Z0-9]{8,32}$`) // !@.-_*

	_ACCOUNT_PasswordContains = []*regexp.Regexp{
		regexp.MustCompile(`[a-z]\+`),
		regexp.MustCompile(`[A-Z]\+`),
		regexp.MustCompile(`[0-9]\+`),
	}
)

func ValidateName(name string, allowEmpty bool) error {
	length := utf8.RuneCountInString(name)
	if length == 0 && allowEmpty {
		return nil
	}

	if length < 1 || length > 32 {
		return fmt.Errorf("name length: %d~%d", 1, 32)
	}

	return nil
}

func ValidateEmail(email string, allowEmpty bool) (err error) {
	length := len(email)

	if length == 0 && allowEmpty {
		return nil
	}

	if length < 5 {
		return fmt.Errorf("email length is too short")
	}

	if length > 64 {
		return fmt.Errorf("email length is too long")
	}

	if _, err = mail.ParseAddress(email); err != nil {
		return err
	}

	return nil
}

func ValidatePhone(phone string, allowEmpty bool) (err error) {
	length := len(phone)

	if length == 0 && allowEmpty {
		return nil
	}

	if !_ACCOUNT_Phone.Match([]byte(phone)) {
		return fmt.Errorf("invalid phone")
	}

	return nil
}

func ValidatePassword(password string) (err error) {
	bts := []byte(password)

	if !_ACCOUNT_Password.Match(bts) {
		return fmt.Errorf("invalid password")
	}

	for i := range _ACCOUNT_PasswordContains {
		if _ACCOUNT_PasswordContains[i].Match(bts) {
			return fmt.Errorf("invalid password")
		}
	}

	return nil
}

// TODO: ValidateLabels: 16, 32
