package auth_mocks

import (
	"github.com/stretchr/testify/mock"
	"goauth/auth/auth_domain"
	"goauth/user/user_domain"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) ValidateAccessToken(token string) (uint, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTokenService) NewPairFromUser(u *user_domain.User,
	prevAccessToken string) (*auth_domain.TokenPair, error) {

	ret := m.Called(u, prevAccessToken)

	var r0 *auth_domain.TokenPair
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*auth_domain.TokenPair)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
