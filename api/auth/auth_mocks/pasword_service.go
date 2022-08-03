package auth_mocks

import "github.com/stretchr/testify/mock"

type MockPasswordService struct {
	mock.Mock
}

func (s *MockPasswordService) CreateHash(password string) (hash string, err error) {
	ret := s.Called(password)
	var r0 string

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(string)
	}
	return r0, nil
}

func (s *MockPasswordService) ComparePasswordAndHash(password string, hash string) (match bool, err error) {
	ret := s.Called(password, hash)
	var r0 bool

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(bool)
	}
	return r0, nil
}
