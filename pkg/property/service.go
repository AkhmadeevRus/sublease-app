package property

import (
	"github.com/google/uuid"
)

type IPropertyService interface {
	CreateProperty(property Property, userId uuid.UUID) error
	DeleteProperty(userId, id uuid.UUID) error
	UpdateProperty(userId, id uuid.UUID, input PropertyUpdate) error
}

type PropertyService struct {
	repo IPropertyRepository
}

func NewPropertyService(repo IPropertyRepository) *PropertyService {
	return &PropertyService{repo: repo}
}

func (s *PropertyService) CreateProperty(property Property, userId uuid.UUID) error {
	if err := property.Validate(); err != nil {
		return err
	}
	return s.repo.CreateProperty(property, userId)
}

func (s *PropertyService) DeleteProperty(userId, id uuid.UUID) error {
	return s.repo.DeleteProperty(userId, id)
}

func (s *PropertyService) UpdateProperty(userId, id uuid.UUID, input PropertyUpdate) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.UpdateProperty(userId, id, input)
}
