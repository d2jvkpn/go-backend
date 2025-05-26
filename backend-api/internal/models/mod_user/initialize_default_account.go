package mod_user

import (
	"context"
	"errors"
	"fmt"

	. "backend-api/internal/models"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

/*
initialize default account

init:

	firstname: d2jvkpn
	lastname: Doe
	email: d2jvkpn@local.dev
	password: cKEvRvCM6Tme1zQ0
*/
func InitializeDefaultAccount(ctx context.Context, vp *viper.Viper) (err error) {
	var (
		count   int64
		bts     []byte
		account CreateAccount
	)

	err = Table(ctx, TABLE_UserAccounts).Select("COUNT(1)").Limit(1).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	bts, err = bcrypt.GenerateFromPassword([]byte(vp.GetString("password")), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	account.Firstname = vp.GetString("firstname")
	account.Lastname = vp.GetString("lastname")
	account.Email = vp.GetString("email")
	account.Password = string(bts)
	account.Level, account.Status = "admin", "activated"

	if account.Firstname == "" || account.Lastname == "" {
		return errors.New("firstname or lastname not set")
	}

	if account.Email == "" || account.Password == "" {
		return errors.New("email or password not set")
	}

	if err = Table(ctx, TABLE_UserAccounts).Create(&account).Error; err != nil {
		return fmt.Errorf("InitializeDefaultAccount: %w", err)
	}

	return nil
}
