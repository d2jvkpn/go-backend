package mod_user

import (
	"context"
	//"fmt"

	. "backend-api/internal/models"
	"backend-api/pkg/structs"

	"github.com/d2jvkpn/errx"
	"github.com/google/uuid"
)

func DeleteAccounts(ctx context.Context, accountIds []uuid.UUID) (count int64, err *errx.ErrX) {
	if len(accountIds) == 0 {
		return 0, nil
	}

	result := Table(ctx, TABLE_UserAccounts).
		Where("id IN ? AND status != 'deleted'", accountIds).
		Update("status", "deleted")
	if result.Error != nil {
		return 0, structs.InternalError(result.Error).WithCode("database")
	}

	return result.RowsAffected, nil
}
