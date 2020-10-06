package trace

import (
	"github.com/nozgurozturk/noo-analytics/entities"
	errors "github.com/nozgurozturk/noo-analytics/internal/utils"
	"time"
)

type Service interface {
	Insert(book *entities.TraceDTO) (*entities.Trace, *errors.ApplicationError)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) Insert(traceDto *entities.TraceDTO) (*entities.Trace, *errors.ApplicationError) {
	trace := entities.ToTrace(traceDto)

	now := time.Now()
	trace.TimeStamp = now.Unix()
	trace.Year = now.Year()
	trace.Month = int(now.Month())
	trace.Day = now.Day()
	trace.Hour = now.Hour()

	createTrace, err := s.repository.Insert(trace)

	if err != nil {
		return nil, err
	}

	return createTrace, nil
}
