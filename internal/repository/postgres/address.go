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
	"github.com/jackc/pgx/v5/pgxpool"
)

type AddressRepository struct {
	Conn *pgxpool.Pool
}

func NewAddressRepository(conn *pgxpool.Pool) *AddressRepository {
	return &AddressRepository{Conn: conn}
}

func (r *AddressRepository) GetAddressByID(ctx context.Context, id string) (*domain.Address, error) {
	query := `
		SELECT * FROM address WHERE id=@id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	rows, _ := r.Conn.Query(ctx, query, args)

	address, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[domain.Address])

	if err != nil {
		return nil, err
	}

	return address, err
}

func (r *AddressRepository) GetAddressList(ctx context.Context, userID string) ([]*domain.Address, error) {
	query := `
		SELECT * FROM address WHERE user_id=@user_id;
	`

	args := pgx.NamedArgs{
		"user_id": userID,
	}

	rows, _ := r.Conn.Query(ctx, query, args)

	addresses, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[domain.Address])

	if err != nil {
		return nil, err
	}

	return addresses, err
}

func (r *AddressRepository) AddAddress(ctx context.Context, userID string, a *domain.AddAddressForm) error {
	query := `
		INSERT INTO address (id, address, country, user_id) 
		VALUES (@id, @address, @country, @user_id)
	`

	newUUID := uuid.New()

	args := pgx.NamedArgs{
		"id":      newUUID,
		"address": a.Address,
		"country": a.Country,
		"user_id": userID,
	}

	_, err := r.Conn.Exec(ctx, query, args)

	if err != nil {
		return err
	}

	return nil
}

func (r *AddressRepository) UpdateAddress(ctx context.Context, id string, a *domain.UpdateAddressForm) error {
	args := pgx.NamedArgs{
		"id": id,
	}

	// iterate through struct to get field need to update
	v := reflect.ValueOf(*a)
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
		UPDATE address
		SET %v
		WHERE id=@id;
	`, strings.Join(fieldString, ","))

	_, err := r.Conn.Exec(ctx, query, args)

	if err != nil {
		return err
	}

	return nil
}

func (r *AddressRepository) DeleteAddress(ctx context.Context, id string) error {
	query := `
		DELETE FROM address
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
