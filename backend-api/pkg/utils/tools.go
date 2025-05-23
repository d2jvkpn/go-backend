package utils

import (
	"cmp"
	"fmt"
	"math"
	"math/rand/v2"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	Rand *rand.Rand
	// _Runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	_Alphanum = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	_RE_FirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	_RE_AllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func init() {
	source := rand.NewPCG(42, uint64(time.Now().UnixNano()))
	Rand = rand.New(source)
}

func UUID_Nil() uuid.UUID {
	return _UUID_Null
}

func UUIDNotNull(ids ...uuid.UUID) bool {
	return !slices.Contains(ids, _UUID_Null)
}

func UUIDUnique(ids []uuid.UUID) []uuid.UUID {
	slices.SortFunc(ids, func(a uuid.UUID, b uuid.UUID) int {
		s1, s2 := a.String(), b.String()
		switch {
		case s1 < s2:
			return -1
		case s1 == s2:
			return 0
		default:
			return 1
		}
	})

	return slices.Compact(ids)
}

func UUIDToStrings(ids []uuid.UUID) (list []string) {
	list = make([]string, len(ids))
	for i := range ids {
		list[i] = ids[i].String()
	}

	return list
}

func UniqueSlice[T cmp.Ordered](input []T) []T {
	slices.Sort(input)
	return slices.Compact(input)
}

func UniqueFilepath(t time.Time, dateDir bool, p ...string) string {
	/*
		if len(a) == 0 {
			a = []time.Time{time.Now()}
		}
	*/

	base := fmt.Sprintf("%s-%d_%s", t.Format("2006-01"), t.Unix(), uuid.New())

	if dateDir {
		p = append(p, t.Format("2006-01"), base)
	} else {
		p = append(p, base)
	}

	return filepath.Join(p...)
}

func Alphanum(length int) (bts []byte) {
	bts = make([]byte, length)

	for i := range bts {
		bts[i] = _Alphanum[Rand.IntN(len(_Alphanum))]
	}

	return bts
}

// https://stackoverflow.com/posts/56616250/revisions
func ToSnakeCase(str string) string {
	snake := _RE_FirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = _RE_AllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func Round3(value float64) float64 {
	return math.Round(value*1e3) / 1e3
}
