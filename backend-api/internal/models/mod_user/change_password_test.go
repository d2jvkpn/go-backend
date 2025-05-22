package mod_user

import (
	"fmt"
	"testing"

	"backend-api/pkg/utils"

	"github.com/d2jvkpn/errx"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// -- 513a20cb-cd80-435f-bdd3-b4ce32392d2b  mSe5KRUACL1wKz2S  Gc7SX7810aGIIW09
func TestChangePassword(t *testing.T) {
	var (
		e              error
		err            *errx.ErrX
		accountId      uuid.UUID
		changePassword ChangePassword
	)

	if _TestFlag.NArg() != 3 {
		return
	}
	fmt.Printf("~~~ args: %v\n", _TestFlag.Args())

	accountId, e = utils.UUIDFromString(_TestFlag.Args()[0])
	require.Nil(t, e)

	changePassword.OldPassword = _TestFlag.Args()[1]
	changePassword.NewPassword = _TestFlag.Args()[2]

	err = changePassword.Do(_TestCtx, accountId)
	require.Nil(t, err)
}
