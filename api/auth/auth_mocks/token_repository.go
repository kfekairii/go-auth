package auth_mocks

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type MockTokenRepository struct {
	mock.Mock
}

func (r *MockTokenRepository) SetRefreshToken(userID string, tokenID string, expiresIn time.Duration) error {
	ret := r.Called(userID, tokenID, expiresIn)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)

	}

	return r0
}

func (r *MockTokenRepository) DeleteRefreshToken(userID string, prevTokenID string) error {
	ret := r.Called(userID, prevTokenID)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)

	}

	return r0
}
