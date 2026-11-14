package optional

import (
	"encoding/json"
	"testing"
)

func TestOptionalJSONStates(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		wantSet   bool
		wantValue *float64
	}{
		{name: "字段缺失", body: `{}`},
		{name: "显式 null", body: `{"score":null}`, wantSet: true},
		{name: "有值", body: `{"score":88.5}`, wantSet: true, wantValue: floatPtr(88.5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got struct {
				Score Optional[float64] `json:"score"`
			}
			if err := json.Unmarshal([]byte(tt.body), &got); err != nil {
				t.Fatalf("Unmarshal(%s): %v", tt.body, err)
			}
			if got.Score.Set != tt.wantSet {
				t.Fatalf("Set = %v, want %v", got.Score.Set, tt.wantSet)
			}
			assertFloatPtr(t, got.Score.Value, tt.wantValue)
		})
	}
}

func TestOptionalInvalidValue(t *testing.T) {
	var got struct {
		Score Optional[float64] `json:"score"`
	}
	if err := json.Unmarshal([]byte(`{"score":"abc"}`), &got); err == nil {
		t.Fatal("值类型不匹配时应当返回错误")
	}
}

// TestOptionalOtherTypes 确认泛型对非数值类型同样可用。
func TestOptionalOtherTypes(t *testing.T) {
	var got struct {
		Name Optional[string] `json:"name"`
	}
	if err := json.Unmarshal([]byte(`{"name":"张三"}`), &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !got.Name.Set || got.Name.Value == nil || *got.Name.Value != "张三" {
		t.Fatalf("Set=%v Value=%v, want Set=true Value=张三", got.Name.Set, got.Name.Value)
	}
}

func floatPtr(v float64) *float64 { return &v }

func assertFloatPtr(t *testing.T, got, want *float64) {
	t.Helper()
	switch {
	case want == nil && got != nil:
		t.Fatalf("Value = %v, want nil", *got)
	case want != nil && got == nil:
		t.Fatalf("Value = nil, want %v", *want)
	case want != nil && *got != *want:
		t.Fatalf("Value = %v, want %v", *got, *want)
	}
}
