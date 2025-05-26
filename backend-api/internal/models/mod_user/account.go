package mod_user

import (
	"context"
	"errors"
	//"fmt"
	"slices"
	"time"

	. "backend-api/internal/models"
	"backend-api/pkg/structs"
	"backend-api/pkg/utils"

	"github.com/d2jvkpn/errx"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	// minLength: 1
	// maxLength: 32
	// example: John
	Firstname string `json:"firstname" gorm:"column:firstname" fake:"{firstname}" extensions:"x-order=05"`

	// lastname
	// minLength: 1
	// maxLength: 32
	// example: John
	Lastname string `json:"lastname" gorm:"column:lastname" fake:"{lastname}" extensions:"x-order=06"`

	// Phone number
	// minLength: 6
	// maxLength: 20
	// example: ^1[3456789][0-9]{9}$
	Phone string `json:"phone" gorm:"column:phone;default:null" fake:"-" extensions:"x-order=07"` // fake:"{phone}"

	// email address
	// minLength: 5
	// maxLength: 64
	// example: john@noreply.local
	Email string `json:"email,omitempty" gorm:"column:email;default:null" fake:"{email}" extensions:"x-order=08"`

	// level
	// enum: admin,editor,reviewer,user,guest
	Level string `json:"level" gorm:"column:level" fake:"{randomstring:[admin,editor,reviewer,user,guest]}" extensions:"x-order=09"`

	// labels, max length of a label is 32
	// maxLength: 16
	Labels pq.StringArray `json:"labels" gorm:"column:labels;default:null;type:varchar[]" fake:"fake" fakesize:"1" swaggertype:"array,string" extensions:"x-order=10"`

	Password string `json:"-" gorm:"column:password" fake:"-" swaggerignore:"true"`
}

func HashPassword(password string) (bts []byte, err *errx.ErrX) {
	var e error

	bts, e = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if e != nil {
		return nil, structs.InternalError(e).WithCode("bcrypt")
	}

	return bts, nil
}

func (self *Account) IsOK() (err *errx.ErrX) {
	if self.Status != "activated" {
		return structs.BizError(errors.New("account hasn't been actived")).WithCode("not_activated")
	}

	return nil
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
