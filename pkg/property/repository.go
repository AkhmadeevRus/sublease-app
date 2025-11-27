package property

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IPropertyRepository interface {
	CreateProperty(property Property, userId uuid.UUID) (uuid.UUID, error)
	GetMyProperties(id uuid.UUID) ([]Property, error)
	DeleteProperty(userId, id uuid.UUID) error
	UpdateProperty(userId, id uuid.UUID, input Property) error
}

type PropertyRepository struct {
	db *sqlx.DB
}

func NewPropertyRepository(db *sqlx.DB) *PropertyRepository {
	return &PropertyRepository{db: db}
}

func (r *PropertyRepository) CreateProperty(property Property, userId uuid.UUID) (uuid.UUID, error) {
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
		return uuid.Nil, err
	}

	return id, nil
}

func (r *PropertyRepository) GetMyProperties(id uuid.UUID) ([]Property, error) {
	return []Property{}, nil
}

func (r *PropertyRepository) DeleteProperty(userId, id uuid.UUID) error {
	return nil
}

func (r *PropertyRepository) UpdateProperty(userId, id uuid.UUID, input Property) error {
	return nil
}
