package user_service

import (
	"goauth/user/user_domain"
)

func (s *Service) Create(u *user_domain.User) error {
	if err := s.UserRepository.Create(u); err != nil {

		return err
	}
	return nil
}
