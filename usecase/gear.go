package usecase

import (
	"context"

	"github.com/goldenfealla/gear-manager/domain"
)

type GearRepository interface {
	GetGearVarietyList(ctx context.Context, category string) ([]string, error)
	GetGearBrandList(ctx context.Context, category string) ([]string, error)
	GetGearListCount(ctx context.Context, filter domain.ListGearFilter) (int64, error)
	GetGearList(ctx context.Context, filter domain.ListGearFilter) ([]*domain.Gear, error)
	GetGearByID(ctx context.Context, id string) (*domain.Gear, error)
	AddGear(ctx context.Context, g *domain.AddGearForm) error
	UpdateGear(ctx context.Context, id string, g *domain.UpdateGearForm) error
	UpdateGearQuantity(ctx context.Context, id string, quantity int64) error
	DeleteGear(ctx context.Context, id string) error
}

type GearUsecase struct {
	r GearRepository
}

func NewGearUsecase(r GearRepository) *GearUsecase {
	return &GearUsecase{
		r,
	}
}

func (u *GearUsecase) GetGearBrandList(ctx context.Context, category string) ([]string, error) {
	result, err := u.r.GetGearBrandList(ctx, category)

	if err != nil {
		return nil, err
	}

	return result, err
}

func (u *GearUsecase) GetGearVarietyList(ctx context.Context, category string) ([]string, error) {
	result, err := u.r.GetGearVarietyList(ctx, category)

	if err != nil {
		return nil, err
	}

	return result, err
}

func (u *GearUsecase) GetGearListCount(ctx context.Context, filter domain.ListGearFilter) (int64, error) {
	result, err := u.r.GetGearListCount(ctx, filter)

	if err != nil {
		return -1, err
	}

	return result, err
}

func (u *GearUsecase) GetGearList(ctx context.Context, filter domain.ListGearFilter) ([]*domain.Gear, error) {
	result, err := u.r.GetGearList(ctx, filter)

	if err != nil {
		return nil, err
	}

	return result, err
}

func (u *GearUsecase) GetGearByID(ctx context.Context, id string) (*domain.Gear, error) {
	result, err := u.r.GetGearByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return result, err
}

func (u *GearUsecase) AddGear(ctx context.Context, f *domain.AddGearForm) error {
	err := u.r.AddGear(ctx, f)

	if err != nil {
		return err
	}

	return nil
}

func (u *GearUsecase) UpdateGear(ctx context.Context, id string, f *domain.UpdateGearForm) error {
	err := u.r.UpdateGear(ctx, id, f)

	if err != nil {
		return err
	}

	return nil
}

func (u *GearUsecase) DeleteGear(ctx context.Context, id string) error {
	err := u.r.DeleteGear(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
