package user_domain

import "goauth/utils/apperrors"

func NewEmailAlreadyExist() *apperrors.Error {
	return &apperrors.Error{
		Type:    apperrors.Conflict,
		Message: "Email already exist",
	}
}
