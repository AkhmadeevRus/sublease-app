package search

import (
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
