package search

import (
	"github.com/AkhmadeevRus/sublease-app/pkg/property"
	sq "github.com/Masterminds/squirrel"
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
	sql, args, err := sq.Select("id", "user_id", "title", "description", "address", "price", "area",
		"rooms_count", "bathrooms_count", "property_type", "deal_type", "material_type",
		"gas", "electricity", "internet", "sewerage", "plumbing",
		"renovation", "floor", "land_area", "floors").
		From("properties").
		ToSql()
	if err != nil {
		return nil, err
	}
	err = r.db.Select(&properties, sql, args)
	if err != nil {
		return nil, err
	}
	return properties, nil
}

func (r *SearchRepository) GetPropertyById(id uuid.UUID) (property.Property, error) {
	var outProperty property.Property
	sql, args, err := sq.Select("id", "user_id", "title", "description", "address", "price", "area",
		"rooms_count", "bathrooms_count", "property_type", "deal_type", "material_type",
		"gas", "electricity", "internet", "sewerage", "plumbing",
		"renovation", "floor", "land_area", "floors").
		From("properties").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return property.Property{}, err
	}
	err = r.db.Get(&outProperty, sql, args...)
	if err != nil {
		return property.Property{}, err
	}
	return outProperty, nil
}

func (r *SearchRepository) SearchByFilters(filter PropertyFilter) ([]property.Property, error) {
	query := sq.Select("id", "user_id", "title", "description", "address", "price", "area",
		"rooms_count", "bathrooms_count", "property_type", "deal_type", "material_type",
		"gas", "electricity", "internet", "sewerage", "plumbing",
		"renovation", "floor", "land_area", "floors",
	).From("properties")

	if filter.MinPrice != nil {
		query = query.Where(sq.GtOrEq{"price": *filter.MinPrice})
	}
	if filter.MaxPrice != nil {
		query = query.Where(sq.LtOrEq{"price": *filter.MaxPrice})
	}
	if filter.MinArea != nil {
		query = query.Where(sq.GtOrEq{"area": *filter.MinArea})
	}
	if filter.MaxArea != nil {
		query = query.Where(sq.LtOrEq{"area": *filter.MaxArea})
	}
	if filter.RoomsCount != nil {
		query = query.Where(sq.Eq{"rooms_count": *filter.RoomsCount})
	}
	if filter.PropertyType != nil {
		query = query.Where(sq.Eq{"property_type": *filter.PropertyType})
	}
	if filter.DealType != nil {
		query = query.Where(sq.Eq{"deal_type": *filter.DealType})
	}
	if filter.BathroomsCount != nil {
		query = query.Where(sq.Eq{"bathrooms_count": *filter.BathroomsCount})
	}
	if filter.Material != nil {
		query = query.Where(sq.Eq{"material_type": *filter.Material})
	}
	if filter.Gas != nil {
		query = query.Where(sq.Eq{"gas": *filter.Gas})
	}
	if filter.Electricity != nil {
		query = query.Where(sq.Eq{"electricity": *filter.Electricity})
	}
	if filter.Internet != nil {
		query = query.Where(sq.Eq{"internet": *filter.Internet})
	}
	if filter.Sewerage != nil {
		query = query.Where(sq.Eq{"sewerage": *filter.Sewerage})
	}
	if filter.Plumbing != nil {
		query = query.Where(sq.Eq{"plumbing": *filter.Plumbing})
	}
	if filter.Renovation != nil {
		query = query.Where(sq.Eq{"renovation": *filter.Renovation})
	}
	if filter.Floor != nil {
		query = query.Where(sq.Eq{"floor": *filter.Floor})
	}
	if filter.MinLandArea != nil {
		query = query.Where(sq.GtOrEq{"land_area": *filter.MinLandArea})
	}
	if filter.MaxLandArea != nil {
		query = query.Where(sq.LtOrEq{"land_area": *filter.MaxLandArea})
	}
	if filter.Floors != nil {
		query = query.Where(sq.Eq{"floors": *filter.Floors})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var properties []property.Property
	err = r.db.Select(&properties, sql, args...)
	if err != nil {
		return nil, err
	}

	return properties, nil
}
