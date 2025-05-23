package utils

import (
	"encoding/json"
	// "fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type IdName struct {
	Id   uuid.UUID `json:"id" gorm:"column:id;type:uuid" extensions:"x-order=01"`
	Name string    `json:"name" gorm:"column:name" extensions:"x-order=02"`
}

type KeyValue[T any] struct {
	Key   string `json:"key" gorm:"column:key" extensions:"x-order=01"`
	Value T      `json:"value" gorm:"column:value" extensions:"x-order=02"`
}

type IdNameValue[T any] struct {
	Id    uuid.UUID `json:"id" gorm:"column:id;type:uuid" extensions:"x-order=01"`
	Name  string    `json:"name" gorm:"column:name" extensions:"x-order=02"`
	Value T         `json:"value" gorm:"column:value" extensions:"x-order=03"`
}

// Query
type QueryPage struct {
	// 1. control parameters
	AId       string    `json:"-" form:"accountId"`
	AccountId uuid.UUID `json:"accountId" form:"-" swaggerignore:"true"`

	// 2. exactly macth

	// 3. general parameters
	PageSize  int    `json:"pageSize" form:"pageSize"  binding:"required,gte=10,lte=100" minimum:"10" maximum:"100" extensions:"x-order=01"`
	PageIndex int    `json:"pageIndex" form:"pageIndex" binding:"required,gt=0" minimum:"1" extensions:"x-order=02"`
	StartDate string `json:"startDate,omitepmty" form:"startDate" binding:"ltefield=EndDate" time_format:"2006-01-02" example:"2006-01-02" extensions:"x-order=03"`
	// created_at >= StartTime
	startTime []byte
	EndDate   string `json:"endDate,omitepmty" form:"endDate" time_format:"2006-01-02" example:"2006-01-02" extensions:"x-order=04"`
	// created_at < EndTime
	endTime []byte
	// 1. binding:"required" means this paramter must exists and not empty
	// 2. oneof="'' createdAt name" means this paramter must exists and but can be empty
	// 3. validate:"required" makes sure sagger marked as "* required"(red)
	Sort  string `form:"sort" binding:"oneof='' createdAt name" enums:",createdAt,name" validate:"required" enums:",createdAt,name" extensions:"x-order=05"`
	Order string `form:"order" binding:"oneof='' asc desc" enums:",asc,desc" validate:"required" enums:",asc,desc" extensions:"x-order=06"`

	// 4. other paramters
	Role string `json:"role,omitempty" form:"role" binding:"required,oneof=all manager employee" enums:"all,manager,employee" extensions:"x-order=07"`

	//
	// NOTE: don't use bool type here as binding can't parse "false" to fase
	Status string `json:"status,omitempty" form:"status" binding:"required,oneof=true false" enums:"true false" extensions:"x-order=08"`

	// 5. full text search
	Search string `json:"search,omitempty" form:"search" extensions:"x-order=09"`
}

/*
curl 'localhost:3061/api/v1/open/test_01?pageIndex=1&pageSize=10&status=true&role=all'
*/

func (self *QueryPage) FromString(jstr string) (err error) {
	return json.Unmarshal([]byte(jstr), self)
}

func (self *QueryPage) Validate() (err error) {
	if self.AId != "" {
		if self.AccountId, err = UUIDFromString(self.AId); err != nil {
			return err
		}
	}

	if len(self.StartDate) > 0 && len(self.EndDate) > 0 {
		start, end, e := DateRange(self.StartDate, self.EndDate)
		if e != nil {
			return e
		}

		self.startTime, self.endTime = pq.FormatTimestamp(start), pq.FormatTimestamp(end)
	}

	return nil
}

func (self *QueryPage) order() string {
	if self.Sort == "" {
		self.Sort = "created_at"
		if self.Order == "" {
			self.Order = "desc"
		}
	}

	if self.Order == "desc" {
		return ToSnakeCase(self.Order) + " DESC, id"
	} else {
		return ToSnakeCase(self.Order) + " ASC, id"
	}
}

// PageResult
type PageResult[T any] struct {
	PageSize  int   `json:"pageSize"`
	PageIndex int   `json:"pageIndex"`
	Total     int64 `json:"total"`
	Items     []T   `json:"items"`
}

func NewPageResult[T any](sizes ...int) (result *PageResult[T]) {
	result = new(PageResult[T])

	if len(sizes) == 0 {
		result.Items = make([]T, 0)
	} else {
		result.Items = make([]T, 0, sizes[0])
	}

	return result
}

type Items[T any] struct {
	Number int `json:"number"`
	Items  []T `json:"items"`
}

func NewItems[T any](elements []T) Items[T] {
	return Items[T]{Number: len(elements), Items: elements}
}

func IntoPageResult[T, V any](from *PageResult[T], into func(*T) V) (to PageResult[V]) {
	to.PageSize, to.PageIndex, to.Total = from.PageSize, from.PageIndex, from.Total

	to.Items = make([]V, len(from.Items))

	for i := range from.Items {
		to.Items[i] = into(&from.Items[i])
	}

	return to
}

func PageResultJoin[K comparable, T, V any](result *PageResult[T], exts map[K]V,
	getKey func(*T) K, setValue func(*T, V)) (n int) {

	if len(exts) == 0 {
		return 0
	}

	for i := range result.Items {
		if v, ok := exts[getKey(&result.Items[i])]; ok {
			n++
			setValue(&result.Items[i], v)
		}
	}

	return n
}
