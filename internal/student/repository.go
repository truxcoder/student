package student

import (
	"context"
	"student/internal/repository"
	"student/internal/sqlc"
)

type StudentRepository interface {
	Get(ctx context.Context, id string) (sqlc.Student, error)
	List(ctx context.Context) ([]sqlc.Student, error)
	Create(ctx context.Context, student sqlc.CreateStudentParams) error
	Update(ctx context.Context, student sqlc.UpdateStudentParams) error
	Delete(ctx context.Context, id string) error
}

type studentRepository struct {
	*repository.Repository
}

func NewStudentRepository(repo *repository.Repository) StudentRepository {
	return &studentRepository{repo}
}

func (r *studentRepository) Get(ctx context.Context, id string) (sqlc.Student, error) {
	student, err := r.DB.GetStudentByID(ctx, id)
	if err != nil {
		return sqlc.Student{}, err
	}
	return student, nil
}

func (r *studentRepository) List(ctx context.Context) ([]sqlc.Student, error) {
	students, err := r.DB.ListStudents(ctx)
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (r *studentRepository) Create(ctx context.Context, student sqlc.CreateStudentParams) error {
	return r.DB.CreateStudent(ctx, student)
}

func (r *studentRepository) Update(ctx context.Context, student sqlc.UpdateStudentParams) error {
	return r.DB.UpdateStudent(ctx, student)
}

func (r *studentRepository) Delete(ctx context.Context, id string) error {
	return r.DB.DeleteStudent(ctx, id)
}

// User accepted prediction:
