package student

import (
	"time"

	"student/internal/platform/optional"
)

// StudentRequest 是创建学生的 HTTP 请求体，只承载传输层数据，
// 不依赖 sqlc 生成类型，也不负责生成主键。
type StudentRequest struct {
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	Birthday  time.Time `json:"birthday"`
	Score     *float64  `json:"score"`
	Classroom string    `json:"classroom"`
}

// StudentUpdateRequest 是 PATCH 局部更新的请求体：
// 未提供的字段保持库中原值，不参与更新。
type StudentUpdateRequest struct {
	ID        string                     `json:"id"`
	Name      *string                    `json:"name"`
	Gender    *string                    `json:"gender"`
	Birthday  *time.Time                 `json:"birthday"`
	Score     optional.Optional[float64] `json:"score"`
	Classroom *string                    `json:"classroom"`
}

// StudentResponse 是学生资源的 HTTP 响应体。
// Score 为指针：库中成绩可空，空值时序列化为 null。
type StudentResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	Birthday  time.Time `json:"birthday"`
	Score     *float64  `json:"score"`
	Classroom string    `json:"classroom"`
}
