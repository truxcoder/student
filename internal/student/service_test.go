package student

import (
	"context"
	"testing"
	"time"
	"uuid"

	"student/internal/platform/optional"
	"student/internal/service"
	"student/internal/sqlc"
)

// fakeRepository 记录 service 交给仓储的参数，用于断言映射结果，不触达数据库。
type fakeRepository struct {
	created sqlc.CreateStudentParams
	updated sqlc.UpdateStudentParams
}

func (f *fakeRepository) Get(context.Context, string) (sqlc.Student, error) {
	return sqlc.Student{}, ErrNotFound
}

func (f *fakeRepository) List(context.Context) ([]sqlc.Student, error) { return nil, nil }

func (f *fakeRepository) Create(_ context.Context, params sqlc.CreateStudentParams) error {
	f.created = params
	return nil
}

func (f *fakeRepository) Update(_ context.Context, params sqlc.UpdateStudentParams) error {
	f.updated = params
	return nil
}

func (f *fakeRepository) Delete(context.Context, string) error { return nil }

func newTestService(repo Repository) Service {
	return NewService(service.New(nil, nil), repo)
}

// TestServiceCreateGeneratesUUIDv7 验证主键由 service 生成，且版本号为 7。
func TestServiceCreateGeneratesUUIDv7(t *testing.T) {
	fake := &fakeRepository{}
	svc := newTestService(fake)
	score := 88.5
	birthday := time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)

	if err := svc.Create(context.Background(), StudentRequest{
		Name:      "张三",
		Gender:    "男",
		Birthday:  birthday,
		Score:     &score,
		Classroom: "一班",
	}); err != nil {
		t.Fatalf("Create: %v", err)
	}

	id, err := uuid.Parse(fake.created.ID)
	if err != nil {
		t.Fatalf("主键 %q 不是合法 UUID: %v", fake.created.ID, err)
	}
	if version := id[6] >> 4; version != 7 {
		t.Fatalf("主键版本 = %d, want 7 (UUIDv7)", version)
	}

	if fake.created.Name != "张三" ||
		fake.created.Gender != "男" ||
		!fake.created.Birthday.Equal(birthday) ||
		fake.created.Score == nil || *fake.created.Score != score ||
		fake.created.Classroom != "一班" {
		t.Fatalf("请求字段未正确映射到仓储参数: %+v", fake.created)
	}
}

// TestServiceUpdateScoreStates 覆盖 PATCH 的三态 score 到仓储参数的映射。
func TestServiceUpdateScoreStates(t *testing.T) {
	score := 99.5
	tests := []struct {
		name      string
		score     optional.Optional[float64]
		wantClear int64
		wantScore *float64
	}{
		{name: "字段缺失保留原值", score: optional.Optional[float64]{}},
		{name: "显式 null 清空", score: optional.Optional[float64]{Set: true}, wantClear: 1},
		{name: "有值覆盖", score: optional.Optional[float64]{Set: true, Value: &score}, wantScore: &score},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &fakeRepository{}
			svc := newTestService(fake)

			if err := svc.Update(context.Background(), StudentUpdateRequest{ID: "stu-1", Score: tt.score}); err != nil {
				t.Fatalf("Update: %v", err)
			}
			if fake.updated.ID != "stu-1" {
				t.Fatalf("ID = %q, want stu-1", fake.updated.ID)
			}
			if fake.updated.ClearScore != tt.wantClear {
				t.Fatalf("ClearScore = %d, want %d", fake.updated.ClearScore, tt.wantClear)
			}
			switch {
			case tt.wantScore == nil && fake.updated.Score != nil:
				t.Fatalf("Score = %v, want nil", *fake.updated.Score)
			case tt.wantScore != nil && fake.updated.Score == nil:
				t.Fatalf("Score = nil, want %v", *tt.wantScore)
			case tt.wantScore != nil && *fake.updated.Score != *tt.wantScore:
				t.Fatalf("Score = %v, want %v", *fake.updated.Score, *tt.wantScore)
			}
		})
	}
}

// TestServiceUpdateOmittedFieldsStayNil 确认未提交的字段不会被改写。
func TestServiceUpdateOmittedFieldsStayNil(t *testing.T) {
	fake := &fakeRepository{}
	svc := newTestService(fake)

	if err := svc.Update(context.Background(), StudentUpdateRequest{ID: "stu-1"}); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if fake.updated.Name != nil || fake.updated.Gender != nil ||
		fake.updated.Birthday != nil || fake.updated.Classroom != nil {
		t.Fatalf("未提交的字段不应写入仓储参数: %+v", fake.updated)
	}
}
