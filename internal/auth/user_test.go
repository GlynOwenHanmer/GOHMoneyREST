package auth_test

import (
	"testing"

	"github.com/glynternet/mon/internal/auth"
	"github.com/glynternet/mon/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestUserClaimsAuthoriser_NewClaims(t *testing.T) {
	a := auth.UserClaimsAuthoriser(model.User{})
	claims := a.NewClaims()
	assert.NotNil(t, claims)
	assert.Equal(t, model.User{}, *claims.(*model.User))
}

func TestUserClaimsAuthoriser(t *testing.T) {
	for _, errorCase := range []struct {
		name       string
		authorised model.User
		claims     interface{}
	}{
		{
			name:   "wrong type",
			claims: "hiyer",
		},
		{
			name:       "wrong email",
			authorised: model.User{Email: "A"},
			claims:     &model.User{Email: "boop"},
		},
		{
			name:   "nil user pointer",
			claims: (*model.User)(nil),
		},
		{
			name:       "expected validated email",
			authorised: model.User{EmailVerified: true},
			claims:     &model.User{},
		},
	} {
		t.Run("unauthorised/"+errorCase.name, func(t *testing.T) {
			a := auth.UserClaimsAuthoriser(errorCase.authorised)
			err := a.Authorise(errorCase.claims)
			assert.Error(t, err)
		})
	}

	t.Run("authorised", func(t *testing.T) {
		a := auth.UserClaimsAuthoriser(model.User{})
		err := a.Authorise(&model.User{})
		assert.NoError(t, err)
	})
}
