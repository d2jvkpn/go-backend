package mod_user

import (
	"fmt"
	"net/mail"
	"regexp"
	"unicode/utf8"
)

var (
	ACCOUNT_Name_Min = 2
	ACCOUNT_Name_Max = 24

	_ACCOUNT_Phone = regexp.MustCompile(`^1[3456789]\d{9}$`)

	_ACCOUNT_Password = regexp.MustCompile(`^[a-zA-Z0-9]{16,64}$`) // !@.-_*

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

	if length < ACCOUNT_Name_Min || length > ACCOUNT_Name_Max {
		return fmt.Errorf("name length: %d~%d", ACCOUNT_Name_Min, ACCOUNT_Name_Max)
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

	if length > 128 {
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
