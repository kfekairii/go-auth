package user_mocks

import (
	"github.com/stretchr/testify/mock"
	"goauth/user/user_domain"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetByEmail(email string) (*user_domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserService) Get(ID uint) (*user_domain.User, error) {
	// Args that will be passed to "Return" in the tests
	ret := m.Called(ID)

	// First return value
	var r0 *user_domain.User

	// Get Returns the argument at the specified index.
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*user_domain.User)
	}
	// Second return value
	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockUserService) Create(u *user_domain.User) error {
	// Args that will be passed to "Return" in the tests
	ret := m.Called(u)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
