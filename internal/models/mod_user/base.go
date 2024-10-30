package mod_user

import (
	"context"
	// "fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	_DB *gorm.DB
	// _Logger *zap.Logger
)

func Init(ctx context.Context, db *gorm.DB) (err error) {
	_DB = db
	// _Logger = settings.Logger.Named("mod_user").WithOptions(zap.WithCaller(true))

	return nil
}

func Table(ctx context.Context, table string, returing ...string) *gorm.DB {
	if len(returing) == 0 {
		return _DB.WithContext(ctx).Table(table)
	}

	cols := make([]clause.Column, 0, len(returing))
	for _, v := range returing {
		cols = append(cols, clause.Column{Name: v})
	}

	return _DB.WithContext(ctx).Table(table).Clauses(clause.Returning{Columns: cols})
}

func DB(ctx context.Context) *gorm.DB {
	return _DB.WithContext(ctx)
}
