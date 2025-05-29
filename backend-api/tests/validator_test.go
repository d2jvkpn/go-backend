package tests

import (
	// "fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestValidator(t *testing.T) {
	type Query struct {
		Size int    `validate:"omitempty,gte=10,lte=100"`
		Name string `validate:"required,min=1,max=4"`
	}

	var (
		query    Query
		validate *validator.Validate
	)

	validate = validator.New()
	query.Size = 10

	require.NotNil(t, validate.Struct(&query))

	query.Name = "你好世界" // rune
	require.Nil(t, validate.Struct(&query))
}
