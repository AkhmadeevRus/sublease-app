package property

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IPropertyRepository interface {
	CreateProperty(property Property, userId uuid.UUID) error
	DeleteProperty(userId, id uuid.UUID) error
	UpdateProperty(userId, id uuid.UUID, input PropertyUpdate) error
}

type PropertyRepository struct {
	db *sqlx.DB
}

func NewPropertyRepository(db *sqlx.DB) *PropertyRepository {
	return &PropertyRepository{db: db}
}

func (r *PropertyRepository) CreateProperty(property Property, userId uuid.UUID) error {
	var id uuid.UUID
	sql, args, err := sq.Insert("properties").
		Columns("user_id", "title", "description", "address", "price", "area",
			"rooms_count", "bathrooms_count", "property_type", "deal_type", "material_type",
			"gas", "electricity", "internet", "sewerage", "plumbing",
			"renovation", "floor", "land_area", "floors").
		Values(userId, property.Title,
			property.Description, property.Address, property.Price,
			property.Area, property.RoomsCount, property.BathroomsCount,
			property.PropertyType, property.DealType, property.Material,
			property.Gas, property.Electricity, property.Internet,
			property.Sewerage, property.Plumbing, property.Renovation,
			property.Floor, property.LandArea, property.Floors).
		Suffix("RETURNING id").ToSql()

	if err != nil {
		return err
	}

	err = r.db.QueryRow(sql, args...).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to exec query: %w", err)
	}

	return nil
}

func (r *PropertyRepository) DeleteProperty(userId, id uuid.UUID) error {
	sql, args, err := sq.Delete("properties").
		Where(sq.Eq{"id": id, "user_id": userId}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *PropertyRepository) UpdateProperty(userId, id uuid.UUID, input PropertyUpdate) error {
	query := sq.Update("properties")
	if input.Title != nil {
		query = query.Set("title", *input.Title)
	}
	if input.Description != nil {
		query = query.Set("description", *input.Description)
	}
	if input.Address != nil {
		query = query.Set("address", *input.Address)
	}
	if input.Price != nil {
		query = query.Set("price", *input.Price)
	}
	if input.Area != nil {
		query = query.Set("area", *input.Area)
	}
	if input.RoomsCount != nil {
		query = query.Set("rooms_count", *input.RoomsCount)
	}
	if input.BathroomsCount != nil {
		query = query.Set("bathrooms_count", *input.BathroomsCount)
	}
	if input.PropertyType != nil {
		query = query.Set("property_type", *input.PropertyType)
	}
	if input.DealType != nil {
		query = query.Set("deal_type", *input.DealType)
	}
	if input.Material != nil {
		query = query.Set("material_type", *input.Material)
	}
	if input.Gas != nil {
		query = query.Set("gas", *input.Gas)
	}
	if input.Electricity != nil {
		query = query.Set("electricity", *input.Electricity)
	}
	if input.Internet != nil {
		query = query.Set("internet", *input.Internet)
	}
	if input.Sewerage != nil {
		query = query.Set("sewerage", *input.Sewerage)
	}
	if input.Plumbing != nil {
		query = query.Set("plumbing", *input.Plumbing)
	}
	if input.Renovation != nil {
		query = query.Set("renovation", *input.Renovation)
	}
	if input.Floor != nil {
		query = query.Set("floor", *input.Floor)
	}
	if input.LandArea != nil {
		query = query.Set("land_area", *input.LandArea)
	}
	if input.Floors != nil {
		query = query.Set("floors", *input.Floors)
	}
	sql, args, err := query.Where(sq.Eq{"id": id, "user_id": userId}).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(sql, args...)
	if err != nil {
		return err
	}

	return nil
}
