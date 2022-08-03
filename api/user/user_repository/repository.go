package user_repository

import (
	"errors"
	"github.com/jackc/pgconn"
	"goauth/user/user_domain"
	"goauth/utils/apperrors"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) user_domain.IUserRepository {
	return &UserRepository{
		DB: db,
	}
}
func (r UserRepository) Create(u *user_domain.User) error {
	result := r.DB.Create(u)
	var perr *pgconn.PgError
	if errors.As(result.Error, &perr) {

		if perr.Code == "23505" {
			return user_domain.NewEmailAlreadyExist()
		}

		return apperrors.NewInternal()
	}
	return result.Error
}

func (r UserRepository) FindByID(ID uint) (*user_domain.User, error) {
	var user user_domain.User
	result := r.DB.First(&user, ID)
	if result.Error != nil {

		return nil, result.Error
	}
	return &user, nil
}

func (r UserRepository) FindByEmail(email string) (*user_domain.User, error) {
	var user user_domain.User
	result := r.DB.Where(&user_domain.User{Email: email}).First(&user)
	if result.Error != nil {

		return nil, result.Error
	}
	return &user, nil
}
