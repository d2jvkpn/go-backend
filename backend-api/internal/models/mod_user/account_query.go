package mod_user

import (
	"context"
	"fmt"
	"strings"

	. "backend-api/internal/models"
	"backend-api/internal/settings"
	"backend-api/pkg/infra"
	"backend-api/pkg/structs"
	"backend-api/pkg/utils"

	"github.com/d2jvkpn/errx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type QueryAccounts struct {
	ok bool

	// minimum: 15
	// maximum: 100
	// default: 15
	PageSize int `json:"pageSize" form:"pageSize" validate:"omitempty,gte=10,lte=100" extensions:"x-order=01"`

	// minimum: 1
	PageIndex int `json:"pageIndex" form:"pageIndex" validate:"omitempty,gt=0" extensions:"x-order=02"`

	// enum: createdAt,updatedAt,firstname,lastname
	// default: createdAt + desc
	Sort string `form:"sort" validate:"omitempty,oneof=createdAt updatedAt firstname lastname" extensions:"x-order=03"`

	// enum: asc,desc
	// default: asc
	Order string `form:"order" validate:"omitempty,oneof=asc desc" extensions:"x-order=04"`

	Keyword string `json:"keyword" form:"keyword" extensions:"x-order=10"`

	// enum: admin,editor,reviewer,user,guest
	// default:
	Level string `json:"level" form:"level" validate:"omitempty,oneof=admin editor reviewer user guest" extensions:"x-order=11"`

	// enum: created,activated,blocked
	// default: activated
	Status string `json:"status" form:"status" validate:"omitempty,oneof=created activated blocked" extensions:"x-order=12"`
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
		settings.Logger.Named("mod_user").Debug("validate", zap.Any("error", e))
		// fmt.Println("!!! validate QueryAccounts:", e)
		return structs.Invalid(e)
	}
	self.SetDefaults()

	self.ok = true
	return nil
}

func (self *QueryAccounts) db(ctx context.Context, flip bool) *gorm.DB {
	tx := Table(ctx, TABLE_UserAccounts+" t1")

	// TODO: elasticsearch
	if self.Keyword != "" {
		tx = tx.Where(
			"t1.firstname LIKE ? OR t1.lastname LIKE ? OR t1.phone LIKE ? OR t1.email LIKE ? OR ? = ANY(labels)",
			"%"+self.Keyword+"%", "%"+self.Keyword+"%", "%"+self.Keyword+"%", "%"+self.Keyword+"%", self.Keyword,
		)
	}

	if self.Level != "" {
		tx = tx.Where("t1.level = ?", self.Level)
	}

	if self.Status != "" {
		tx = tx.Where("t1.status = ?", self.Status)
	} else {
		tx = tx.Where("t1.status != 'deleted'")
	}

	if flip {
		orderSeg := fmt.Sprintf("t1.%s %s, t1.id", utils.ToSnakeCase(self.Sort), strings.ToUpper(self.Order))
		tx = infra.GormFlip(tx.Order(orderSeg), self.PageSize, self.PageIndex)
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
