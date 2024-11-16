package postgres

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/goldenfealla/gear-manager/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type GearRepository struct {
	Conn *pgx.Conn
}

func NewGearRepository(conn *pgx.Conn) *GearRepository {
	return &GearRepository{Conn: conn}
}

func (r *GearRepository) GetGearList(ctx context.Context) ([]*domain.Gear, error) {
	rows, _ := r.Conn.Query(ctx, "SELECT id, name, type, price, discount, quantity FROM gear")

	gears, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[domain.Gear])

	if err != nil {
		return nil, err
	}

	return gears, err
}

func (r *GearRepository) GetGearByID(ctx context.Context, id string) (*domain.Gear, error) {
	query := `
		SELECT id, name, type, price, discount, quantity FROM gear WHERE id=@id
	`
	args := pgx.NamedArgs{
		"id": id,
	}

	rows, _ := r.Conn.Query(ctx, query, args)

	gear, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[domain.Gear])

	if err != nil {
		return nil, err
	}

	return gear, err
}

func (r *GearRepository) AddGear(ctx context.Context, g *domain.AddGearForm) error {
	query := `
		INSERT INTO gear (id, name, type, price, discount, quantity) 
		VALUES (@gearID, @gearName, @gearType, @gearPrice, @gearDiscount, @gearQuantity)
	`

	newUUID := uuid.New()

	args := pgx.NamedArgs{
		"gearID":       newUUID,
		"gearName":     g.Name,
		"gearType":     g.Type,
		"gearPrice":    g.Price,
		"gearDiscount": g.Discount,
		"gearQuantity": g.Quantity,
	}

	_, err := r.Conn.Exec(ctx, query, args)

	if err != nil {
		return err
	}

	return nil
}

func (r *GearRepository) UpdateGear(ctx context.Context, id string, g *domain.UpdateGearForm) error {
	args := pgx.NamedArgs{
		"id": id,
	}

	// iterate through struct to get field need to update
	v := reflect.ValueOf(*g)
	typeOfG := v.Type()

	fieldString := []string{}

	for i := 0; i < v.NumField(); i++ {
		field := typeOfG.Field(i).Tag.Get("db")
		value := v.Field(i)

		if !value.IsNil() {
			args[field] = value.Elem()
			fieldString = append(fieldString, fmt.Sprintf("%v='%v'", field, value.Elem()))
		}
	}

	if len(fieldString) == 0 {
		return errors.New("field to update is required")
	}

	query := fmt.Sprintf(`
		UPDATE gear
		SET %v
		WHERE id=@id;
	`, strings.Join(fieldString, ","))

	_, err := r.Conn.Exec(ctx, query, args)

	if err != nil {
		return err
	}

	return nil
}

func (r *GearRepository) DeleteGear(ctx context.Context, id string) error {
	query := `
		DELETE FROM gear
		WHERE id=@id
	`
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := r.Conn.Exec(ctx, query, args)

	if err != nil {
		return err
	}

	return nil
}
