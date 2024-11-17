package domain

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Phone     string    `json:"phone" db:"phone"`
	Password  string    `json:"password" db:"password"`
	Verified  bool      `json:"verified" db:"verified"`
}

type UserInfo struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Phone     string    `json:"phone" db:"phone"`
}

type UserCredential struct {
	Token    string    `json:"token"`
	UserInfo *UserInfo `json:"user"`
}

type RegisterUserForm struct {
	Username  string `json:"username" validate:"required,gte=6,lte=20"`
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" db:"first_name" validate:"required,gte=2,lte=30"`
	LastName  string `json:"last_name" db:"last_name" validate:"required,gte=2,lte=30"`
	Phone     string `json:"phone" db:"phone" validate:"required,gte=0"`
	Password  string `json:"password" validate:"required,gte=8,lte=24"`
}

type LoginUserForm struct {
	UsernameOrEmail string `json:"username_or_email" validate:"required"`
	Password        string `json:"password" validate:"required,gte=8,lte=24"`
}

type UpdateUserForm struct {
	Username  *string `json:"username,omitempty" db:"username" validate:"omitempty,gte=6,lte=20"`
	Email     *string `json:"email,omitempty" db:"email" validate:"omitempty,email"`
	FirstName *string `json:"first_name,omitempty" db:"first_name" validate:"omitempty,gte=2,lte=30"`
	LastName  *string `json:"last_name,omitempty" db:"last_name" validate:"omitempty,gte=2,lte=30"`
	Phone     *string `json:"phone,omitempty" db:"phone"`
}
