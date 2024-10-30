package mod_user

import (
	"context"
	// "errors"
	"fmt"
	"strings"
	"time"

	. "github.com/d2jvkpn/go-backend/internal/models"
	"github.com/d2jvkpn/go-backend/pkg/erri"
	"github.com/d2jvkpn/go-backend/pkg/infra"

	// "gorm.io/gorm"
	"github.com/d2jvkpn/errx"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	// "go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Account struct {
	// account id
	Id uuid.UUID `json:"id,omitempty" gorm:"column:id;type:uuid;default:gen_random_uuid();->" fake:"-" extensions:"x-order=01"`
	// created time RFC3339
	CreatedAt time.Time `json:"createdAt,omitempty" gorm:"column:created_at;autoCreateTime;->" fake:"-" extensions:"x-order=02"`
	// updated time
	UpdatedAt time.Time `json:"updatedAt,omitempty" gorm:"column:updated_at;autoUpdateTime;->" fake:"-" extensions:"x-order=03"`
	// account status
	Status string `json:"status,omitempty" gorm:"column:status;->;<-:create" fake:"activated" extensions:"x-order=04"`

	// firstname
	Firstname string `json:"firstname" gorm:"column:firstname" minLength:"2" maxLength:"24" example:"John" fake:"{firstname}" extensions:"x-order=05"`
	// lastname
	Lastname string `json:"lastname" gorm:"column:lastname" minLength:"2" maxLength:"24" example:"John" fake:"{lastname}" extensions:"x-order=06"`
	// Phone number
	Phone string `json:"phone" gorm:"column:phone;default:null" minLength:"6" maxLength:"20" example:"^1[3456789][0-9]{9}$" fake:"-" extensions:"x-order=07"` // fake:"{phone}"
	// email address
	Email string `json:"email,omitempty" gorm:"column:email;default:null" minLength:"5" maxLength:"128"  example:"john@noreply.local" fake:"{email}" extensions:"x-order=08"`
	// role
	Role string `json:"role" gorm:"column:role" fake:"{randomstring:[manager,employee,customer,contractor]}" extensions:"x-order=09"`
	// labels
	Labels pq.StringArray `json:"labels" gorm:"column:labels" fake:"fake" fakesize:"1" extensions:"x-order=10"`

	Password string `json:"-" gorm:"column:password" fake:"-" swaggerignore:"true"`
}

type CreateAccount struct {
	Status string `json:"status,omitempty" gorm:"column:status;->;<-:create" validate:"oneof=created activated" fake:"created" extensions:"x-order=04"`
	// firstname
	Firstname string `json:"firstname" gorm:"column:firstname" binding:"required" example:"John" fake:"{firstname}" extensions:"x-order=05"`
	// lastname
	Lastname string `json:"lastname" gorm:"column:lastname" binding:"required" example:"John" fake:"{lastname}" extensions:"x-order=06"`
	// Phone number
	Phone string `json:"phone" gorm:"column:phone;default:null" binding:"required,min=6,max=20" minLength:"6" maxLength:"20" example:"^1[3456789][0-9]{9}$" fake:"-" extensions:"x-order=07"`
	// email address
	Email string `json:"email,omitempty" gorm:"column:email;default:null" binding:"max=128" minLength:"5" maxLength:"128"  example:"john@noreply.local" fake:"{email}" extensions:"x-order=08"`
	// role
	Role string `json:"role" gorm:"column:role" validate:"oneof=manager employee customer contractor" fake:"{randomstring:[manager,employee,customer,contractor]}" extensions:"x-order=09"`
	// labels
	Labels pq.StringArray `json:"labels" gorm:"column:labels;type:varchar[]" fake:"fake" fakesize:"1" extensions:"x-order=10"`

	// password: [a-z][A-Z][0-9][!@.-_*]
	Password string `json:"password" gorm:"column:password" binding:"required,min=8,max=32" fake:"-" extensions:"x-order=11"`
}

func (self *CreateAccount) Validate() *errx.ErrX {
	var e error

	if e = ValidateName(self.Firstname, false); e != nil {
		return erri.Invalid(e, "fristname")
	}

	if e = ValidateName(self.Lastname, false); e != nil {
		return erri.Invalid(e, "lastname")
	}

	if self.Phone == "" && self.Email == "" {
		e = fmt.Errorf("phone or email is unset")
		return erri.Invalid(e, "phone or email is unset")
	}
	if e = ValidatePhone(self.Phone, true); e != nil {
		return erri.Invalid(e, "phone")
	}

	if e = ValidateEmail(self.Email, true); e != nil {
		return erri.Invalid(e, "email")
	}

	if e = ValidatePassword(self.Password); e != nil {
		return erri.Invalid(e, "password")
	}

	if e = _Validate.Struct(self); e != nil {
		if errs, ok := e.(validator.ValidationErrors); ok {
			fields := make([]string, len(errs))
			for i := range errs {
				fields[i] = errs[i].Error()
			}

			return erri.Invalid(e, strings.Join(fields, ","))
		} else {
			return erri.InternalErr(e, "validator")
		}
	}

	return nil
}

func (self *CreateAccount) hashPassword() (err *errx.ErrX) {
	var (
		bts []byte
		e   error
	)

	bts, e = bcrypt.GenerateFromPassword([]byte(self.Password), bcrypt.DefaultCost)
	if e != nil {
		return erri.InternalErr(e, "bcrypt")
	}
	self.Password = string(bts)

	return nil
}

func (self *CreateAccount) Do(ctx context.Context) (err *errx.ErrX) {
	var (
		bts []byte
		e   error

		tracer trace.Tracer
		span   trace.Span
	)

	tracer = otel.Tracer("mod_user.CreateAccount")

	if err = self.Validate(); err != nil {
		return err
	}

	_, span = tracer.Start(ctx, "bcrypt.GenerateFromPassword")
	bts, e = bcrypt.GenerateFromPassword([]byte(self.Password), bcrypt.DefaultCost)
	span.End()
	if e != nil {
		return erri.InternalErr(e, "bcrypt")
	}
	self.Password = string(bts)

	_, span = tracer.Start(ctx, "gorm.Create")
	e = Table(ctx, TABLE_UserAccounts, "id").Create(self).Error
	span.End()
	if e == nil {
		return nil
	}

	if infra.PgUniqueViolation(e) {
		errStr := e.Error()

		err = erri.BizErr(e, "already_exists")
		switch {
		case strings.Contains(errStr, "_email_key\""):
			err.WithMsg("email already exists")
		case strings.Contains(errStr, "_phone_key\""):
			err.WithMsg("email already exists")
		default:
			// TODO:
		}
		return err
	}

	return erri.InternalErr(e, "database")
}
