package usecase

import (
	"context"

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

func (u *AddressUsecase) GetAddressByID(ctx context.Context, id string) (*domain.Address, error) {
	result, err := u.r.GetAddressByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return result, err
}

func (u *AddressUsecase) GetAddressList(ctx context.Context, userID string) ([]*domain.Address, error) {
	result, err := u.r.GetAddressList(ctx, userID)

	if err != nil {
		return nil, err
	}

	return result, err
}

func (u *AddressUsecase) AddAddress(ctx context.Context, userID string, f *domain.AddAddressForm) error {
	err := u.r.AddAddress(ctx, userID, f)

	if err != nil {
		return err
	}

	return nil
}

func (u *AddressUsecase) UpdateAddress(ctx context.Context, id string, f *domain.UpdateAddressForm) error {
	err := u.r.UpdateAddress(ctx, id, f)

	if err != nil {
		return err
	}

	return nil
}

func (u *AddressUsecase) DeleteAddress(ctx context.Context, id string) error {
	err := u.r.DeleteAddress(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
