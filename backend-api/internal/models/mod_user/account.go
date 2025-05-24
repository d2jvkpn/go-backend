package mod_user

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	. "backend-api/internal/models"
	"backend-api/pkg/infra"
	"backend-api/pkg/structs"
	"backend-api/pkg/utils"

	"github.com/d2jvkpn/errx"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	// minLength: 2
	// maxLength: 24
	// example: John
	Firstname string `json:"firstname" gorm:"column:firstname" fake:"{firstname}" extensions:"x-order=05"`

	// lastname
	// minLength: 2
	// maxLength: 24
	// example: John
	Lastname string `json:"lastname" gorm:"column:lastname" fake:"{lastname}" extensions:"x-order=06"`

	// Phone number
	// minLength: 6
	// maxLength: 20
	// example: ^1[3456789][0-9]{9}$
	Phone string `json:"phone" gorm:"column:phone;default:null" fake:"-" extensions:"x-order=07"` // fake:"{phone}"

	// email address
	// minLength: 5
	// maxLength: 128
	// example: john@noreply.local
	Email string `json:"email,omitempty" gorm:"column:email;default:null" fake:"{email}" extensions:"x-order=08"`

	// level
	// enum: admin,editor,reviewer,user,guest
	Level string `json:"level" gorm:"column:level" fake:"{randomstring:[admin,editor,reviewer,user,guest]}" extensions:"x-order=09"`

	// labels
	Labels pq.StringArray `json:"labels" gorm:"column:labels;type:varchar[]" fake:"fake" fakesize:"1" swaggertype:"array,string" extensions:"x-order=10"`

	Password string `json:"-" gorm:"column:password" fake:"-" swaggerignore:"true"`
}

type CreateAccount struct {
	// required: true
	// enum: created,activated
	Status string `json:"status,omitempty" gorm:"column:status;->;<-:create" validate:"required,oneof=created activated" fake:"activated" extensions:"x-order=01"`

	// required: true
	// minLength: 2
	// maxLength: 32
	// example: John
	Firstname string `json:"firstname" gorm:"column:firstname" validate:"required,min=2,max=32" fake:"{firstname}" extensions:"x-order=02"`

	// minLength: 2
	// maxLength: 32
	// example: Doe
	Lastname string `json:"lastname" gorm:"column:lastname" validate:"required,min=2,max=32" fake:"{lastname}" extensions:"x-order=03"`

	// Phone number
	// minLength: 6
	// maxLength: 20
	// example: ^1[3456789][0-9]{9}$
	Phone string `json:"phone" gorm:"column:phone;default:null" validate:"omitempty,min=6,max=20" fake:"-" extensions:"x-order=04"`

	// minLength: 5
	// maxLength: 64
	// example: john@noreply.local
	Email string `json:"email,omitempty" gorm:"column:email;default:null" validate:"omitempty,min=5,max=64" fake:"{email}" extensions:"x-order=05"`

	// required: true
	// enum: admin,editor,reviewer,user,guest
	Level string `json:"level" gorm:"column:level" validate:"required,oneof=admin editor reviewer user guest" fake:"{randomstring:[admin,editor,reviewer,user,guest]}" swaggertype:"array,string" extensions:"x-order=06"`

	// labels
	Labels pq.StringArray `json:"labels" gorm:"column:labels;type:varchar[]" validate:"max=16" fake:"fake" fakesize:"1" extensions:"x-order=07"`

	// password: ^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[^]{8,32}$
	// minLength: 8
	// maxLength: 32
	Password string `json:"password" gorm:"column:password" validate:"required,min=8,max=32" fake:"-" extensions:"x-order=08"`

	// return account id
	Id uuid.UUID `json:"-" gorm:"column:id;type:uuid;default:gen_random_uuid();->" fake:"-" extensions:"x-order=09"`
}

func (self *Account) IsOK() (err *errx.ErrX) {
	if self.Status != "activated" {
		return structs.BizError(errors.New("account hasn't been actived")).WithCode("not_activated")
	}

	return nil
}

func (self *CreateAccount) Validate() *errx.ErrX {
	var e error

	if e = ValidateName(self.Firstname, false); e != nil {
		return structs.Invalid(e).WithCode("invalid_first").WithMsg("fristname")
	}

	if e = ValidateName(self.Lastname, false); e != nil {
		return structs.Invalid(e).WithCode("invalid_lastname").WithMsg("lastname")
	}

	if self.Phone == "" && self.Email == "" {
		e = fmt.Errorf("phone or email is unset")
		return structs.Invalid(e).WithCode("no_contact").WithMsg("phone or email is unset")
	}
	if e = ValidatePhone(self.Phone, true); e != nil {
		return structs.Invalid(e).WithCode("invalid_phone").WithMsg("phone")
	}

	if e = ValidateEmail(self.Email, true); e != nil {
		return structs.Invalid(e).WithCode("invalid_email").WithMsg("email")
	}

	if self.Password != "" {
		if e = ValidatePassword(self.Password); e != nil {
			return structs.Invalid(e).WithCode("invalid_password").WithMsg("invalid password")
		}
		self.Status = "activated"
	} else {
		self.Status = "created"
	}

	if e = _Validate.Struct(self); e != nil {
		if errs, ok := e.(validator.ValidationErrors); ok {
			/* can't complie with go1.24
			fields := make([]string, len(errs))
			for i := range errs {
				fields[i] = errs[i].Error()
			}

			return structs.Invalid(e).WithMsg(strings.Join(fields, ","))
			*/
			return structs.Invalid(e).WithMsg("%v", errs)
		} else {
			return structs.InternalError(e).WithCode("validator")
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
		return structs.InternalError(e).WithCode("bcrypt")
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
		return structs.InternalError(e).WithCode("bcrypt")
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

		err = structs.BizError(e).WithCode("already_exists")
		switch {
		case strings.Contains(errStr, "_email_key\""):
			err.WithCode("exists_email").WithMsg("email already exists")
		case strings.Contains(errStr, "_phone_key\""):
			err.WithCode("exists_phone").WithMsg("phone already exists")
		default:
			// TODO:
		}
		return err
	}

	return structs.InternalError(e).WithCode("database")
}

type UpdateStatus struct {
	// account id(uuid)
	AccountId string `form:"accountId"`
	accountId uuid.UUID

	// origin status
	// enum: created,activated,blocked
	// required: true
	Status string `form:"status" validate:"oneof=created activated blocked"`

	// new status
	// enum: activated,blocked,deleted
	// required: true
	NewStatus string `form:"newStatus" validate:"oneof=activated blocked deleted"`
}

func (self *UpdateStatus) Validate() (err *errx.ErrX) {
	var (
		e              error
		validNewStatus bool
	)

	if e = _Validate.Struct(self); e != nil {
		return structs.Invalid(e).WithCode("validate_failed")
	}

	if self.accountId, e = utils.UUIDFromString(self.AccountId); e != nil {
		return structs.Invalid(e).WithCode("invalid_accountId")
	}

	switch self.Status {
	case "created":
		validNewStatus = slices.Contains([]string{"activated", "deleted"}, self.NewStatus)
	case "activated":
		validNewStatus = slices.Contains([]string{"blocked", "deleted"}, self.NewStatus)
	case "blocked":
		validNewStatus = slices.Contains([]string{"activated", "deleted"}, self.NewStatus)
	default:
		return structs.Invalid(e).WithCode("invalid_status")
	}

	if !validNewStatus {
		return structs.Invalid(e).WithCode("invalid_newStatus")
	}

	return nil
}

/*
UPDATE your_table SET status = 'new_value' WHERE id = 123 RETURNING OLD.status;

BEGIN;
SELECT status FROM your_table WHERE id = 123 FOR UPDATE;
UPDATE your_table SET status = 'new_value' WHERE id = 123;
COMMIT;
*/
func (self *UpdateStatus) Do(ctx context.Context) (originStatus string, err *errx.ErrX) {
	if err = self.Validate(); err != nil {
		return "", err
	}

	var e1, e2 error

	e2 = Table(ctx, TABLE_UserAccounts).Transaction(func(tx *gorm.DB) error {
		if e1 = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", self.AccountId).Pluck("status", &originStatus).Error; e1 != nil {
			return e1
		}

		if e1 = tx.Where("id = ?", self.AccountId).
			Update("status", self.NewStatus).Error; e1 != nil {
			return e1
		}

		return nil
	})

	if e2 != nil {
		return "", structs.InternalError(e2).WithCode("database")
	}

	return originStatus, nil
}
