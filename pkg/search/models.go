package search

import (
	"errors"

	"github.com/AkhmadeevRus/sublease-app/pkg/property"
)

type PropertyFilter struct {
	MinPrice       *int                   `form:"min_price"`
	MaxPrice       *int                   `form:"max_price"`
	Area           *int                   `form:"area"`
	RoomsCount     *int                   `form:"rooms_count"`
	BathroomsCount *int                   `form:"bathrooms_count"`
	PropertyType   *property.PropertyType `form:"property_type"`
	DealType       *property.DealType     `form:"deal_type"`
	Material       *property.MaterialType `form:"material_type"`
	Gas            *bool                  `form:"gas"`
	Electricity    *bool                  `form:"electricity"`
	Internet       *bool                  `form:"internet"`
	Sewerage       *bool                  `form:"sewerage"`
	Plumbing       *bool                  `form:"plumbing"`
	// for apartment
	Renovation *property.RenovationType `form:"renovation"`
	Floor      *int                     `form:"floor"`
	// for house
	LandArea *int `form:"land_area"`
	Floors   *int `form:"floors"`
}

func (p *PropertyFilter) Validate() error {
	if p.MinPrice != nil && p.MaxPrice != nil && *p.MinPrice > *p.MaxPrice {
		return errors.New("min_price can't be greater than max_price")
	}
	if p.MinPrice != nil && *p.MinPrice < 0 {
		return errors.New("min_price must be positive")
	}
	if p.MaxPrice != nil && *p.MaxPrice < 0 {
		return errors.New("max_price must be positive")
	}
	if p.Area != nil && *p.Area <= 0 {
		return errors.New("area must be positive")
	}
	if p.RoomsCount != nil && *p.RoomsCount <= 0 {
		return errors.New("rooms_count must be positive")
	}
	if p.BathroomsCount != nil && *p.BathroomsCount < 0 {
		return errors.New("bathrooms_count must be positive")
	}
	if p.Floor != nil && *p.Floor < 1 {
		return errors.New("floor must be at least 1")
	}
	if p.Floors != nil && *p.Floors < 1 {
		return errors.New("floors must be at least 1")
	}
	return nil
}
