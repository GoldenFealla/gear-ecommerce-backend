package usecase

import (
	"context"
	"errors"

	"github.com/goldenfealla/gear-manager/domain"
)

type OrderRepository interface {
	HasCart(ctx context.Context, userID string) bool
	GetFullCart(ctx context.Context, userID string) (*domain.FullOrder, error)
	GetCartInfo(ctx context.Context, userID string) (*domain.Order, error)
	CreateCart(ctx context.Context, userID string) error
	AddProductToCart(ctx context.Context, cart *domain.Order, gearID string) error
	SetGearQuantityCart(ctx context.Context, cart *domain.Order, gearID string, quantity int64) error
	RemoveProductToCart(ctx context.Context, cart *domain.Order, gearID string) error
}

type OrderUsercase struct {
	or OrderRepository
	ur UserRepository
}

func NewOrderUsercase(or OrderRepository, ur UserRepository) *OrderUsercase {
	return &OrderUsercase{
		or,
		ur,
	}
}

func (u *OrderUsercase) GetCart(ctx context.Context, userID string) (*domain.FullOrder, error) {
	isUserExisted, err := u.ur.CheckIDExist(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !isUserExisted {
		return nil, errors.New("user not existed")
	}

	if !u.or.HasCart(ctx, userID) {
		u.or.CreateCart(ctx, userID)
	}

	cart, err := u.or.GetFullCart(ctx, userID)

	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (u *OrderUsercase) AddGearToCart(ctx context.Context, userID string, gearID string) error {
	if !u.or.HasCart(ctx, userID) {
		u.or.CreateCart(ctx, userID)
	}

	cart, err := u.or.GetCartInfo(ctx, userID)
	if err != nil {
		return err
	}

	err = u.or.AddProductToCart(ctx, cart, gearID)
	if err != nil {
		return err
	}

	return nil
}

func (u *OrderUsercase) SetGearQuantityCart(ctx context.Context, userID string, gearID string, quantity int64) error {
	if quantity <= 0 {
		return errors.New("quantity must be bigger than 0")
	}

	cart, err := u.or.GetCartInfo(ctx, userID)
	if err != nil {
		return err
	}

	err = u.or.SetGearQuantityCart(ctx, cart, gearID, quantity)
	if err != nil {
		return err
	}

	return nil

}

func (u *OrderUsercase) RemoveGearFromCart(ctx context.Context, userID string, gearID string) error {
	if !u.or.HasCart(ctx, userID) {
		u.or.CreateCart(ctx, userID)
	}

	cart, err := u.or.GetCartInfo(ctx, userID)
	if err != nil {
		return err
	}

	err = u.or.RemoveProductToCart(ctx, cart, gearID)
	if err != nil {
		return err
	}

	return nil
}

func (u *OrderUsercase) GetOrder(ctx context.Context, id string) (*domain.FullOrder, error) {
	return nil, nil

}

func (u *OrderUsercase) GetOrderList(ctx context.Context, userID string) ([]*domain.FullOrder, error) {
	return nil, nil
}
