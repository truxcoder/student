package student

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"student/internal/repository"
	"student/internal/sqlc"
)

type Repository interface {
	Get(ctx context.Context, id string) (sqlc.Student, error)
	List(ctx context.Context) ([]sqlc.Student, error)
	Create(ctx context.Context, student sqlc.CreateStudentParams) error
	Update(ctx context.Context, student sqlc.UpdateStudentParams) error
	Delete(ctx context.Context, id string) error
}

type repo struct {
	*repository.Base
}

func NewRepository(base *repository.Base) Repository {
	return &repo{base}
}

func (r *repo) Get(ctx context.Context, id string) (sqlc.Student, error) {
	student, err := r.Q(ctx).GetStudentByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Student{}, fmt.Errorf("%w: id=%s", ErrNotFound, id)
		}
		return sqlc.Student{}, err
	}
	return student, nil
}

func (r *repo) List(ctx context.Context) ([]sqlc.Student, error) {
	students, err := r.Q(ctx).ListStudents(ctx)
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (r *repo) Create(ctx context.Context, student sqlc.CreateStudentParams) error {
	return r.Q(ctx).CreateStudent(ctx, student)
}

func (r *repo) Update(ctx context.Context, student sqlc.UpdateStudentParams) error {
	rows, err := r.Q(ctx).UpdateStudent(ctx, student)
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("%w: id=%s", ErrNotFound, student.ID)
	}
	return nil
}

func (r *repo) Delete(ctx context.Context, id string) error {
	rows, err := r.Q(ctx).DeleteStudent(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("%w: id=%s", ErrNotFound, id)
	}
	return nil
}

// User accepted prediction:
