package mod_user

import (
	"context"
	"fmt"
	"strings"

	. "backend-api/internal/models"
	"backend-api/pkg/infra"
	"backend-api/pkg/structs"
	"backend-api/pkg/utils"

	"github.com/d2jvkpn/errx"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	//"gorm.io/gorm"
	//"gorm.io/gorm/clause"
)

type EditAccount struct {
	// account id(uuid)
	// reqired: true
	AccountId string    `json:"-" form:"accountId" gorm:"-" fake:"-" extensions:"x-order=01"`
	id        uuid.UUID `json:"-" fake:"-"`

	// required: true
	// minLength: 1
	// maxLength: 32
	// example: John
	Firstname string `json:"firstname" gorm:"column:firstname" validate:"required,min=1,max=32" fake:"{firstname}" extensions:"x-order=02"`

	// minLength: 1
	// maxLength: 32
	// example: Doe
	Lastname string `json:"lastname" gorm:"column:lastname" validate:"required,min=1,max=32" fake:"{lastname}" extensions:"x-order=03"`

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
	// default:
	Password string `json:"password" gorm:"column:password" validate:"omitempty,min=8,max=32" fake:"-" extensions:"x-order=08"`
}

func (self *EditAccount) Validate() *errx.ErrX {
	var e error

	if self.id, e = utils.UUIDFromString(self.AccountId); e != nil {
		return structs.Invalid(e).WithCode("invalid_id").WithMsg("firstname")
	}

	if e = ValidateName(self.Firstname, false); e != nil {
		return structs.Invalid(e).WithCode("invalid_firstname").WithMsg("firstname")
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

func (self *EditAccount) Do(ctx context.Context) (err *errx.ErrX) {
	var (
		fields []interface{}
		bts    []byte
		e      error
	)

	self.Email = strings.ToLower(self.Email)
	if err = self.Validate(); err != nil {
		return err
	}

	fields = make([]interface{}, 0, 7)
	fields = append(fields, "firstname", "lastname", "email", "level", "labels")

	if self.Password != "" {
		bts, e = bcrypt.GenerateFromPassword([]byte(self.Password), bcrypt.DefaultCost)
		if e != nil {
			return structs.InternalError(e).WithCode("bcrypt")
		}
		self.Password = string(bts)
		fields = append(fields, "password")
	}

	e = Table(ctx, TABLE_UserAccounts).Where("id = ?", self.AccountId).
		Select(fields[0], fields[1:]...).
		UpdateColumns(self).Error
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
