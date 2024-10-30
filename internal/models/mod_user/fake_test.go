package mod_user

import (
	"fmt"
	"testing"

	. "github.com/d2jvkpn/go-backend/internal/models"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestFake01_Accounts(t *testing.T) {
	var (
		e        error
		accounts []CreateAccount
		tx       *gorm.DB
	)

	// 1.
	tx = Table(_TestCtx, TABLE_UserAccounts).
		Where("'fake' = any(labels)").Delete(nil)

	require.Nil(t, tx.Error)
	fmt.Printf("==> deleted fake account(s): %d\n", tx.RowsAffected)

	// 2.
	num := 10
	accounts = make([]CreateAccount, num)
	gofakeit.Slice(&accounts)
	for i := range accounts {
		accounts[i].Password = gofakeit.Password(true, true, true, true, false, 16)
		accounts[i].Labels = append(
			accounts[i].Labels,
			"passowrd="+accounts[i].Password,
		)
		e = accounts[i].hashPassword()
		require.Nil(t, e)
	}
	fmt.Printf("==> fake accounts: %v\n", accounts)

	e = Table(_TestCtx, TABLE_UserAccounts).CreateInBatches(accounts, len(accounts)).Error
	require.Nil(t, e)
}
