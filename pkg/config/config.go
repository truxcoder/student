package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Supported 配置值允许的类型约束
// 编译器会在实例化时拒绝不支持的类型，而非运行时静默失败
type Supported interface {
	string | int | int8 | int16 | int32 | int64 |
		float32 | float64 | bool |
		[]string | []int | []float64 | []bool |
		map[string]interface{}
}

// Config 配置管理器（线程安全）
type Config struct {
	mu        sync.RWMutex
	data      map[string]interface{}
	delimiter string
}

// Option 配置选项函数
type Option func(*Config)

// WithDelimiter 设置键分隔符，默认为 "."
func WithDelimiter(delimiter string) Option {
	return func(c *Config) {
		c.delimiter = delimiter
	}
}

// New 创建新的配置实例
func New(opts ...Option) *Config {
	c := &Config{
		data:      make(map[string]interface{}),
		delimiter: ".",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// NewFromFile 从 JSON 文件加载配置
func NewFromFile(filename string, opts ...Option) (*Config, error) {
	c := New(opts...)

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	if err := json.Unmarshal(data, &c.data); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	return c, nil
}

// ============================================================
// 核心 API —— 泛型获取
// ============================================================

// Get 泛型获取配置值，支持自动类型转换
//
//	config.Get[string](cfg, "app.name")
//	config.Get[int](cfg, "app.port")
//	config.Get[bool](cfg, "app.debug")
func Get[T Supported](c *Config, key string) T {
	c.mu.RLock()
	val := c.lookup(key)
	c.mu.RUnlock()
	return convertTo[T](val)
}

// GetOr 泛型获取配置值，键不存在或为零值时返回默认值
//
//	config.GetOr(cfg, "app.host", "localhost")
//	config.GetOr(cfg, "app.port", 8080)
func GetOr[T Supported](c *Config, key string, defaultVal T) T {
	c.mu.RLock()
	val := c.lookup(key)
	c.mu.RUnlock()

	if val == nil {
		return defaultVal
	}
	result := convertTo[T](val)
	if isZeroValue(result) {
		return defaultVal
	}
	return result
}

// ============================================================
// 便捷方法
// ============================================================

// GetString 获取字符串配置
func (c *Config) GetString(key string) string {
	return Get[string](c, key)
}

// GetInt 获取整数配置
func (c *Config) GetInt(key string) int {
	return Get[int](c, key)
}

// GetInt64 获取 int64 配置
func (c *Config) GetInt64(key string) int64 {
	return Get[int64](c, key)
}

// GetFloat64 获取浮点数配置
func (c *Config) GetFloat64(key string) float64 {
	return Get[float64](c, key)
}

// GetBool 获取布尔配置
func (c *Config) GetBool(key string) bool {
	return Get[bool](c, key)
}

// GetDuration 获取时间间隔配置
func (c *Config) GetDuration(key string) time.Duration {
	s := c.GetString(key)
	if s == "" {
		return 0
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return d
}

// GetDurationOr 获取时间间隔配置，带默认值
func (c *Config) GetDurationOr(key string, defaultVal time.Duration) time.Duration {
	d := c.GetDuration(key)
	if d == 0 {
		return defaultVal
	}
	return d
}

// GetStringSlice 获取字符串切片配置
func (c *Config) GetStringSlice(key string) []string {
	return Get[[]string](c, key)
}

// GetStringMap 获取字符串映射配置
func (c *Config) GetStringMap(key string) map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val := c.lookup(key)
	if val == nil {
		return nil
	}
	if m, ok := val.(map[string]interface{}); ok {
		return m
	}
	return nil
}

// ============================================================
// 读写操作
// ============================================================

// Set 设置配置值，支持点号分隔的键路径
func (c *Config) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.set(key, value)
}

// Has 检查键是否存在
func (c *Config) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := strings.Split(key, c.delimiter)
	current := c.data

	for i, k := range keys {
		if i == len(keys)-1 {
			_, ok := current[k]
			return ok
		}
		if nested, ok := current[k].(map[string]interface{}); ok {
			current = nested
		} else {
			return false
		}
	}
	return false
}

// Delete 删除配置键
func (c *Config) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	keys := strings.Split(key, c.delimiter)
	current := c.data

	for i, k := range keys {
		if i == len(keys)-1 {
			delete(current, k)
			return
		}
		if nested, ok := current[k].(map[string]interface{}); ok {
			current = nested
		} else {
			return
		}
	}
}

// All 返回所有配置的深拷贝
func (c *Config) All() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return deepCopy(c.data).(map[string]interface{})
}

// Keys 返回所有顶级键
func (c *Config) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.data))
	for k := range c.data {
		keys = append(keys, k)
	}
	return keys
}

// Merge 合并另一个配置（other 的值覆盖 c 的同名键）
func (c *Config) Merge(other *Config) {
	if other == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	other.mu.RLock()
	defer other.mu.RUnlock()

	mergeMaps(c.data, other.data)
}

// ============================================================
// 环境变量加载
// ============================================================

// LoadFromEnv 从环境变量加载配置
// 环境变量名去掉 prefix 后转小写，下划线替换为分隔符
func (c *Config) LoadFromEnv(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, prefix) {
			continue
		}

		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			continue
		}

		key := strings.TrimPrefix(pair[0], prefix)
		key = strings.ToLower(key)
		key = strings.ReplaceAll(key, "_", c.delimiter)

		c.set(key, pair[1])
	}
}

// ============================================================
// 反序列化
// ============================================================

// Unmarshal 将全部配置反序列化到结构体
func (c *Config) Unmarshal(v interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return unmarshalMap(c.data, v)
}

// UnmarshalKey 将指定键的配置反序列化到结构体
func (c *Config) UnmarshalKey(key string, v interface{}) error {
	c.mu.RLock()
	val := c.lookup(key)
	c.mu.RUnlock()

	if val == nil {
		return fmt.Errorf("键 '%s' 不存在", key)
	}
	if m, ok := val.(map[string]interface{}); ok {
		return unmarshalMap(m, v)
	}
	return fmt.Errorf("键 '%s' 的值不是对象", key)
}

// ============================================================
// 持久化
// ============================================================

// SaveToFile 保存配置到 JSON 文件
func (c *Config) SaveToFile(filename string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, err := json.MarshalIndent(c.data, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	return os.WriteFile(filename, data, 0644)
}

// ============================================================
// 内部方法
// ============================================================

// lookup 内部无锁查找，调用者须已持有读锁
func (c *Config) lookup(key string) interface{} {
	keys := strings.Split(key, c.delimiter)
	current := c.data

	for i, k := range keys {
		if i == len(keys)-1 {
			return current[k]
		}
		if nested, ok := current[k].(map[string]interface{}); ok {
			current = nested
		} else {
			return nil
		}
	}
	return nil
}

// set 内部无锁设置，调用者须已持有写锁
func (c *Config) set(key string, value interface{}) {
	keys := strings.Split(key, c.delimiter)
	current := c.data

	for i, k := range keys {
		if i == len(keys)-1 {
			current[k] = value
			return
		}
		if _, ok := current[k]; !ok {
			current[k] = make(map[string]interface{})
		}
		if nested, ok := current[k].(map[string]interface{}); ok {
			current = nested
		} else {
			current[k] = make(map[string]interface{})
			current = current[k].(map[string]interface{})
		}
	}
}

// ============================================================
// 泛型类型转换
// ============================================================

func convertTo[T Supported](val interface{}) T {
	var zero T
	if val == nil {
		return zero
	}

	if v, ok := val.(T); ok {
		return v
	}

	targetType := reflect.TypeOf(zero)
	if targetType == nil {
		return zero
	}

	valReflect := reflect.ValueOf(val)

	switch targetType.Kind() {
	case reflect.String:
		return any(fmt.Sprintf("%v", val)).(T)
	case reflect.Int:
		return any(int(toFloat64(val))).(T)
	case reflect.Int8:
		return any(int8(toFloat64(val))).(T)
	case reflect.Int16:
		return any(int16(toFloat64(val))).(T)
	case reflect.Int32:
		return any(int32(toFloat64(val))).(T)
	case reflect.Int64:
		return any(toInt64(val)).(T)
	case reflect.Float32:
		return any(float32(toFloat64(val))).(T)
	case reflect.Float64:
		return any(toFloat64(val)).(T)
	case reflect.Bool:
		return any(toBool(val)).(T)
	case reflect.Slice:
		return convertSlice[T](val, targetType)
	case reflect.Map:
		if m, ok := val.(map[string]interface{}); ok {
			return any(m).(T)
		}
	default:
		// 不支持的目标类型，尝试 reflect 兜底
		if valReflect.CanConvert(targetType) {
			return valReflect.Convert(targetType).Interface().(T)
		}
	}

	return zero
}

func convertSlice[T Supported](val interface{}, targetType reflect.Type) T {
	var zero T

	switch v := val.(type) {
	case []string:
		slice := reflect.MakeSlice(targetType, len(v), len(v))
		for i, s := range v {
			setSliceElement(slice.Index(i), s)
		}
		return slice.Interface().(T)

	case []interface{}:
		slice := reflect.MakeSlice(targetType, len(v), len(v))
		for i, item := range v {
			elemVal := reflect.ValueOf(item)
			elemType := targetType.Elem()
			if elemVal.Type().AssignableTo(elemType) {
				slice.Index(i).Set(elemVal)
			} else if elemVal.CanConvert(elemType) {
				slice.Index(i).Set(elemVal.Convert(elemType))
			} else {
				setSliceElement(slice.Index(i), fmt.Sprintf("%v", item))
			}
		}
		return slice.Interface().(T)

	case string:
		if strings.Contains(v, ",") {
			parts := strings.Split(v, ",")
			if targetType.Elem().Kind() == reflect.String {
				result := make([]string, len(parts))
				for i, p := range parts {
					result[i] = strings.TrimSpace(p)
				}
				return any(result).(T)
			}
			slice := reflect.MakeSlice(targetType, len(parts), len(parts))
			for i, p := range parts {
				setSliceElement(slice.Index(i), strings.TrimSpace(p))
			}
			return slice.Interface().(T)
		}
		slice := reflect.MakeSlice(targetType, 1, 1)
		setSliceElement(slice.Index(0), v)
		return slice.Interface().(T)

	default:
		// 无法转换，返回零值
	}

	return zero
}

func setSliceElement(elem reflect.Value, s string) {
	switch elem.Kind() {
	case reflect.String:
		elem.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			elem.SetInt(i)
		}
	case reflect.Float32, reflect.Float64:
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			elem.SetFloat(f)
		}
	case reflect.Bool:
		if b, err := strconv.ParseBool(s); err == nil {
			elem.SetBool(b)
		}
	default:
		// 不支持的切片元素类型，忽略
	}
}

// isZeroValue 安全判断零值，避免对不可比较类型（slice/map 等）panic
func isZeroValue[T Supported](val T) bool {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	default:
		return v.IsZero()
	}
}

// ============================================================
// 基础类型转换
// ============================================================

func toFloat64(val interface{}) float64 {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case int32:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	default:
		// 不支持的类型
	}
	return 0
}

func toInt64(val interface{}) int64 {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	default:
		// 不支持的类型
	}
	return 0
}

func toBool(val interface{}) bool {
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case string:
		switch strings.ToLower(v) {
		case "true", "yes", "on", "1":
			return true
		default:
			return false
		}
	default:
		return false
	}
}

// ============================================================
// 工具函数
// ============================================================

func mergeMaps(dst, src map[string]interface{}) {
	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			srcMap, srcOk := srcVal.(map[string]interface{})
			dstMap, dstOk := dstVal.(map[string]interface{})
			if srcOk && dstOk {
				mergeMaps(dstMap, srcMap)
				continue
			}
		}
		dst[key] = srcVal
	}
}

func deepCopy(src interface{}) interface{} {
	if src == nil {
		return nil
	}

	srcVal := reflect.ValueOf(src)
	switch srcVal.Kind() {
	case reflect.Map:
		dst := make(map[string]interface{})
		for _, key := range srcVal.MapKeys() {
			dst[key.String()] = deepCopy(srcVal.MapIndex(key).Interface())
		}
		return dst
	case reflect.Slice:
		dst := make([]interface{}, srcVal.Len())
		for i := 0; i < srcVal.Len(); i++ {
			dst[i] = deepCopy(srcVal.Index(i).Interface())
		}
		return dst
	default:
		return src
	}
}

// ============================================================
// 反序列化工具
// ============================================================

func unmarshalMap(data map[string]interface{}, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("v 必须是非空指针")
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("v 必须指向结构体")
	}

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldValue := rv.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		tag := field.Tag.Get("config")
		if tag == "" {
			tag = field.Tag.Get("json")
		}
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}
		if idx := strings.Index(tag, ","); idx != -1 {
			tag = tag[:idx]
		}
		if tag == "-" {
			continue
		}

		val, ok := data[tag]
		if !ok || val == nil {
			continue
		}

		// 解引用指针字段
		actualField := fieldValue
		if field.Type.Kind() == reflect.Ptr {
			if actualField.IsNil() {
				actualField.Set(reflect.New(field.Type.Elem()))
			}
			actualField = actualField.Elem()
		}

		// 嵌套结构体递归
		if actualField.Kind() == reflect.Struct {
			if nestedMap, nestedOk := val.(map[string]interface{}); nestedOk {
				if err := unmarshalMap(nestedMap, actualField.Addr().Interface()); err != nil {
					return err
				}
				continue
			}
		}

		setFieldValue(actualField, val)
	}

	return nil
}

func setFieldValue(field reflect.Value, val interface{}) {
	if !field.IsValid() || !field.CanSet() {
		return
	}

	valReflect := reflect.ValueOf(val)

	if valReflect.Type() == field.Type() {
		field.Set(valReflect)
		return
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(fmt.Sprintf("%v", val))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.SetInt(int64(toFloat64(val)))
	case reflect.Float32, reflect.Float64:
		field.SetFloat(toFloat64(val))
	case reflect.Bool:
		field.SetBool(toBool(val))
	case reflect.Struct:
		if m, ok := val.(map[string]interface{}); ok {
			_ = unmarshalMap(m, field.Addr().Interface())
		}
	case reflect.Slice:
		if arr, ok := val.([]interface{}); ok {
			slice := reflect.MakeSlice(field.Type(), len(arr), len(arr))
			for i, item := range arr {
				setFieldValue(slice.Index(i), item)
			}
			field.Set(slice)
		}
	case reflect.Map:
		if m, ok := val.(map[string]interface{}); ok {
			mapVal := reflect.MakeMap(field.Type())
			for k, v := range m {
				mapVal.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
			}
			field.Set(mapVal)
		}
	default:
		// 不支持的字段类型，忽略
	}
}
