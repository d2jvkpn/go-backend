package infra

import (
	"context"
	// "errors"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

/*
dsn:
- postgres://{USERANME}:{PASSWORD}@tcp({IP})/{DATABASE}?sslmode=disable
- "host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai"
*/
func PgConnect(vp *viper.Viper, release bool) (gormDB *gorm.DB, sqlDB *sql.DB, err error) {
	var conf *gorm.Config

	conf = &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	}

	if release {
		conf.Logger = logger.Default.LogMode(logger.Silent)
	}

	if gormDB, err = gorm.Open(postgres.Open(vp.GetString("dsn")), conf); err != nil {
		return nil, nil, err
	}

	if !release {
		gormDB = gormDB.Debug()
	}

	// Get generic database object sql.DB to use its functions
	if sqlDB, err = gormDB.DB(); err != nil {
		return nil, nil, err
	}

	if v := vp.GetInt("max_idle_conns"); v > 0 {
		sqlDB.SetMaxIdleConns(v)
	}
	if v := vp.GetInt("max_open_conns"); v > 0 {
		sqlDB.SetMaxOpenConns(v)
	}
	if v := vp.GetDuration("conn_max_idle_time"); v > 0 {
		sqlDB.SetConnMaxLifetime(v)
	}
	if v := vp.GetDuration("conn_max_lifetime"); v > 0 {
		sqlDB.SetConnMaxLifetime(v)
	}

	return gormDB, sqlDB, err
}

func PgNotFound(err error) bool {
	return err.Error() == "record not found"
	// return errors.Is(err, gorm.ErrDuplicatedKey) // !!! this doesn't work
}

// !!! works as expected only version(gorm.io/driver/postgres) <= v1.4.5
//
//	and errors.Is(e, gorm.ErrDuplicatedKey) doesn't work as expected
func PgUniqueViolation(err error) bool {
	/* when when gorm.Config.TranslateError is true
	println("~~~", err.Error())
	return err.Error() == "duplicated key not allowed"

	return errors.Is(err, gorm.ErrDuplicatedKey)
	*/

	/* when gorm.Config.TranslateError is false
	// ERROR: duplicate key value violates unique constraint \"TABLE_FIELD_key\" (SQLSTATE 23505)
	fmt.Printf("~~~ err: %[1]s\n    %#+[1]v, type: %[1]T\n", err)
	e, ok := err.(*pgx.PgError)
	fmt.Printf("~~~ e: %v, ok: %t\n", e, ok)
	return ok && e.Code == "23505"
	*/

	return strings.Contains(err.Error(), "(SQLSTATE 23505)")
}

func PgTimestamptz(t ...time.Time) []byte {
	if len(t) == 0 {
		t = []time.Time{time.Now()}
	}

	return pq.FormatTimestamp(t[0])
}

// #### concat to array
func PgConcateArrays[T any](tx *gorm.DB, field string, arr []T) clause.Expr {
	// gorm.Expr("? || ?", field, pq.Array(arr))
	return gorm.Expr(
		fmt.Sprintf("array_cat(%s, ?)", field),
		pq.Array(arr),
	)
}

// #### push to array and remove duplicates
func PgMergeArrays[T any](field string, arr []T) clause.Expr {
	return gorm.Expr(
		fmt.Sprintf("array(SELECT distinct unnest(%s || ?))", field),
		pq.Array(arr),
	)
}

func PgArrayContains[T any](tx *gorm.DB, field string, value T) *gorm.DB {
	return tx.Where(fmt.Sprintf("? = any(%s)", field), value)
}

func PgArrayIncludes[T any](tx *gorm.DB, field string, arr []T) *gorm.DB {
	return tx.Where(
		fmt.Sprintf("%s @> ?", field),
		pq.Array(arr),
	)
}

func PgArrayPush[T any](field string, value T) clause.Expr {
	return gorm.Expr(
		fmt.Sprintf(`array(SELECT distinct unnest(array_append(%s, ?)))`, field),
		value,
	)
}

func pg_e01_returning(ctx context.Context) {
	var tx *gorm.DB

	tx.WithContext(ctx).Table("table").Clauses(
		clause.Returning{Columns: []clause.Column{
			clause.Column{Name: "field1"},
			clause.Column{Name: "field2"},
		}},
	)
}

func pg_e02_on_conflict(ctx context.Context) {
	var (
		value any
		tx    *gorm.DB
	)

	tx.WithContext(ctx).Table("table").Clauses(
		clause.OnConflict{DoNothing: true},
	)

	tx.WithContext(ctx).Table("table").Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "key1"}, {Name: "key1"}},
			DoUpdates: clause.Assignments(map[string]any{"field1": value}),
		},
	)

	tx.WithContext(ctx).Table("table").Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "key1"}, {Name: "key2"}},
			DoUpdates: clause.AssignmentColumns([]string{"field1"}),
		},
	)
}

func pg_e03_json_field(table, key, field, typ string) string {
	return fmt.Sprintf("(%s.%s->>'%s')::%s", table, key, field, typ)
}

func pg_e04_array_has_overlap[T any](db *gorm.DB, field string, array []T) *gorm.DB {
	return db.Where(fmt.Sprintf("(%s && ?)", field), pq.Array(array))
}

func pg_e05_period_mins(field string, mins int) string {
	return fmt.Sprintf(
		`date_trunc('hour', %[1]s) + `+
			`(floor(extract(minute from %[1]s)/%[2]d) * interval '%[2]d min')`,
		field, mins,
	)
}

func PgTableExists(db *gorm.DB, table string) (existence bool, err error) {
	data := struct {
		Existence bool `gorm:"column:existence"`
	}{}

	err = db.Raw(`SELECT EXISTS (
	  SELECT 1 FROM information_schema.tables
	  WHERE table_schema = 'public' AND table_type = 'BASE TABLE' AND table_name = ?
	) existence;`, table,
	).Take(&data).Error

	if err != nil {
		return false, err
	}

	return data.Existence, nil
}

// "sys_.*", "sys_\\d+$"
func PgFindTables(db *gorm.DB, m string) (tables []string, err error) {
	tables = make([]string, 0)

	err = db.Raw(
		"SELECT table_name FROM information_schema.tables "+
			"WHERE table_schema = 'public' AND table_type = 'BASE TABLE' AND table_name ~ ?",
	).Pluck("table_name", &tables).Error

	if err != nil {
		return nil, err
	}

	return tables, nil
}
