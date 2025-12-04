package search

import (
	"fmt"
	"strings"

	"github.com/AkhmadeevRus/sublease-app/pkg/property"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ISearchRepository interface {
	GetAllProperties() ([]property.Property, error)
	GetPropertyById(id uuid.UUID) (property.Property, error)
	SearchByFilters(filter PropertyFilter) ([]property.Property, error)
}

type SearchRepository struct {
	db *sqlx.DB
}

func NewSearchRepository(db *sqlx.DB) *SearchRepository {
	return &SearchRepository{db: db}
}

func (r *SearchRepository) GetAllProperties() ([]property.Property, error) {
	var properties []property.Property
	query := `SELECT id, user_id, title, description, address, price, area, rooms_count, 
			bathrooms_count, property_type, deal_type, material_type, 
			gas, electricity, internet, sewerage, plumbing,
			renovation, floor, land_area, floors
			FROM properties`
	err := r.db.Select(&properties, query)
	if err != nil {
		return nil, err
	}
	return properties, nil
}

func (r *SearchRepository) GetPropertyById(id uuid.UUID) (property.Property, error) {
	var outProperty property.Property
	query := `SELECT id, user_id, title, description, address, price, area, rooms_count, 
			bathrooms_count, property_type, deal_type, material_type, 
			gas, electricity, internet, sewerage, plumbing,
			renovation, floor, land_area, floors
			FROM properties
			WHERE id = $1`
	err := r.db.Get(&outProperty, query, id)
	if err != nil {
		return property.Property{}, err
	}
	return outProperty, nil
}

func (r *SearchRepository) SearchByFilters(filter PropertyFilter) ([]property.Property, error) {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argsId := 1

	if filter.MinPrice != nil {
		setValues = append(setValues, fmt.Sprintf("price >= $%d", argsId))
		args = append(args, *filter.MinPrice)
		argsId++
	}
	if filter.MaxPrice != nil {
		setValues = append(setValues, fmt.Sprintf("price <= $%d", argsId))
		args = append(args, *filter.MaxPrice)
		argsId++
	}
	if filter.Area != nil {
		setValues = append(setValues, fmt.Sprintf("area >= $%d", argsId))
		args = append(args, *filter.Area)
		argsId++
	}
	if filter.RoomsCount != nil {
		setValues = append(setValues, fmt.Sprintf("rooms_count = $%d", argsId))
		args = append(args, *filter.RoomsCount)
		argsId++
	}
	if filter.BathroomsCount != nil {
		setValues = append(setValues, fmt.Sprintf("bathrooms_count = $%d", argsId))
		args = append(args, *filter.BathroomsCount)
		argsId++
	}
	if filter.PropertyType != nil {
		setValues = append(setValues, fmt.Sprintf("property_type = $%d", argsId))
		args = append(args, *filter.PropertyType)
		argsId++
	}
	if filter.DealType != nil {
		setValues = append(setValues, fmt.Sprintf("deal_type = $%d", argsId))
		args = append(args, *filter.DealType)
		argsId++
	}
	if filter.Material != nil {
		setValues = append(setValues, fmt.Sprintf("material_type = $%d", argsId))
		args = append(args, *filter.Material)
		argsId++
	}
	if filter.Gas != nil {
		setValues = append(setValues, fmt.Sprintf("gas = $%d", argsId))
		args = append(args, *filter.Gas)
		argsId++
	}
	if filter.Electricity != nil {
		setValues = append(setValues, fmt.Sprintf("electricity = $%d", argsId))
		args = append(args, *filter.Electricity)
		argsId++
	}
	if filter.Internet != nil {
		setValues = append(setValues, fmt.Sprintf("internet = $%d", argsId))
		args = append(args, *filter.Internet)
		argsId++
	}
	if filter.Sewerage != nil {
		setValues = append(setValues, fmt.Sprintf("sewerage = $%d", argsId))
		args = append(args, *filter.Sewerage)
		argsId++
	}
	if filter.Plumbing != nil {
		setValues = append(setValues, fmt.Sprintf("plumbing = $%d", argsId))
		args = append(args, *filter.Plumbing)
		argsId++
	}
	if filter.Renovation != nil {
		setValues = append(setValues, fmt.Sprintf("renovation = $%d", argsId))
		args = append(args, *filter.Renovation)
		argsId++
	}
	if filter.Floor != nil {
		setValues = append(setValues, fmt.Sprintf("floor = $%d", argsId))
		args = append(args, *filter.Floor)
		argsId++
	}
	if filter.LandArea != nil {
		setValues = append(setValues, fmt.Sprintf("land_area >= $%d", argsId))
		args = append(args, *filter.LandArea)
		argsId++
	}
	if filter.Floors != nil {
		setValues = append(setValues, fmt.Sprintf("floors = $%d", argsId))
		args = append(args, *filter.Floors)
		argsId++
	}
	query := `SELECT id, user_id, title, description, address, price, area, rooms_count, 
			bathrooms_count, property_type, deal_type, material_type, 
			gas, electricity, internet, sewerage, plumbing,
			renovation, floor, land_area, floors
			FROM properties`

	if len(setValues) > 0 {
		query += " WHERE " + strings.Join(setValues, " AND ")
	}
	var properties []property.Property
	err := r.db.Select(&properties, query, args...)
	if err != nil {
		return nil, err
	}

	return properties, nil
}
