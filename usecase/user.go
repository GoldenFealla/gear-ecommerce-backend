package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/goldenfealla/gear-manager/domain"
	"github.com/goldenfealla/gear-manager/internal/password"
	"github.com/google/uuid"
)

type UserRepository interface {
	CheckEmailExist(ctx context.Context, email string) (bool, error)
	CheckUsernameExist(ctx context.Context, username string) (bool, error)
	CheckUsernameOrEmailExist(ctx context.Context, usernameOrEmail string) (bool, error)
	GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*domain.User, error)
	AddUser(ctx context.Context, user *domain.User) error
}

type UserUsecase struct {
	r UserRepository
}

func NewUserUsecase(r UserRepository) *UserUsecase {
	return &UserUsecase{
		r,
	}
}

func (u *UserUsecase) RegisterUser(f *domain.RegisterUserForm) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	existedUsername, err := u.r.CheckUsernameExist(ctx, f.Username)

	if err != nil {
		return err
	}

	if existedUsername {
		return fmt.Errorf("username %v has already been used", f.Username)
	}

	existedEmail, err := u.r.CheckEmailExist(ctx, f.Email)

	if err != nil {
		return err
	}

	if existedEmail {
		return fmt.Errorf("email %v has already been used", f.Email)
	}

	hashedPassword, err := password.Generate(f.Password)
	if err != nil {
		return fmt.Errorf("Error while hashing password")
	}

	user := &domain.User{
		ID:       uuid.New(),
		Username: f.Username,
		Email:    f.Email,
		Password: hashedPassword,
	}

	err = u.r.AddUser(ctx, user)

	if err != nil {
		return fmt.Errorf("Error while creating user. Detail: %v", err.Error())
	}

	return nil
}

func (u *UserUsecase) LoginUser(f *domain.LoginUserForm) (*domain.UserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	existedUser, err := u.r.CheckUsernameOrEmailExist(ctx, f.UsernameOrEmail)

	if err != nil {
		return nil, err
	}

	if !existedUser {
		return nil, fmt.Errorf("user not existed")
	}

	user, err := u.r.GetByUsernameOrEmail(ctx, f.UsernameOrEmail)

	if err != nil {
		return nil, err
	}

	err = password.Compare(user.Password, f.Password)

	if err != nil {
		return nil, fmt.Errorf("incorrect password")
	}

	return &domain.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
