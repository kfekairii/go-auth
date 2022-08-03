package user_mocks

import (
	"github.com/stretchr/testify/mock"
	"goauth/user/user_domain"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(email string) (*user_domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepository) FindByID(ID uint) (*user_domain.User, error) {
	ret := m.Called(ID)

	var r0 *user_domain.User
	var r1 error

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*user_domain.User)
	}

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1

}

func (m *MockUserRepository) Create(u *user_domain.User) error {
	ret := m.Called(u)

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r1

}
