package utils_test

import (
	"encoder/framework/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsJson(t *testing.T) {
	json := `
    {
      "id": "some",
      "test": true
    }
  `

	err := utils.IsJson(json)
	require.Nil(t, err)

	json = "some"

	err = utils.IsJson(json)
	require.Error(t, err)
}
