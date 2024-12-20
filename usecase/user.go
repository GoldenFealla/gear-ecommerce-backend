package usecase

import (
	"context"
	"fmt"

	"github.com/goldenfealla/gear-manager/domain"
	"github.com/goldenfealla/gear-manager/internal/password"
	"github.com/google/uuid"
)

type UserRepository interface {
	CheckIDExist(ctx context.Context, id string) (bool, error)
	CheckEmailExist(ctx context.Context, email string) (bool, error)
	CheckUsernameExist(ctx context.Context, username string) (bool, error)
	CheckUsernameOrEmailExist(ctx context.Context, usernameOrEmail string) (bool, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*domain.User, error)
	AddUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, id string, user *domain.UpdateUserForm) error
}

type UserUsecase struct {
	r UserRepository
}

func NewUserUsecase(r UserRepository) *UserUsecase {
	return &UserUsecase{
		r,
	}
}

func (u *UserUsecase) RegisterUser(ctx context.Context, f *domain.RegisterUserForm) (*domain.UserInfo, error) {
	existedUsername, err := u.r.CheckUsernameExist(ctx, f.Username)

	if err != nil {
		return nil, err
	}

	if existedUsername {
		return nil, fmt.Errorf("username %v has already been used", f.Username)
	}

	existedEmail, err := u.r.CheckEmailExist(ctx, f.Email)

	if err != nil {
		return nil, err
	}

	if existedEmail {
		return nil, fmt.Errorf("email %v has already been used", f.Email)
	}

	hashedPassword, err := password.Generate(f.Password)
	if err != nil {
		return nil, fmt.Errorf("error while hashing password")
	}

	user := &domain.User{
		ID:        uuid.New(),
		Username:  f.Username,
		Email:     f.Email,
		FirstName: f.FirstName,
		LastName:  f.LastName,
		Phone:     f.Phone,
		Password:  hashedPassword,
	}

	err = u.r.AddUser(ctx, user)

	if err != nil {
		return nil, fmt.Errorf("error while creating user. Detail: %v", err.Error())
	}

	return &domain.UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
	}, nil
}

func (u *UserUsecase) LoginUser(ctx context.Context, f *domain.LoginUserForm) (*domain.UserInfo, error) {
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
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
	}, nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id string, f *domain.UpdateUserForm) (*domain.UserInfo, error) {
	existedUser, err := u.r.CheckIDExist(ctx, id)

	if err != nil {
		return nil, err
	}

	if !existedUser {
		return nil, fmt.Errorf("user not existed")
	}

	err = u.r.UpdateUser(ctx, id, f)

	if err != nil {
		return nil, err
	}

	user, err := u.r.GetUserByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return &domain.UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
	}, nil
}
