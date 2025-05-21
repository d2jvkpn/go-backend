package mod_user

import (
	"context"
	"errors"
	// "fmt"

	. "backend-api/internal/models"
	"backend-api/pkg/structs"

	"github.com/d2jvkpn/errx"
	// "gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	// "go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type AccountLogin struct {
	Phone    string `json:"phone,omitempty" example:"^1[3456789][0-9]{9}$" extensions:"x-order=01"`
	Email    string `json:"email,omitempty" example:"john@noreply.local" extensions:"x-order=02"`
	Password string `json:"password" extensions:"x-order=03"`
}

func (self *AccountLogin) Validate() (password []byte, err *errx.ErrX) {
	var e error

	password = []byte(self.Password)

	if self.Phone == "" && self.Email == "" {
		const msg = "no account info: phone, email"
		return nil, structs.Invalid(errors.New(msg)).WithMsg(msg)
	}

	/*
		if err = ValidatePassword(self.NewPassword) {
			return err
		}
	*/

	if e = ValidatePhone(self.Phone, true); e != nil {
		return nil, structs.Invalid(e).WithMsg("invalid phone or password")
	}

	if e = ValidateEmail(self.Email, true); e != nil {
		return nil, structs.Invalid(e).WithMsg("invalid email or password")
	}

	return password, nil
}

/*
#### steps
1. validate parameters
2. find target account
3. check account.Status
4. bcrypt verify password
5. check account status and expiration
*/
func (self *AccountLogin) Do(ctx context.Context) (account *Account, err *errx.ErrX) {
	var (
		e        error
		password []byte

		tracer trace.Tracer
		span   trace.Span
	)

	if password, err = self.Validate(); err != nil {
		return nil, err
	}

	tracer = otel.Tracer("mod_user.AccountLogin")

	tx := Table(ctx, TABLE_UserAccounts)
	if self.Phone != "" {
		tx = tx.Where("phone = ?", self.Phone)
	} else if self.Email != "" {
		tx = tx.Where("email = ?", self.Email)
	}

	account = new(Account)
	// tx.Select("id", "name", "phone", "email", "status", "level", "status", "password", "expiration")

	_, span = tracer.Start(ctx, "GetAccount")
	e = tx.Take(account).Error
	span.End()

	if e != nil {
		/*
			if infra.GormPgNotFound(e) {
				return nil, BizError(e, "account_not_found")
			}

			return nil, InternalError(e, "find_account")
		*/
		return nil, structs.PgNotFound(e).WithMsg("account or password is incorrect")
	}

	if err = account.IsOK(); err != nil {
		return nil, err
	}

	if account.Password == "" {
		return account, structs.BizError(errors.New("hashed password is empty")).
			WithCode("account_not_available")
	}

	_, span = tracer.Start(ctx, "bcrypt.CompareHashAndPassword")
	e = bcrypt.CompareHashAndPassword([]byte(account.Password), password)
	span.End()

	if e != nil {
		return account, structs.AuthError(e).
			WithCode("wrong_account_or_password").
			WithMsg("")
	}

	return account, nil
}
