package domain

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Email    string    `json:"email" db:"email"`
	Password string    `json:"password" db:"password"`
	Verified bool      `json:"verified" db:"verified"`
}

type UserCredential struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type UserInfo struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Email    string    `json:"email" db:"email"`
}

type RegisterUserForm struct {
	Username string `json:"username" validate:"required,gte=6,lte=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8,lte=24"`
}

type LoginUserForm struct {
	UsernameOrEmail string `json:"username_or_email" validate:"required"`
	Password        string `json:"password" validate:"required,gte=8,lte=24"`
}

type Profile struct {
	ID        uuid.UUID `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Phone     string    `json:"phone" db:"phone"`
}
