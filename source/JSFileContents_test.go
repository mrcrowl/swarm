package source

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRegisterLineForBundleNilImports(t *testing.T) {
	actual := getRegisterLineForBundle("bob", nil)
	expected := getRegisterLineForBundle("bob", []string{})
	assert.Equal(t, expected, actual)
}
