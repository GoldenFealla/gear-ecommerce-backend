package postgres

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/goldenfealla/gear-manager/domain"
	f "github.com/goldenfealla/gear-manager/internal/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GearRepository struct {
	Conn     *pgxpool.Pool
	S3Client *s3.Client
}

func NewGearRepository(conn *pgxpool.Pool, s3Client *s3.Client) *GearRepository {
	return &GearRepository{Conn: conn, S3Client: s3Client}
}

func (r *GearRepository) GetGearBrandList(ctx context.Context, category string) ([]string, error) {
	query := `SELECT DISTINCT brand FROM gear WHERE type=@type`

	key := strings.ToLower(category)

	if _, ok := domain.GearTypeMap[key]; !ok {
		return nil, errors.New("category not exist")
	}

	args := pgx.NamedArgs{
		"type": domain.GearTypeMap[key],
	}

	rows, _ := r.Conn.Query(ctx, query, args)

	type Brand struct {
		Brand string `db:"brand"`
	}

	brands, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Brand])

	if err != nil {
		return nil, err
	}

	n := len(brands)
	result := make([]string, n)

	for i := 0; i < n; i++ {
		result[i] = brands[i].Brand
	}

	return result, err
}

func (r *GearRepository) GetGearListCount(ctx context.Context, filter domain.ListGearFilter) (int64, error) {
	key := strings.ToLower(*filter.Category)

	if _, ok := domain.GearTypeMap[key]; !ok {
		return -1, errors.New("category not exist")
	}

	args := pgx.NamedArgs{}
	w := []string{}

	if key != "all" {
		args["category"] = domain.GearTypeMap[key]
		w = append(w, "type=@category")
	}

	if filter.Brand != nil {
		args["brand"] = *filter.Brand
		w = append(w, "brand=@brand")
	}

	if filter.StartPrice != nil && *filter.StartPrice != -1 {
		args["start_price"] = *filter.StartPrice
		w = append(w, "price>@start_price")
	}

	if filter.EndPrice != nil && *filter.EndPrice != -1 {
		args["end_price"] = *filter.EndPrice
		w = append(w, "price<@end_price")
	}

	query := fmt.Sprintf(
		`
            SELECT count(id) FROM gear 
            WHERE %v
		`,
		strings.Join(w, " AND "),
	)

	rows, _ := r.Conn.Query(ctx, query, args)

	type Count struct {
		Count int64 `db:"count"`
	}

	count, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Count])

	if err != nil {
		return -1, err
	}

	return count.Count, err

}

func (r *GearRepository) GetGearList(ctx context.Context, filter domain.ListGearFilter) ([]*domain.Gear, error) {
	key := strings.ToLower(*filter.Category)

	if _, ok := domain.GearTypeMap[key]; !ok {
		return nil, errors.New("category not exist")
	}

	args := pgx.NamedArgs{}
	w := []string{}

	if key != "all" {
		args["category"] = domain.GearTypeMap[key]
		w = append(w, "type=@category")
	}

	if filter.Brand != nil {
		args["brand"] = *filter.Brand
		w = append(w, "brand=@brand")
	}

	if filter.StartPrice != nil && *filter.StartPrice != -1 {
		args["start_price"] = *filter.StartPrice
		w = append(w, "price>@start_price")
	}

	if filter.EndPrice != nil && *filter.EndPrice != -1 {
		args["end_price"] = *filter.EndPrice
		w = append(w, "price<@end_price")
	}

	sort := ""

	if filter.Sort != nil {
		d := strings.ToUpper(*filter.Sort)
		sort = fmt.Sprintf("ORDER BY discount %v", d)
	}

	query := fmt.Sprintf(
		`
            SELECT * FROM gear 
            WHERE %v
            %v
            LIMIT %v OFFSET %v
		`,
		strings.Join(w, " AND "),
		sort,
		*filter.Limit,
		(*filter.Limit)*(*filter.Page-1),
	)

	rows, _ := r.Conn.Query(ctx, query, args)

	gears, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[domain.Gear])

	if err != nil {
		return nil, err
	}

	return gears, err
}

func (r *GearRepository) GetGearByID(ctx context.Context, id string) (*domain.Gear, error) {
	query := `
		SELECT * FROM gear WHERE id=@id
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
		INSERT INTO gear (id, name, type, price, discount, quantity, image_url, brand) 
		VALUES (@gearID, @gearName, @gearType, @gearPrice, @gearDiscount, @gearQuantity, @gearImageURL, @gearBrand)
	`

	newUUID, err := uuid.NewV7()

	if err != nil {
		return err
	}

	key := strings.ToLower(g.Type)

	args := pgx.NamedArgs{
		"gearID":       newUUID,
		"gearName":     g.Name,
		"gearType":     domain.GearTypeMap[key],
		"gearPrice":    g.Price,
		"gearDiscount": g.Discount,
		"gearQuantity": g.Quantity,
		"gearImageURL": "",
		"gearBrand":    g.Brand,
	}

	if g.ImageBase64 != nil {
		image_url, err := f.UploadImageJpeg(
			r.S3Client,
			*g.ImageBase64,
			fmt.Sprintf("%v.jpg", newUUID.String()),
		)

		if err != nil {
			return err
		}

		args["gearImageURL"] = *image_url
	}

	_, err = r.Conn.Exec(ctx, query, args)

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
			if field == "type" {
				val := domain.GearTypeMap[value.Elem().String()]
				args[field] = val
				fieldString = append(fieldString, fmt.Sprintf("%v='%v'", field, val))
			}

			if field != "type" {
				args[field] = value.Elem()
				fieldString = append(fieldString, fmt.Sprintf("%v='%v'", field, value.Elem()))
			}
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
