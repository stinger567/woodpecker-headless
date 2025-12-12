package permissions

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.woodpecker-ci.org/woodpecker/v3/server/model"
)

func TestAdmins(t *testing.T) {
	a := NewAdmins([]string{"woodpecker-ci"})
	assert.True(t, a.IsAdmin(&model.Account{AccountName: "woodpecker-ci"}))
	assert.False(t, a.IsAdmin(&model.Account{AccountName: "not-woodpecker-ci"}))
	empty := NewAdmins([]string{})
	assert.False(t, empty.IsAdmin(&model.Account{AccountName: "woodpecker-ci"}))
	assert.False(t, empty.IsAdmin(&model.Account{AccountName: "not-woodpecker-ci"}))
}
