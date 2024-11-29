package usecase

import (
	"context"
	"time"

	"github.com/goldenfealla/gear-manager/domain"
)

type GearRepository interface {
	GetGearBrandList(ctx context.Context, category string) ([]string, error)
	GetGearList(ctx context.Context) ([]*domain.Gear, error)
	GetGearByID(ctx context.Context, id string) (*domain.Gear, error)
	AddGear(ctx context.Context, g *domain.AddGearForm) error
	UpdateGear(ctx context.Context, id string, g *domain.UpdateGearForm) error
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

func (u *GearUsecase) GetGearBrandList(category string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result, err := u.r.GetGearBrandList(ctx, category)

	if err != nil {
		return nil, err
	}

	return result, err
}

func (u *GearUsecase) GetGearList() ([]*domain.Gear, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result, err := u.r.GetGearList(ctx)

	if err != nil {
		return nil, err
	}

	return result, err
}

func (u *GearUsecase) GetGearByID(id string) (*domain.Gear, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result, err := u.r.GetGearByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return result, err
}

func (u *GearUsecase) AddGear(f *domain.AddGearForm) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := u.r.AddGear(ctx, f)

	if err != nil {
		return err
	}

	return nil
}

func (u *GearUsecase) UpdateGear(id string, f *domain.UpdateGearForm) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := u.r.UpdateGear(ctx, id, f)

	if err != nil {
		return err
	}

	return nil
}

func (u *GearUsecase) DeleteGear(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := u.r.DeleteGear(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
