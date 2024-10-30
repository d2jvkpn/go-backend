package mod_user

import (
	// "context"
	// "errors"
	// "fmt"
	"time"

	// "gorm.io/gorm"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	Phone string `json:"phone" gorm:"column:phone;default:null" binding:"required,min=6,max=20" minLength:"6" maxLength:"20" example:"^1[3456789][0-9]{9}$" fake:"-" extensions:"x-order=07"` // fake:"{phone}"
	// email address
	Email string `json:"email,omitempty" gorm:"column:email;default:null" binding:"max=128" minLength:"5" maxLength:"128"  example:"john@noreply.local" fake:"{email}" extensions:"x-order=07"`
	// role
	Role   string         `json:"role" gorm:"column:role" fake:"{randomstring:[manager,employee,customer,contractor]}" extensions:"x-order=08"`
	Labels pq.StringArray `json:"labels" gorm:"column:labels" fake:"fake" fakesize:"1" extensions:"x-order=09"`

	Password string `json:"-" gorm:"column:password" fake:"-" swaggerignore:"true"`
}
