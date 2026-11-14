// Package optional 提供能区分"字段缺失 / 显式 null / 有值"三态的 JSON 解码类型，
// 用于 PATCH 等局部更新场景，避免把"未提交"与"置空"混为一谈。
package optional

import (
	"bytes"
	"encoding/json"
)

// Optional 表示一个 JSON 字段的三种状态：
// 缺失（Set=false）、显式 null（Set=true 且 Value=nil）、有值（Set=true 且 Value!=nil）。
// encoding/json 只在字段出现在 JSON 中时才调用 UnmarshalJSON，
// 因此 Set 即"客户端是否提交了该字段"。
type Optional[T any] struct {
	Set   bool
	Value *T
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	o.Set = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		o.Value = nil
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	o.Value = &v
	return nil
}
