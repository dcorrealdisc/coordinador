package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/dcorreal/coordinador/internal/models"
	"github.com/dcorreal/coordinador/internal/repositories"
)

// StudentRepository is a mock implementation of repositories.StudentRepository.
type StudentRepository struct {
	mock.Mock
}

func (m *StudentRepository) Create(ctx context.Context, student *models.Student) error {
	args := m.Called(ctx, student)
	return args.Error(0)
}

func (m *StudentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Student, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Student), args.Error(1)
}

func (m *StudentRepository) List(ctx context.Context, filters repositories.StudentFilters) ([]*models.Student, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Student), args.Error(1)
}

func (m *StudentRepository) Update(ctx context.Context, student *models.Student) error {
	args := m.Called(ctx, student)
	return args.Error(0)
}

func (m *StudentRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy *uuid.UUID) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}

func (m *StudentRepository) Count(ctx context.Context, filters repositories.StudentFilters) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}
