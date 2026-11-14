package student

import (
	"context"
	"uuid"

	"student/internal/platform/optional"
	"student/internal/service"
	"student/internal/sqlc"
)

type Service interface {
	Get(ctx context.Context, id string) (sqlc.Student, error)
	List(ctx context.Context) ([]sqlc.Student, error)
	Create(ctx context.Context, req StudentRequest) error
	Update(ctx context.Context, req StudentUpdateRequest) error
	Delete(ctx context.Context, id string) error
}

type svc struct {
	*service.Base
	repo Repository
}

func NewService(base *service.Base, repo Repository) Service {
	return &svc{
		Base: base,
		repo: repo,
	}
}

func (s *svc) Get(ctx context.Context, id string) (sqlc.Student, error) {
	return s.repo.Get(ctx, id)
}

func (s *svc) List(ctx context.Context) ([]sqlc.Student, error) {
	return s.repo.List(ctx)
}

// Create 生成主键并落库：主键生成属于应用层职责，不下放到传输层 DTO。
func (s *svc) Create(ctx context.Context, req StudentRequest) error {
	return s.repo.Create(ctx, sqlc.CreateStudentParams{
		ID:        uuid.NewV7().String(),
		Name:      req.Name,
		Gender:    req.Gender,
		Birthday:  req.Birthday,
		Score:     req.Score,
		Classroom: req.Classroom,
	})
}

// Update 执行 PATCH 局部更新：未提供的字段传 nil，由 SQL 层保留原值；
// score 显式传 null 时清空为 NULL。
func (s *svc) Update(ctx context.Context, req StudentUpdateRequest) error {
	return s.repo.Update(ctx, sqlc.UpdateStudentParams{
		ID:         req.ID,
		Name:       req.Name,
		Gender:     req.Gender,
		Birthday:   req.Birthday,
		Score:      req.Score.Value,
		ClearScore: clearScore(req.Score),
		Classroom:  req.Classroom,
	})
}

// clearScore 返回 1 表示客户端显式提交了 score: null，需要清空成绩；
// 字段缺失或提供了数值时返回 0，由 SQL 的 COALESCE 保留原值或写入新值。
func clearScore(score optional.Optional[float64]) int64 {
	if score.Set && score.Value == nil {
		return 1
	}
	return 0
}

func (s *svc) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
