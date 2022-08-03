package user_service

import "goauth/user/user_domain"

type Service struct {
	UserRepository user_domain.IUserRepository
}

type ServiceConfig struct {
	UserRepository user_domain.IUserRepository
}

func NewUserService(c *ServiceConfig) user_domain.IUserService {
	return &Service{
		UserRepository: c.UserRepository,
	}
}
