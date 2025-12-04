package property

import (
	"errors"

	"github.com/google/uuid"
)

type PropertyType string
type DealType string
type MaterialType string
type RenovationType string

const (
	PropertyTypeHouse     PropertyType = "house"
	PropertyTypeApartment PropertyType = "apartment"

	DealTypeDailyRent DealType = "daily_rent"
	DealTypeLongRent  DealType = "long_rent"
	DealTypePurchase  DealType = "purchase"

	MaterialTypeBrick    MaterialType = "brick"
	MaterialTypeConcrete MaterialType = "concrete"
	MaterialTypeWood     MaterialType = "wood"
	MaterialTypeSIP      MaterialType = "sip"
	MaterialTypePannel   MaterialType = "pannel"

	RenovationTypeClean          RenovationType = "clean"
	RenovationTypeWithoutRepairs RenovationType = "without_repairs"
	RenovationTypeCosmetic       RenovationType = "cosmetic"
)

type Property struct {
	Id             uuid.UUID    `json:"id" db:"id"`
	UserId         uuid.UUID    `json:"user_id" db:"user_id"`
	Title          string       `json:"title" binding:"required" db:"title"`
	Description    string       `json:"description" binding:"required" db:"description"`
	Address        string       `json:"address" binding:"required" db:"address"`
	Price          int          `json:"price" binding:"required" db:"price"`
	Area           int          `json:"area" binding:"required" db:"area"`
	RoomsCount     int          `json:"rooms_count" binding:"required" db:"rooms_count"`
	BathroomsCount int          `json:"bathrooms_count" binding:"required" db:"bathrooms_count"`
	PropertyType   PropertyType `json:"property_type" binding:"required" db:"property_type"`
	DealType       DealType     `json:"deal_type" binding:"required" db:"deal_type"`
	Material       MaterialType `json:"material_type" binding:"required" db:"material_type"`
	Gas            bool         `json:"gas" db:"gas"`
	Electricity    bool         `json:"electricity" db:"electricity"`
	Internet       bool         `json:"internet" db:"internet"`
	Sewerage       bool         `json:"sewerage" db:"sewerage"`
	Plumbing       bool         `json:"plumbing" db:"plumbing"`
	// for apartment
	Renovation *RenovationType `json:"renovation,omitempty" db:"renovation"`
	Floor      *int            `json:"floor,omitempty" db:"floor"`
	// for house
	LandArea *int `json:"land_area,omitempty" db:"land_area"`
	Floors   *int `json:"floors,omitempty" db:"floors"`
}

func (p *Property) Validate() error {
	if p.Area <= 0 {
		return errors.New("area must be positive")
	}
	if p.RoomsCount <= 0 {
		return errors.New("rooms_count must be positive")
	}
	if p.BathroomsCount < 0 {
		return errors.New("bathrooms_count must be")
	}
	if p.Price <= 0 {
		return errors.New("price must be positive")
	}

	switch p.PropertyType {
	case PropertyTypeHouse:
		if p.LandArea == nil || *p.LandArea <= 0 {
			return errors.New("land_area is required for house")
		}
		if p.Floors == nil || *p.Floors <= 0 {
			return errors.New("floors is required for house")
		}
	case PropertyTypeApartment:
		if p.Floor == nil || *p.Floor <= 0 {
			return errors.New("floor is required for apartment")
		}
	default:
		return errors.New("unknown property type")
	}

	return nil
}

type PropertyUpdate struct {
	Title          *string         `json:"title,omitempty" db:"title"`
	Description    *string         `json:"description,omitempty" db:"description"`
	Address        *string         `json:"address,omitempty" db:"address"`
	Price          *int            `json:"price,omitempty" db:"price"`
	Area           *int            `json:"area,omitempty" db:"area"`
	RoomsCount     *int            `json:"rooms_count,omitempty" db:"rooms_count"`
	BathroomsCount *int            `json:"bathrooms_count,omitempty" db:"bathrooms_count"`
	PropertyType   *PropertyType   `json:"property_type,omitempty" db:"property_type"`
	DealType       *DealType       `json:"deal_type,omitempty" db:"deal_type"`
	Material       *MaterialType   `json:"material_type,omitempty" db:"material_type"`
	Gas            *bool           `json:"gas,omitempty" db:"gas"`
	Electricity    *bool           `json:"electricity,omitempty" db:"electricity"`
	Internet       *bool           `json:"internet,omitempty" db:"internet"`
	Sewerage       *bool           `json:"sewerage,omitempty" db:"sewerage"`
	Plumbing       *bool           `json:"plumbing,omitempty" db:"plumbing"`
	Renovation     *RenovationType `json:"renovation,omitempty" db:"renovation"`
	Floor          *int            `json:"floor,omitempty" db:"floor"`
	LandArea       *int            `json:"land_area,omitempty" db:"land_area"`
	Floors         *int            `json:"floors,omitempty" db:"floors"`
}

func (p *PropertyUpdate) Validate() error {
	if p.Area != nil && *p.Area <= 0 {
		return errors.New("area must be positive")
	}
	if p.RoomsCount != nil && *p.RoomsCount <= 0 {
		return errors.New("rooms_count must be positive")
	}
	if p.BathroomsCount != nil && *p.BathroomsCount < 0 {
		return errors.New("bathrooms_count must be non-negative")
	}
	if p.Price != nil && *p.Price <= 0 {
		return errors.New("price must be positive")
	}

	if p.PropertyType != nil {
		switch *p.PropertyType {
		case PropertyTypeHouse:
			if p.LandArea != nil && *p.LandArea <= 0 {
				return errors.New("land_area is required for house")
			}
			if p.Floors != nil && *p.Floors <= 0 {
				return errors.New("floors is required for house")
			}
		case PropertyTypeApartment:
			if p.Floor != nil && *p.Floor <= 0 {
				return errors.New("floor is required for apartment")
			}
		default:
			return errors.New("unknown property type")
		}
	}

	return nil
}
