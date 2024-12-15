package postgres

import (
	"context"
	"errors"

	"github.com/goldenfealla/gear-manager/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	Conn *pgxpool.Pool
}

func NewOrderRepository(conn *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{Conn: conn}
}

func (r *OrderRepository) HasCart(ctx context.Context, userID string) bool {
	query := `SELECT EXISTS(SELECT 1 FROM "order" WHERE user_id=@user_id AND status=@status) `
	args := &pgx.NamedArgs{
		"user_id": userID,
		"status":  domain.CART,
	}

	var b bool
	r.Conn.QueryRow(ctx, query, args).Scan(&b)

	return b
}

func (r *OrderRepository) getOrderGearList(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderGear, error) {
	query := `
		SELECT Gear.*, OrderGear.quantity
		FROM "gear_order" OrderGear
		JOIN "gear" Gear ON OrderGear.gear_id=Gear.id
		WHERE order_id=@orderID
	`
	args := &pgx.NamedArgs{
		"orderID": orderID,
	}

	rows, err := r.Conn.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*domain.OrderGear, error) {
		var gear domain.Gear
		var quantity int64

		err := row.Scan(
			&gear.ID,
			&gear.Name,
			&gear.Type,
			&gear.Price,
			&gear.Discount,
			&gear.Quantity,
			&gear.ImageURL,
			&gear.Brand,
			&gear.Variety,
			&quantity,
		)

		if err != nil {
			return nil, err
		}

		return &domain.OrderGear{
			Gear:     &gear,
			Quantity: quantity,
		}, nil
	})
}

func (r *OrderRepository) GetFullCartByUserID(ctx context.Context, userID string) (*domain.FullOrder, error) {
	query := `SELECT * FROM "order" WHERE user_id=@user_id AND status=@status`
	args := &pgx.NamedArgs{
		"user_id": userID,
		"status":  domain.CART,
	}

	rows, err := r.Conn.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	order, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[domain.Order])
	if err != nil {
		return nil, err
	}

	orderGear, err := r.getOrderGearList(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	fullOrder := &domain.FullOrder{
		Order:     order,
		OrderGear: orderGear,
	}

	return fullOrder, nil
}

func (r *OrderRepository) GetFullCartByID(ctx context.Context, orderID string) (*domain.FullOrder, error) {
	query := `SELECT * FROM "order" WHERE id=@id AND status=@status`
	args := &pgx.NamedArgs{
		"id":     orderID,
		"status": domain.CART,
	}

	rows, err := r.Conn.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	order, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[domain.Order])
	if err != nil {
		return nil, err
	}

	orderGear, err := r.getOrderGearList(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	fullOrder := &domain.FullOrder{
		Order:     order,
		OrderGear: orderGear,
	}

	return fullOrder, nil
}

func (r *OrderRepository) GetCartInfo(ctx context.Context, userID string) (*domain.Order, error) {
	query := `SELECT * FROM "order" WHERE user_id=@user_id AND status=@status`
	args := &pgx.NamedArgs{
		"user_id": userID,
		"status":  domain.CART,
	}

	rows, err := r.Conn.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	order, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[domain.Order])
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) CreateCart(ctx context.Context, userID string) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid uuid")
	}

	cart := &domain.Order{
		ID:     id,
		Status: domain.CART,
		UserID: userUUID,
	}

	query := `
		INSERT INTO "order" (id, status, user_id, total)
		VALUES (@id, @status, @user_id, @total)
	`

	args := pgx.NamedArgs{
		"id":      cart.ID,
		"status":  cart.Status,
		"user_id": cart.UserID,
		"total":   0,
	}

	_, err = r.Conn.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) AddProductToCart(ctx context.Context, cart *domain.Order, gearID string) error {
	query := `
		INSERT INTO gear_order (order_id, gear_id, quantity)
		VALUES (@order_id, @gear_id, @quantity)
	`

	gearUUID, err := uuid.Parse(gearID)
	if err != nil {
		return errors.New("invalid gear uuid")
	}

	args := pgx.NamedArgs{
		"order_id": cart.ID,
		"gear_id":  gearUUID,
		"quantity": 1,
	}
	_, err = r.Conn.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) SetGearQuantityCart(ctx context.Context, cart *domain.Order, gearID string, quantity int64) error {
	query := `
		UPDATE gear_order
		SET quantity=@quantity
		WHERE order_id=@order_id AND gear_id=@gear_id
	`

	gearUUID, err := uuid.Parse(gearID)
	if err != nil {
		return errors.New("invalid gear uuid")
	}

	args := pgx.NamedArgs{
		"quantity": quantity,
		"order_id": cart.ID,
		"gear_id":  gearUUID,
	}
	_, err = r.Conn.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) RemoveProductToCart(ctx context.Context, cart *domain.Order, gearID string) error {
	query := `
		DELETE FROM gear_order
		WHERE order_id=@order_id AND gear_id=@gear_id
	`

	gearUUID, err := uuid.Parse(gearID)
	if err != nil {
		return errors.New("invalid gear uuid")
	}

	args := pgx.NamedArgs{
		"order_id": cart.ID,
		"gear_id":  gearUUID,
		"quantity": 1,
	}
	_, err = r.Conn.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status domain.OrderStatus) error {
	err := uuid.Validate(orderID)
	if err != nil {
		return errors.New("invalid uuid")
	}

	query := `
		UPDATE "order"
		SET status=@status
		WHERE id=@id
	`

	args := pgx.NamedArgs{
		"id":     orderID,
		"status": status,
	}

	_, err = r.Conn.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) UpdateOrderTotalPrice(ctx context.Context, orderID string, price int64) error {
	err := uuid.Validate(orderID)
	if err != nil {
		return errors.New("invalid uuid")
	}

	query := `
		UPDATE "order"
		SET total=@total
		WHERE id=@id
	`

	args := pgx.NamedArgs{
		"id":    orderID,
		"total": price,
	}

	_, err = r.Conn.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}
