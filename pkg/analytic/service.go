package analytics

import (
	"github.com/nozgurozturk/noo-analytics/entities"
	errors "github.com/nozgurozturk/noo-analytics/internal/utils"
)

type Service interface {
	FindActionsByDate(entity *entities.AnalyticsActionRequest) ([]entities.AnalyticsActionResponse, *errors.ApplicationError)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) FindActionsByDate(entity *entities.AnalyticsActionRequest) ([]entities.AnalyticsActionResponse, *errors.ApplicationError) {
	actions, err := s.repository.FindActionsByDate(entity)
	if err != nil {
		return nil, err
	}

	return actions, nil
}
