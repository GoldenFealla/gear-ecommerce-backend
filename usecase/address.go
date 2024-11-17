package usecase

import (
	"context"
	"time"

	"github.com/goldenfealla/gear-manager/domain"
)

type AddressRepository interface {
	GetAddressByID(ctx context.Context, id string) (*domain.Address, error)
	GetAddressList(ctx context.Context, userID string) ([]*domain.Address, error)
	AddAddress(ctx context.Context, userID string, a *domain.AddAddressForm) error
	UpdateAddress(ctx context.Context, id string, a *domain.UpdateAddressForm) error
	DeleteAddress(ctx context.Context, id string) error
}

type AddressUsecase struct {
	r AddressRepository
}

func NewAddressUsecase(r AddressRepository) *AddressUsecase {
	return &AddressUsecase{
		r,
	}
}

func (u *AddressUsecase) GetAddressByID(id string) (*domain.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result, err := u.r.GetAddressByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return result, err
}

func (u *AddressUsecase) GetAddressList(userID string) ([]*domain.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result, err := u.r.GetAddressList(ctx, userID)

	if err != nil {
		return nil, err
	}

	return result, err
}

func (u *AddressUsecase) AddAddress(userID string, f *domain.AddAddressForm) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := u.r.AddAddress(ctx, userID, f)

	if err != nil {
		return err
	}

	return nil
}

func (u *AddressUsecase) UpdateAddress(id string, f *domain.UpdateAddressForm) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := u.r.UpdateAddress(ctx, id, f)

	if err != nil {
		return err
	}

	return nil
}

func (u *AddressUsecase) DeleteAddress(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := u.r.DeleteAddress(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
