package utils

import (
	"fmt"

	"github.com/google/uuid"
)

var (
	_UUID_Null  uuid.UUID
	ErrIdIsNull = fmt.Errorf("id is null")
)

func UUIDIsNull(id uuid.UUID) bool {
	return id == _UUID_Null
}

func UUIDFromString(id string) (ans uuid.UUID, err error) {
	if ans, err = uuid.Parse(id); err != nil {
		return ans, err
	}
	if ans == _UUID_Null {
		return ans, ErrIdIsNull
	}

	return
}

func UUIDFromStrings(ids []string) (list []uuid.UUID, err error) {
	var uid uuid.UUID

	list = make([]uuid.UUID, 0, len(ids))
	for i := range ids {
		if uid, err = uuid.Parse(ids[i]); err != nil {
			return nil, err
		}
		if uid == _UUID_Null {
			return list, ErrIdIsNull
		}
		list = append(list, uid)
	}

	return list, nil
}
