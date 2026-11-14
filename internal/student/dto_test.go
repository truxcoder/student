package student

import (
	"encoding/json"
	"testing"
)

// TestStudentUpdateRequestScoreStates 覆盖 PATCH 请求体中 score 的三种 JSON 状态。
func TestStudentUpdateRequestScoreStates(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		wantSet   bool
		wantValue *float64
	}{
		{name: "字段缺失", body: `{"id":"stu-1"}`},
		{name: "显式 null", body: `{"id":"stu-1","score":null}`, wantSet: true},
		{name: "有值", body: `{"id":"stu-1","score":88.5}`, wantSet: true, wantValue: float64Ptr(88.5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req StudentUpdateRequest
			if err := json.Unmarshal([]byte(tt.body), &req); err != nil {
				t.Fatalf("Unmarshal(%s): %v", tt.body, err)
			}
			if req.ID != "stu-1" {
				t.Fatalf("ID = %q, want stu-1", req.ID)
			}
			if req.Score.Set != tt.wantSet {
				t.Fatalf("Score.Set = %v, want %v", req.Score.Set, tt.wantSet)
			}
			switch {
			case tt.wantValue == nil && req.Score.Value != nil:
				t.Fatalf("Score.Value = %v, want nil", *req.Score.Value)
			case tt.wantValue != nil && req.Score.Value == nil:
				t.Fatalf("Score.Value = nil, want %v", *tt.wantValue)
			case tt.wantValue != nil && *req.Score.Value != *tt.wantValue:
				t.Fatalf("Score.Value = %v, want %v", *req.Score.Value, *tt.wantValue)
			}
		})
	}
}

func TestStudentUpdateRequestInvalidScore(t *testing.T) {
	var req StudentUpdateRequest
	if err := json.Unmarshal([]byte(`{"id":"stu-1","score":"abc"}`), &req); err == nil {
		t.Fatal("score 类型不匹配时应当返回解码错误")
	}
}

func float64Ptr(v float64) *float64 { return &v }
