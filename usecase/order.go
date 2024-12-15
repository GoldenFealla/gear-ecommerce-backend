package usecase

import (
	"context"
	"errors"

	"github.com/goldenfealla/gear-manager/domain"
)

type OrderRepository interface {
	HasCart(ctx context.Context, userID string) bool
	GetFullCartByUserID(ctx context.Context, userID string) (*domain.FullOrder, error)
	GetFullOrderByID(ctx context.Context, orderID string) (*domain.FullOrder, error)
	GetFullOrderList(ctx context.Context, userID string, page int64, limit int64) ([]*domain.Order, error)
	GetCartInfo(ctx context.Context, userID string) (*domain.Order, error)
	CreateCart(ctx context.Context, userID string) error
	AddProductToCart(ctx context.Context, cart *domain.Order, gearID string) error
	SetGearQuantityCart(ctx context.Context, cart *domain.Order, gearID string, quantity int64) error
	RemoveProductToCart(ctx context.Context, cart *domain.Order, gearID string) error
	UpdateOrderStatus(ctx context.Context, cartID string, status domain.OrderStatus) error
	UpdateOrderTotalPrice(ctx context.Context, cartID string, price int64) error
}

type OrderUsercase struct {
	or OrderRepository
	ur UserRepository
	gr GearRepository
}

func NewOrderUsercase(or OrderRepository, ur UserRepository, gr GearRepository) *OrderUsercase {
	return &OrderUsercase{
		or,
		ur,
		gr,
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

	cart, err := u.or.GetFullCartByUserID(ctx, userID)

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

func (u *OrderUsercase) updateGearQuantityWorker(
	ctx context.Context,
	ordergear <-chan *domain.OrderGear,
	pricegear chan<- int64,
) {
	for og := range ordergear {
		r := og.Gear.Quantity - og.Quantity
		u.gr.UpdateGearQuantity(ctx, og.Gear.ID.String(), r)
		pricegear <- int64(og.Gear.Price) * og.Quantity
	}
}

func (u *OrderUsercase) PayCart(ctx context.Context, orderID string) error {
	cart, err := u.or.GetFullOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	err = u.or.UpdateOrderStatus(ctx, cart.Order.ID.String(), domain.PAID)
	if err != nil {
		return err
	}

	numOfOrderGear := len(cart.OrderGear)
	orderGearsChan := make(chan *domain.OrderGear, numOfOrderGear)
	priceGearsChan := make(chan int64, numOfOrderGear)
	// 3 worker
	for w := 1; w <= 3; w++ {
		go u.updateGearQuantityWorker(ctx, orderGearsChan, priceGearsChan)
	}

	for j := 0; j < numOfOrderGear; j++ {
		orderGearsChan <- cart.OrderGear[j]
	}
	close(orderGearsChan)

	totalPrice := int64(0)

	for j := 0; j < numOfOrderGear; j++ {
		priceGear := <-priceGearsChan
		totalPrice += priceGear
	}

	u.or.UpdateOrderTotalPrice(ctx, cart.Order.ID.String(), totalPrice)

	return nil
}

func (u *OrderUsercase) GetOrder(ctx context.Context, id string) (*domain.FullOrder, error) {
	order, err := u.or.GetFullOrderByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (u *OrderUsercase) GetOrderList(ctx context.Context, userID string, page int64, limit int64) ([]*domain.Order, error) {
	orders, err := u.or.GetFullOrderList(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
