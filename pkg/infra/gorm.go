package infra

import (
	"database/sql"
	"fmt"
	// "errors"
	"bytes"
	"encoding/json"
	"strings"
	"time"

	"database/sql/driver"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GormCloseDB(db *gorm.DB) (err error) {
	var sqlDB *sql.DB

	if db == nil {
		return nil
	}

	if sqlDB, err = db.DB(); err != nil {
		return err
	}

	return sqlDB.Close()
}

func GormFlip(tx *gorm.DB, pageSize, pageIndex int) *gorm.DB {
	if pageIndex <= 0 || pageSize <= 0 {
		return tx.Limit(0)
	}

	return tx.Limit(pageSize).Offset((pageIndex - 1) * pageSize)
}

// Abandoned, GormEquals(_DB.Table("accounts"), "id", ids).Delete(nil)
func GormEquals[T comparable](tx *gorm.DB, field string, values []T) *gorm.DB {
	if len(values) == 0 {
		return tx.Where("0 = 1")
	}

	equals, items := make([]string, 0, len(values)), make([]any, 0, len(values))

	for _, k := range values {
		equals = append(equals, fmt.Sprintf("%s = ?", field))
		items = append(items, k)
	}

	return tx.Where(strings.Join(equals, " OR "), items...)
}

// Abandoned
func SQLxOrStrings[T uuid.UUID | string](field string, ids []T) string {
	if len(ids) <= 0 {
		return "0 = 1"
	}

	equals := make([]string, len(ids))
	for i := 0; i < len(ids); i++ {
		equals[i] = fmt.Sprintf("%s = '%s'", field, ids[i])
	}

	return strings.Join(equals, " OR ")
}

// func SQLxInStrings[T uuid.UUID | string](ids []T) string {
func SQLxInStrings[T uuid.UUID](ids []T) string {
	if len(ids) <= 0 {
		return "(NULL)"
	}

	equals := make([]string, len(ids))
	for i := 0; i < len(ids); i++ {
		equals[i] = fmt.Sprintf("'%s'", ids[i])
	}

	return "(" + strings.Join(equals, ", ") + ")"
}

func GormFmtExpr(state string, a []any, b ...any) clause.Expr {
	if len(a) > 0 {
		state = fmt.Sprintf(state, a...)
	}

	return gorm.Expr(state, b...)
}

func SQLxCol(table string) func(...string) string {
	return func(fields ...string) string {
		if len(fields) > 0 {
			return fmt.Sprintf("%s.%s", table, fields[0])
		} else {
			return table
		}
	}
}

func GormTransaction(db *gorm.DB, state string, a ...any) (err error) {
	tx := db.Begin()

	if err = tx.Exec(state, a...).Error; err != nil {
		return tx.Rollback().Error
	} else {
		return tx.Commit().Error
	}
}

// https://github.com/jackc/pgx/issues/1781
type GormUUIDs []uuid.UUID

func NewGormUUIDs(a ...int) GormUUIDs {
	if len(a) > 0 {
		return make([]uuid.UUID, 0, a[0])
	} else {
		return make([]uuid.UUID, 0)
	}
}

func (self *GormUUIDs) FromStrings(ids []string) (err error) {
	var id uuid.UUID

	for _, v := range ids {
		if v == "" {
			continue
		}

		if id, err = uuid.Parse(v); err != nil {
			return err
		}
		*self = append(*self, id)
	}

	return nil
}

// Scan implements the sql.Scanner interface for postgres uuid[] deserialization.
// e.g. single: {23e4aa0c-8591-4544-8cc1-8348d9d7f901}
// e.g. multiple: {23e4aa0c-8591-4544-8cc1-8348d9d7f901,aee911e4-3fa8-4866-95ac-e4878d3e06e5}
func (self *GormUUIDs) Scan(value any) error {
	var (
		ok   bool
		err  error
		str  string
		id   uuid.UUID
		strs []string
	)

	if str, ok = value.(string); !ok {
		return fmt.Errorf("failed to scan: %v", value)
	}

	if len(str) < 3 { // "{}", ""
		*self = make([]uuid.UUID, 0)
		return nil
	}

	strs = strings.Split(str[1:len(str)-1], ",")
	*self = make([]uuid.UUID, 0, len(strs))
	for _, v := range strs {
		if id, err = uuid.Parse(v); err != nil {
			return err
		}

		*self = append(*self, id)
	}

	return nil
}

// Value implements the driver.Valuer interface for postgres uuid[] serialization.
func (self GormUUIDs) Value() (driver.Value, error) {
	var strs []string

	strs = make([]string, len(self))

	for i := range self {
		strs[i] = (self)[i].String()
	}

	if len(strs) == 0 {
		return "{}", nil
	}

	return fmt.Sprintf("{%s}", strings.Join(strs, ",")), nil
}

// Date
type GormDate string

func (self *GormDate) Validate() (err error) {
	if _, err = time.Parse(time.DateOnly, string(*self)); err != nil {
		return err
	}

	return nil
}

func (self *GormDate) Scan(value interface{}) error {
	var (
		ok bool
		at time.Time
	)

	if at, ok = value.(time.Time); !ok {
		return fmt.Errorf("failed to scan: %v", value)
	}

	*self = GormDate(at.Format(time.DateOnly))

	return nil
}

func (self GormDate) Value() (driver.Value, error) {
	var (
		err error
		at  time.Time
	)

	if at, err = time.Parse(time.DateOnly, string(self)); err != nil {
		return nil, err
	}

	return at.Format(time.DateOnly), nil
}

// GormJSONArray, json array
type GormJSONArray[T any] []T

func (self *GormJSONArray[T]) Scan(value any) (err error) {
	bts, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan: %v", value)
	}
	bts = bytes.TrimSpace(bts)

	result := make(GormJSONArray[T], 0, 5)
	if len(bts) > 0 { // empty []uint8 cause error
		err = json.Unmarshal(bts, &result)
	}

	*self = result
	return err
}

func (self GormJSONArray[T]) Value() (driver.Value, error) {
	return json.Marshal(self)
}

// GormJSONMap, json map
type GormJSONMap[T any] map[string]T

func (self *GormJSONMap[T]) Scan(value any) (err error) {
	bts, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan: %v", value)
	}
	bts = bytes.TrimSpace(bts)

	if len(bts) > 0 { // empty []uint8 causes error
		err = json.Unmarshal(bts, &self)
	}

	return err
}

func (self GormJSONMap[T]) Value() (driver.Value, error) {
	return json.Marshal(self)
}

func (self GormJSONMap[T]) JSON() (bts []byte) {
	bts, _ = json.Marshal(self)
	return bts
}

// GormJSON
type GormJSONAny[T any] struct {
	Item T
}

func (self *GormJSONAny[T]) Scan(value any) (err error) {
	bts, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan: %v", value)
	}
	bts = bytes.TrimSpace(bts)

	if len(bts) > 0 { // empty []uint8 causes error
		err = json.Unmarshal(bts, &self.Item)
	}

	return err
}

func (self GormJSONAny[T]) Value() (driver.Value, error) {
	return json.Marshal(self.Item)
}

func (self GormJSONAny[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(&self.Item)
}

func (self *GormJSONAny[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &self.Item)
}
