package property

import (
	"fmt"
	"strings"

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
	query := `INSERT INTO properties (
			user_id, title, description, address, price, area, rooms_count, 
			bathrooms_count, property_type, deal_type, material_type, 
			gas, electricity, internet, sewerage, plumbing,
			renovation, floor, land_area, floors)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)RETURNING id`

	row := r.db.QueryRow(
		query,
		userId,
		property.Title,
		property.Description,
		property.Address,
		property.Price,
		property.Area,
		property.RoomsCount,
		property.BathroomsCount,
		property.PropertyType,
		property.DealType,
		property.Material,
		property.Gas,
		property.Electricity,
		property.Internet,
		property.Sewerage,
		property.Plumbing,
		property.Renovation,
		property.Floor,
		property.LandArea,
		property.Floors,
	)

	if err := row.Scan(&id); err != nil {
		return err
	}

	return nil
}

func (r *PropertyRepository) DeleteProperty(userId, id uuid.UUID) error {
	query := `DELETE FROM properties WHERE user_id = $1 and id = $2`
	_, err := r.db.Exec(query, userId, id)
	return err
}

func (r *PropertyRepository) UpdateProperty(userId, id uuid.UUID, input PropertyUpdate) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Address != nil {
		setValues = append(setValues, fmt.Sprintf("address=$%d", argId))
		args = append(args, *input.Address)
		argId++
	}

	if input.Price != nil {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argId))
		args = append(args, *input.Price)
		argId++
	}

	if input.Area != nil {
		setValues = append(setValues, fmt.Sprintf("area=$%d", argId))
		args = append(args, *input.Area)
		argId++
	}

	if input.RoomsCount != nil {
		setValues = append(setValues, fmt.Sprintf("rooms_count=$%d", argId))
		args = append(args, *input.RoomsCount)
		argId++
	}

	if input.BathroomsCount != nil {
		setValues = append(setValues, fmt.Sprintf("bathrooms_count=$%d", argId))
		args = append(args, *input.BathroomsCount)
		argId++
	}

	if input.PropertyType != nil {
		setValues = append(setValues, fmt.Sprintf("property_type=$%d", argId))
		args = append(args, *input.PropertyType)
		argId++
	}

	if input.DealType != nil {
		setValues = append(setValues, fmt.Sprintf("deal_type=$%d", argId))
		args = append(args, *input.DealType)
		argId++
	}

	if input.Material != nil {
		setValues = append(setValues, fmt.Sprintf("material_type=$%d", argId))
		args = append(args, *input.Material)
		argId++
	}

	if input.Gas != nil {
		setValues = append(setValues, fmt.Sprintf("gas=$%d", argId))
		args = append(args, *input.Gas)
		argId++
	}

	if input.Electricity != nil {
		setValues = append(setValues, fmt.Sprintf("electricity=$%d", argId))
		args = append(args, *input.Electricity)
		argId++
	}

	if input.Internet != nil {
		setValues = append(setValues, fmt.Sprintf("internet=$%d", argId))
		args = append(args, *input.Internet)
		argId++
	}

	if input.Sewerage != nil {
		setValues = append(setValues, fmt.Sprintf("sewerage=$%d", argId))
		args = append(args, *input.Sewerage)
		argId++
	}

	if input.Plumbing != nil {
		setValues = append(setValues, fmt.Sprintf("plumbing=$%d", argId))
		args = append(args, *input.Plumbing)
		argId++
	}

	if input.Renovation != nil {
		setValues = append(setValues, fmt.Sprintf("renovation=$%d", argId))
		args = append(args, *input.Renovation)
		argId++
	}

	if input.LandArea != nil {
		setValues = append(setValues, fmt.Sprintf("land_area=$%d", argId))
		args = append(args, *input.LandArea)
		argId++
	}

	if input.Floor != nil {
		setValues = append(setValues, fmt.Sprintf("floor=$%d", argId))
		args = append(args, *input.Floor)
		argId++
	}

	if input.Floors != nil {
		setValues = append(setValues, fmt.Sprintf("floors=$%d", argId))
		args = append(args, *input.Floors)
		argId++
	}

	if len(setValues) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE properties SET %s WHERE id=$%d AND user_id=$%d", setQuery, argId, argId+1)
	args = append(args, id, userId)
	_, err := r.db.Exec(query, args...)
	return err
}
