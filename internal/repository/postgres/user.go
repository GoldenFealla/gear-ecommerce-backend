package postgres

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/goldenfealla/gear-manager/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	Conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) *UserRepository {
	return &UserRepository{Conn: conn}
}

func (r *UserRepository) CheckIDExist(ctx context.Context, id string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM "user" WHERE id=@id)
	`
	args := &pgx.NamedArgs{
		"id": id,
	}

	var b bool
	err := r.Conn.QueryRow(ctx, query, args).Scan(&b)

	if err != nil {
		return false, err
	}

	return b, nil
}

func (r *UserRepository) CheckEmailExist(ctx context.Context, e string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM "user" WHERE email=@email)
	`
	args := &pgx.NamedArgs{
		"email": e,
	}

	var b bool
	err := r.Conn.QueryRow(ctx, query, args).Scan(&b)

	if err != nil {
		return false, err
	}

	return b, nil
}

func (r *UserRepository) CheckUsernameExist(ctx context.Context, un string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM "user" WHERE username=@username)
	`
	args := &pgx.NamedArgs{
		"username": un,
	}

	var b bool
	err := r.Conn.QueryRow(ctx, query, args).Scan(&b)

	if err != nil {
		return false, err
	}

	return b, nil
}

func (r *UserRepository) CheckUsernameOrEmailExist(ctx context.Context, unoe string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM "user" WHERE (email=@email OR username=@username))
	`
	args := &pgx.NamedArgs{
		"email":    unoe,
		"username": unoe,
	}

	var b bool
	err := r.Conn.QueryRow(ctx, query, args).Scan(&b)

	if err != nil {
		return false, err
	}

	return b, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, username, email, first_name, last_name, phone, password, verified FROM "user" WHERE id=@id
	`
	args := &pgx.NamedArgs{
		"id": id,
	}

	var user domain.User
	err := r.Conn.QueryRow(ctx, query, args).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Password,
		&user.Verified,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByUsernameOrEmail(ctx context.Context, unoe string) (*domain.User, error) {
	query := `
		SELECT id, username, email, first_name, last_name, phone, password, verified FROM "user" WHERE (email=@email OR username=@username)
	`
	args := &pgx.NamedArgs{
		"email":    unoe,
		"username": unoe,
	}

	var user domain.User
	err := r.Conn.QueryRow(ctx, query, args).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Password,
		&user.Verified,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) AddUser(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO "user" 
			(
				id, 
				username, 
				email, 
				first_name, 
				last_name, 
				phone, 
				password, 
				verified
			)
		VALUES 
			(
				@userID, 
				@username, 
				@userEmail, 
				@userFirstName, 
				@userLastName,
				@userPhone,
				@userPassword, 
				@userVerified
			)
	`

	args := &pgx.NamedArgs{
		"userID":        u.ID,
		"username":      u.Username,
		"userEmail":     u.Email,
		"userFirstName": u.FirstName,
		"userLastName":  u.LastName,
		"userPhone":     u.Phone,
		"userPassword":  u.Password,
		"userVerified":  false,
	}

	_, err := r.Conn.Exec(ctx, query, args)

	if err != nil {
		return err
	}

	return nil

}

func (r *UserRepository) UpdateUser(ctx context.Context, id string, u *domain.UpdateUserForm) error {
	args := pgx.NamedArgs{
		"id": id,
	}

	// iterate through struct to get field need to update
	v := reflect.ValueOf(*u)
	typeOfG := v.Type()

	fieldString := []string{}

	for i := 0; i < v.NumField(); i++ {
		field := typeOfG.Field(i).Tag.Get("db")
		value := v.Field(i)

		if !value.IsNil() {
			fieldString = append(fieldString, fmt.Sprintf("%v='%v'", field, value.Elem()))
		}
	}

	if len(fieldString) == 0 {
		return errors.New("field to update is required")
	}

	query := fmt.Sprintf(`
		UPDATE "user" 
		SET %v 
		WHERE id=@id;
	`, strings.Join(fieldString, ","))

	fmt.Println(query)

	_, err := r.Conn.Exec(ctx, query, args)

	if err != nil {
		return err
	}

	return nil
}
