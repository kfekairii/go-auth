package user_domain

type IUserService interface {
	Get(ID uint) (*User, error)
	GetByEmail(email string) (*User, error)
	Create(u *User) error
}

type IUserRepository interface {
	FindByID(ID uint) (*User, error)
	FindByEmail(email string) (*User, error)
	Create(u *User) error
}
