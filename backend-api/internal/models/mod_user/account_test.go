package mod_user

import (
	"fmt"
	"testing"

	"github.com/d2jvkpn/errx"
	"github.com/stretchr/testify/require"
)

func TestAccountValidation(t *testing.T) {
	var (
		err     *errx.ErrX
		account *CreateAccount
	)

	account = &CreateAccount{
		Firstname: "Jane",
		Lastname:  "Doe",
	}

	err = account.Validate()
	require.NotNil(t, err)
	fmt.Printf("==> error: %v\n", err)
}
