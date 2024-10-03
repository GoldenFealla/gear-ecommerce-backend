package postgres

import (
	"context"

	"github.com/goldenfealla/gear-manager/domain"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	Conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) *UserRepository {
	return &UserRepository{Conn: conn}
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

func (r *UserRepository) GetByUsernameOrEmail(ctx context.Context, unoe string) (*domain.User, error) {
	query := `
		SELECT id, email, username, password, verified FROM "user" WHERE (email=@email OR username=@username)
	`
	args := &pgx.NamedArgs{
		"email":    unoe,
		"username": unoe,
	}

	var user domain.User
	err := r.Conn.QueryRow(ctx, query, args).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
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
		INSERT INTO "user" (id, username, email, password, verified)
		VALUES (@userID, @username, @userEmail, @userPassword, @userVerified)
	`

	args := &pgx.NamedArgs{
		"userID":       u.ID,
		"username":     u.Username,
		"userEmail":    u.Email,
		"userPassword": u.Password,
		"userVerified": false,
	}

	_, err := r.Conn.Exec(ctx, query, args)

	if err != nil {
		return err
	}

	return nil

}
