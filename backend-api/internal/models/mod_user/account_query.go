package mod_user

import (
	"context"
	// "fmt"

	. "backend-api/internal/models"
	"backend-api/pkg/infra"
	"backend-api/pkg/structs"
	"backend-api/pkg/utils"

	"github.com/d2jvkpn/errx"
	"gorm.io/gorm"
)

type QueryAccounts struct {
	ok bool

	// minimum: 15
	// maximum: 100
	// default: 15
	PageSize int `json:"pageSize" form:"pageSize" validate:"omitempty,gte=15,lte=100" extensions:"x-order=01"`

	// minimum: 1
	PageIndex int `json:"pageIndex" form:"pageIndex" validate:"omitempty,gt=0" extensions:"x-order=02"`

	// enum: createdAt,updatedAt,firstname,lastname
	// default: createdAt
	Sort string `form:"sort" validate:"omitempty,oneof=createdAt updatedAt firstname lastname" extensions:"x-order=03"`

	// enum: asc,desc
	// default: asc
	Order string `form:"order" validate:"omitempty,oneof=asc desc" extensions:"x-order=04"`

	Search string `json:"search" form:"search" extensions:"x-order=10"`

	// enum: admin,editor,reviewer,user,guest
	// default:
	Level string `json:"level" form:"level" validate:"omitempty,oneof=admin editor reviewer user guest" extensions:"x-order=11"`

	// enum: created,activated,blocked,deleted
	// default: activated
	Status string `json:"status" form:"status" validate:"omitempty,oneof=created activated blocked deleted" extensions:"x-order=12"`
}

func (self *QueryAccounts) SetDefaults() {
	if self.PageSize <= 0 {
		self.PageSize = 15
	}

	if self.PageIndex <= 0 {
		self.PageIndex = 1
	}

	if self.Sort == "" {
		self.Sort = "createdAt"
		self.Order = "desc"
	}

	if self.Order == "" {
		self.Order = "asc"
	}

	if self.Status == "" {
		self.Status = "activated"
	}
}

func (self *QueryAccounts) Validate() (err *errx.ErrX) {
	var e error

	if self.ok {
		return nil
	}

	if e = _Validate.Struct(self); e != nil {
		return structs.Invalid(e)
	}
	self.SetDefaults()

	self.ok = true
	return nil
}

func (self *QueryAccounts) order() string {
	if self.Order == "desc" {
		return "t1." + utils.ToSnakeCase(self.Sort) + " DESC, t1.id"
	} else {
		return "t1." + utils.ToSnakeCase(self.Sort) + " ASC, t1.id"
	}
}

func (self *QueryAccounts) db(ctx context.Context, flip bool) *gorm.DB {
	tx := Table(ctx, TABLE_UserAccounts+" t1")

	if self.Search != "" {
		tx = tx.Where(
			"t1.firstname LIKE ? OR t1.lastname LIKE ? OR t1.email LIKE ?",
			"%"+self.Search+"%", "%"+self.Search+"%", "%"+self.Search+"%",
		)
	}

	if self.Level != "" {
		tx = tx.Where("t1.level = ?", self.Level)
	}

	if self.Status != "" {
		tx = tx.Where("t1.status = ?", self.Status)
	}

	if flip {
		tx = infra.GormFlip(tx.Order(self.order()), self.PageSize, self.PageIndex)
	}

	return tx
}

func (self *QueryAccounts) Do(ctx context.Context) (result *utils.PageResult[Account], err *errx.ErrX) {
	var e error

	if err = self.Validate(); err != nil {
		return nil, err
	}

	result = utils.NewPageResult[Account]()

	//
	result.PageIndex, result.PageSize, result.Total = self.PageIndex, self.PageSize, 0

	if self.PageIndex == 1 {
		if e = self.db(ctx, false).Select("COUNT(1)").Count(&result.Total).Error; e != nil {
			return nil, structs.InternalError(e).WithCode("query_accounts")
		}

		if result.Total == 0 {
			return result, nil
		}
	}

	e = self.db(ctx, true).Omit("password").Find(&result.Items).Error
	if e != nil {
		return nil, structs.InternalError(e).WithCode("query_accounts")
	}

	return result, nil
}
