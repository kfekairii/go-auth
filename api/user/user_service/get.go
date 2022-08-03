package user_service

import (
	"goauth/user/user_domain"
)

func (s *Service) Get(ID uint) (*user_domain.User, error) {
	fetchedUser, err := s.UserRepository.FindByID(ID)

	return fetchedUser, err
}

func (s *Service) GetByEmail(email string) (*user_domain.User, error) {
	fetchedUser, err := s.UserRepository.FindByEmail(email)
	return fetchedUser, err
}
