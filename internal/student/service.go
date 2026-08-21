package student

import (
	"context"
	"student/internal/service"
	"student/internal/sqlc"
)

type StudentService interface {
	Get(ctx context.Context, id string) (sqlc.Student, error)
	List(ctx context.Context) ([]sqlc.Student, error)
	Create(ctx context.Context, student *sqlc.CreateStudentParams) error
	Update(ctx context.Context, student *sqlc.UpdateStudentParams) error
	Delete(ctx context.Context, id string) error
}

type studentService struct {
	*service.Service
	repo StudentRepository
}

func NewStudentService(svc *service.Service, repo StudentRepository) StudentService {
	return &studentService{
		Service: svc,
		repo:    repo,
	}
}

func (s *studentService) Get(ctx context.Context, id string) (sqlc.Student, error) {
	return s.repo.Get(ctx, id)
}

func (s *studentService) List(ctx context.Context) ([]sqlc.Student, error) {
	return s.repo.List(ctx)
}

func (s *studentService) Create(ctx context.Context, student *sqlc.CreateStudentParams) error {
	return s.repo.Create(ctx, *student)
}

func (s *studentService) Update(ctx context.Context, student *sqlc.UpdateStudentParams) error {
	return s.repo.Update(ctx, *student)
}

func (s *studentService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
