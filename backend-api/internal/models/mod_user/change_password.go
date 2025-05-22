package mod_user

import (
	"context"
	"fmt"

	. "backend-api/internal/models"
	"backend-api/pkg/structs"

	"github.com/d2jvkpn/errx"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	// "go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type ChangePassword struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (self *ChangePassword) Do(ctx context.Context, accountId uuid.UUID) (err *errx.ErrX) {
	var (
		bts []byte
		e   error

		tracer trace.Tracer
		span   trace.Span
	)

	tracer = otel.Tracer("mod_user.ChangePassword")
	// TODO: check self.NewPassword

	account := struct {
		Password string `gorm:"column:password"`
		Status   string `gorm:"column:status"`
	}{}

	_, span = tracer.Start(ctx, "getAccount")
	e = Table(ctx, TABLE_UserAccounts).Select("password", "status").Where("id = ?", accountId).Take(&account).Error
	span.End()
	if e != nil {
		return structs.PgNotFound(e)
	}

	if account.Status != "activated" {
		return structs.BizError(fmt.Errorf("account not activated"))
	}

	_, span = tracer.Start(ctx, "bcrypt.CompareHashAndPassword")
	e = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(self.OldPassword))
	span.End()
	if e != nil {
		return structs.BizError(e).WithCode("incorrect_password")
	}

	_, span = tracer.Start(ctx, "bcrypt.GenerateFromPassword")
	bts, e = bcrypt.GenerateFromPassword([]byte(self.NewPassword), bcrypt.DefaultCost)
	span.End()
	if e != nil {
		return structs.InternalError(e)
	}

	_, span = tracer.Start(ctx, "updatePassword")
	e = Table(ctx, TABLE_UserAccounts).Where("id = ?", accountId).Update("password", bts).Error
	span.End()
	if e != nil {
		return structs.InternalError(e).WithCode("change_password")
	}

	return nil
}
