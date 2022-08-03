package user_service_test

import (
	"fmt"
	"goauth/user/user_domain"
	"goauth/user/user_mocks"
	"goauth/user/user_service"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGet(t *testing.T) {
	const USERID uint = 101
	t.Run("Success", func(t *testing.T) {
		mockResp := &user_domain.User{
			Model: gorm.Model{
				ID: USERID,
			},
			Email: "dev.kfek@gmail.com",
		}
		mockUserRepository := new(user_mocks.MockUserRepository)
		mockUserRepository.On("FindByID", USERID).
			Return(mockResp, nil)

		us := user_service.NewUserService(&user_service.ServiceConfig{
			UserRepository: mockUserRepository,
		})

		u, err := us.Get(USERID)
		assert.NoError(t, err)
		assert.Equal(t, u, mockResp)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockUserRepository := new(user_mocks.MockUserRepository)
		mockUserRepository.On("FindByID", USERID).
			Return(nil, fmt.Errorf("error"))

		us := user_service.NewUserService(&user_service.ServiceConfig{
			UserRepository: mockUserRepository,
		})

		u, err := us.Get(USERID)
		assert.Nil(t, u)
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})
}
