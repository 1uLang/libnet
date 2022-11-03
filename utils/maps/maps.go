package maps

import (
	"encoding/json"
	"github.com/1uLang/libnet/utils/types"
	"reflect"
)

type Map map[string]interface{}

// NewMap 新建Map
func NewMap(maps ...interface{}) Map {
	m := Map{}
	for _, mp := range maps {
		v := reflect.ValueOf(mp)
		if v.Kind() != reflect.Map {
			continue
		}

		for _, k := range v.MapKeys() {
			m[types.String(k.Interface())] = v.MapIndex(k).Interface()
		}
	}
	return m
}

// DecodeJSON 从字节数据中解码map
func DecodeJSON(jsonData []byte) (Map, error) {
	m := Map{}
	err := json.Unmarshal(jsonData, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

// Keys 取得所有键
func (this Map) Keys() []interface{} {
	m := []interface{}{}
	for key := range this {
		m = append(m, key)
	}
	return m
}

// Values 取得所有值
func (this Map) Values() []interface{} {
	m := []interface{}{}
	for _, value := range this {
		m = append(m, value)
	}
	return m
}

// Has 判断是否有某个键值
func (this Map) Has(key string) bool {
	_, found := this[key]
	return found
}

// Get 取得键值
func (this Map) Get(key string) interface{} {
	return this[key]
}

// GetBool 取得bool类型的键值
func (this Map) GetBool(key string) bool {
	return types.Bool(this[key])
}

// GetUint 取得uint类型的键值
func (this Map) GetUint(key string) uint {
	return types.Uint(this[key])
}

// GetUint8 取得uint8类型的键值
func (this Map) GetUint8(key string) uint8 {
	return types.Uint8(this[key])
}

// GetUint16 取得uint16类型的键值
func (this Map) GetUint16(key string) uint16 {
	return types.Uint16(this[key])
}

// GetUint32 取得uint32类型的键值
func (this Map) GetUint32(key string) uint32 {
	return types.Uint32(this[key])
}

// GetUint64 取得uint64类型的键值
func (this Map) GetUint64(key string) uint64 {
	return types.Uint64(this[key])
}

// GetInt 取得int类型的键值
func (this Map) GetInt(key string) int {
	return types.Int(this[key])
}

// GetInt8 取得int8类型的键值
func (this Map) GetInt8(key string) int8 {
	return types.Int8(this[key])
}

// GetInt16 取得int16类型的键值
func (this Map) GetInt16(key string) int16 {
	return types.Int16(this[key])
}

// GetInt32 取得int32类型的键值
func (this Map) GetInt32(key string) int32 {
	return types.Int32(this[key])
}

// GetInt64 取得int64类型的键值
func (this Map) GetInt64(key string) int64 {
	return types.Int64(this[key])
}

// GetFloat32 取得float32类型的键值
func (this Map) GetFloat32(key string) float32 {
	return types.Float32(this[key])
}

// GetFloat64 取得float64类型的键值
func (this Map) GetFloat64(key string) float64 {
	return types.Float64(this[key])
}

// Increase 给某个键值增加数值（可以为负），并返回操作后的值
func (this Map) Increase(key string, delta interface{}) interface{} {
	value, found := this[key]
	if !found || value == nil {
		this[key] = delta
	} else {
		switch value := value.(type) {
		case uint:
			this[key] = value + types.Uint(delta)
		case uint8:
			this[key] = value + types.Uint8(delta)
		case uint16:
			this[key] = value + types.Uint16(delta)
		case uint32:
			this[key] = value + types.Uint32(delta)
		case uint64:
			this[key] = value + types.Uint64(delta)
		case int:
			this[key] = value + types.Int(delta)
		case int8:
			this[key] = value + types.Int8(delta)
		case int16:
			this[key] = value + types.Int16(delta)
		case int32:
			this[key] = value + types.Int32(delta)
		case int64:
			this[key] = value + types.Int64(delta)
		case float32:
			this[key] = value + types.Float32(delta)
		case float64:
			this[key] = value + types.Float64(delta)
		}
	}

	return this[key]
}

// GetString 取得字符串类型的键值
func (this Map) GetString(key string) string {
	return types.String(this[key])
}

// GetMap 取得Map类型的键值
func (this Map) GetMap(key string) Map {
	value, found := this[key]
	if !found || value == nil || reflect.TypeOf(value).Kind() != reflect.Map {
		return nil
	}
	return NewMap(value)
}

// GetSlice 取得Slice类型的键值
func (this Map) GetSlice(key string) []interface{} {
	value, found := this[key]
	if !found || value == nil || reflect.TypeOf(value).Kind() != reflect.Slice {
		return nil
	}
	result := []interface{}{}
	v := reflect.ValueOf(value)
	count := v.Len()
	for i := 0; i < count; i++ {
		item := v.Index(i).Interface()
		result = append(result, item)
	}
	return result
}

// GetBytes 取得Byte Slice键值
func (this Map) GetBytes(key string) []byte {
	value, ok := this[key]
	if !ok || value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if ok {
		return bytes
	}

	s, ok := value.(string)
	if ok {
		return []byte(s)
	}

	return nil
}

// Delete 删除键
func (this Map) Delete(key ...string) {
	for _, oneKey := range key {
		delete(this, oneKey)
	}
}

// Put 添加键值
func (this Map) Put(key string, value interface{}) {
	this[key] = value
}

// Len 取得键值数量
func (this Map) Len() int {
	return len(this)
}

// GoMap 转换为map[string]interface{}
func (this Map) GoMap() map[string]interface{} {
	return this
}

// AsJSON 转换为JSON
func (this Map) AsJSON() []byte {
	data, err := json.Marshal(this)
	if err != nil {
		return []byte{}
	}
	return data
}

// AsPrettyJSON 转换为格式化后的JSON
func (this Map) AsPrettyJSON() []byte {
	data, err := json.MarshalIndent(this, "", "   ")
	if err != nil {
		return []byte{}
	}
	return data
}
