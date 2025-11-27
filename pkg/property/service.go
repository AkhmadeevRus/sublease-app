package property

import (
	"github.com/google/uuid"
)

type IPropertyService interface {
	CreateProperty(property Property, userId uuid.UUID) (uuid.UUID, error)
	GetMyProperties(id uuid.UUID) ([]Property, error)
	DeleteProperty(userId, id uuid.UUID) error
	UpdateProperty(userId, id uuid.UUID, input Property) error
}

type PropertyService struct {
	repo IPropertyRepository
}

func NewPropertyService(repo IPropertyRepository) *PropertyService {
	return &PropertyService{repo: repo}
}

func (s *PropertyService) CreateProperty(property Property, userId uuid.UUID) (uuid.UUID, error) {
	if err := property.Validate(); err != nil {
		return uuid.Nil, err
	}
	return s.repo.CreateProperty(property, userId)
}

func (s *PropertyService) GetMyProperties(id uuid.UUID) ([]Property, error) {
	return s.repo.GetMyProperties(id)
}

func (s *PropertyService) DeleteProperty(userId, id uuid.UUID) error {
	return s.repo.DeleteProperty(userId, id)
}

func (s *PropertyService) UpdateProperty(userId, id uuid.UUID, input Property) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.UpdateProperty(userId, id, input)
}
