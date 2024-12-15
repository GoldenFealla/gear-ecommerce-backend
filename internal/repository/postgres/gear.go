package postgres

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
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

func (r *GearRepository) getGearFilterList(ctx context.Context, category string, field string) ([]string, error) {
	query := fmt.Sprintf("SELECT DISTINCT %v FROM gear WHERE type=@type", field)

	key := strings.ToLower(category)

	if _, ok := domain.GearTypeMap[key]; !ok {
		return nil, errors.New("category not exist")
	}

	args := pgx.NamedArgs{
		"type": domain.GearTypeMap[key],
	}

	rows, _ := r.Conn.Query(ctx, query, args)

	list, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (string, error) {
		var item string
		err := row.Scan(&item)
		return item, err
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (r *GearRepository) GetGearBrandList(ctx context.Context, category string) ([]string, error) {
	return r.getGearFilterList(ctx, category, "brand")
}

func (r *GearRepository) GetGearVarietyList(ctx context.Context, category string) ([]string, error) {
	return r.getGearFilterList(ctx, category, "variety")
}

func (r *GearRepository) processWhereFilter(args pgx.NamedArgs, filter domain.ListGearFilter) (*string, error) {
	key := strings.ToLower(*filter.Category)

	if _, ok := domain.GearTypeMap[key]; !ok {
		return nil, errors.New("category not exist")
	}

	w := []string{}

	if key != "all" {
		args["category"] = domain.GearTypeMap[key]
		w = append(w, "type=@category")
	}

	if filter.Brand != nil {
		args["brand"] = *filter.Brand
		w = append(w, "brand=@brand")
	}

	if filter.Variety != nil {
		args["variety"] = *filter.Variety
		w = append(w, "variety=@variety")
	}

	if filter.Price != nil {
		prices := strings.Split(*filter.Price, ",")

		startPrice, err := strconv.ParseInt(prices[0], 10, 64)

		if err != nil {
			return nil, err
		}

		if startPrice != -1 {
			args["start_price"] = startPrice
			w = append(w, "price>@start_price")
		}

		endPrice, err := strconv.ParseInt(prices[1], 10, 64)

		if err != nil {
			return nil, err
		}

		if endPrice != -1 {
			args["end_price"] = endPrice
			w = append(w, "price<@end_price")
		}
	}

	where := ""
	if len(w) > 0 {
		where = fmt.Sprintf("WHERE %v", strings.Join(w, " AND "))
	}

	return &where, nil
}

func (r *GearRepository) GetGearListCount(ctx context.Context, filter domain.ListGearFilter) (int64, error) {
	args := pgx.NamedArgs{}

	where, err := r.processWhereFilter(args, filter)

	if err != nil {
		return -1, err
	}

	query := fmt.Sprintf(
		`
            SELECT count(id) FROM gear 
            %v
		`,
		*where,
	)

	rows, _ := r.Conn.Query(ctx, query, args)

	count, err := pgx.CollectOneRow(rows, func(row pgx.CollectableRow) (int64, error) {
		var count int64
		err := row.Scan((&count))
		return count, err
	})

	if err != nil {
		return -1, err
	}

	return count, err
}

func (r *GearRepository) GetGearList(ctx context.Context, filter domain.ListGearFilter) ([]*domain.Gear, error) {
	args := pgx.NamedArgs{}

	where, err := r.processWhereFilter(args, filter)

	if err != nil {
		return nil, err
	}

	sort := ""

	if filter.Sort != nil {
		d := strings.ToUpper(*filter.Sort)
		sort = fmt.Sprintf("ORDER BY discount %v", d)
	}

	query := fmt.Sprintf(
		`
            SELECT * FROM gear 
            %v
            %v
            LIMIT %v OFFSET %v
		`,
		*where,
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
		INSERT INTO gear (id, name, type, price, discount, quantity, image_url, brand, variety) 
		VALUES (@gearID, @gearName, @gearType, @gearPrice, @gearDiscount, @gearQuantity, @gearImageURL, @gearBrand, @gearVariety)
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
		"gearVariety":  g.Variety,
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
			if field == "image_base64" {
				continue
			}

			if field == "type" {
				key := strings.ToLower(value.Elem().String())
				val := domain.GearTypeMap[key]
				args[field] = val
				fieldString = append(fieldString, fmt.Sprintf("%v='%v'", field, val))
			}

			if field != "type" {
				args[field] = value.Elem()
				fieldString = append(fieldString, fmt.Sprintf("%v='%v'", field, value.Elem()))
			}
		}
	}

	if g.ImageBase64 != nil {
		imageURL, err := f.UploadImageJpeg(
			r.S3Client,
			*g.ImageBase64,
			fmt.Sprintf("%v.jpg", id),
		)

		if err != nil {
			return err
		}

		fieldString = append(fieldString, fmt.Sprintf("%v='%v'", "image_url", *imageURL))
		args["gearImageURL"] = *imageURL
	}

	if len(fieldString) == 0 {
		return errors.New("field to update is required")
	}

	query := fmt.Sprintf(`
		UPDATE gear
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

func (r *GearRepository) UpdateGearQuantity(ctx context.Context, gearID string, quantity int64) error {
	query := `
		UPDATE gear
		SET quantity=@quantity
		WHERE @id=id
	`

	args := pgx.NamedArgs{
		"id":       gearID,
		"quantity": quantity,
	}

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
