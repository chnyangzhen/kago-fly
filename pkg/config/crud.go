package config

import (
	"github.com/chnyangzhen/kago-fly/pkg/constant"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"os"
	"strings"
	"time"
)

func GetString(key string) string {
	return root.GetString(key)
}

func GetStringWithDefault(key string, defaultValue string) string {
	v := root.GetString(key)
	if v == "" {
		return defaultValue
	}
	return v
}

// MakeKey 根据Key列表，构建Configuration的查询Key。
// Note: Key列表任意单个Key不允许为空字符。
func MakeKey(keys ...string) string {
	if len(keys) == 0 {
		return ""
	}
	for _, key := range keys {
		if key == "" {
			panic("configuration: not allow empty key")
		}
	}
	return strings.Join(keys, ".")
}

func (c *Configuration) makeKey(key string) string {
	return MakeKey(c.namespace, key)
}

// ToStringMap 将当前配置实例（命名空间）下所有配置，转换成 map[string]any 类型的字典。
func (c *Configuration) ToStringMap() map[string]interface{} {
	if "" == c.namespace {
		return c.root.AllSettings()
	}
	return cast.ToStringMap(c.root.Get(c.namespace))
}

// Keys 获取当前配置实例（命名空间）下所有配置的键列表
func (c *Configuration) Keys() []string {
	v := c.root.Sub(c.namespace)
	if v != nil {
		return v.AllKeys()
	}
	return []string{}
}

func (c *Configuration) Get(key string) interface{} {
	return c.doGet(c.makeKey(key), nil)
}

func (c *Configuration) GetOrDefault(key string, def interface{}) interface{} {
	return c.doGet(c.makeKey(key), def)
}

// Set 向当前配置实例以覆盖的方式设置Key-Value键值。
func (c *Configuration) Set(key string, value interface{}) {
	c.root.Set(c.makeKey(key), value)
}

// SetKeyAlias 设置当前配置实例的Key与GlobalAlias的映射
func (c *Configuration) SetKeyAlias(keyAlias map[string]string) {
	// 兼容逻辑（读yml数组配置时候，alias默认转换为nil，这里做兼容）
	if c.alias == nil {
		c.alias = make(map[string]string)
	}

	for key, alias := range keyAlias {
		c.alias[c.makeKey(key)] = alias
	}
}

// SetDefault 为当前配置实例设置单个默认值。与Viper的SetDefault一致，作用于当前配置实例。
func (c *Configuration) SetDefault(key string, value interface{}) {
	c.root.SetDefault(c.makeKey(key), value)
}

// SetDefaults 为当前配置实例设置一组默认值。与Viper的SetDefault一致，作用于当前配置实例。
func (c *Configuration) SetDefaults(defaults map[string]interface{}) {
	for key, val := range defaults {
		c.root.SetDefault(c.makeKey(key), val)
	}
}

// IsSet 判定当前配置实例是否设置指定Key（多个）。与Viper的IsSet一致，查询范围为当前配置实例。
func (c *Configuration) IsSet(keys ...string) bool {
	if len(keys) == 0 {
		return false
	}
	// Any not set, return false
	for _, key := range keys {
		if !c.root.IsSet(c.makeKey(key)) {
			return false
		}
	}
	return true
}

// GetString returns the value associated with the key as a string.
func (c *Configuration) GetString(key string) string {
	return cast.ToString(c.Get(key))
}

// GetBool returns the value associated with the key as a boolean.
func (c *Configuration) GetBool(key string) bool {
	return cast.ToBool(c.Get(key))
}

// GetInt returns the value associated with the key as an integer.
func (c *Configuration) GetInt(key string) int {
	return cast.ToInt(c.Get(key))
}

// GetInt32 returns the value associated with the key as an integer.
func (c *Configuration) GetInt32(key string) int32 {
	return cast.ToInt32(c.Get(key))
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Configuration) GetInt64(key string) int64 {
	return cast.ToInt64(c.Get(key))
}

// GetUint returns the value associated with the key as an unsigned integer.
func (c *Configuration) GetUint(key string) uint {
	return cast.ToUint(c.Get(key))
}

// GetUint32 returns the value associated with the key as an unsigned integer.
func (c *Configuration) GetUint32(key string) uint32 {
	return cast.ToUint32(c.Get(key))
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func (c *Configuration) GetUint64(key string) uint64 {
	return cast.ToUint64(c.Get(key))
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Configuration) GetFloat64(key string) float64 {
	return cast.ToFloat64(c.Get(key))
}

// GetTime returns the value associated with the key as time.
func (c *Configuration) GetTime(key string) time.Time {
	return cast.ToTime(c.Get(key))
}

// GetDuration returns the value associated with the key as a duration.
func (c *Configuration) GetDuration(key string) time.Duration {
	return cast.ToDuration(c.Get(key))
}

// GetIntSlice returns the value associated with the key as a slice of int values.
func (c *Configuration) GetIntSlice(key string) []int {
	return cast.ToIntSlice(c.Get(key))
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Configuration) GetStringSlice(key string) []string {
	return cast.ToStringSlice(c.Get(key))
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Configuration) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(c.Get(key))
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Configuration) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(c.Get(key))
}

func (c *Configuration) GetStruct(key string, outptr interface{}) error {
	return c.GetStructTag(key, "json", outptr)
}

func (c *Configuration) GetStructTag(key, structTag string, outptr interface{}) error {
	key = c.makeKey(key)
	if !c.root.IsSet(key) {
		return nil
	}
	return c.root.UnmarshalKey(key, outptr, func(opt *mapstructure.DecoderConfig) {
		opt.TagName = structTag
	})
}

func (c *Configuration) doGet(key string, indef interface{}) interface{} {
	val := c.root.Get(key)
	if expr, ok := val.(string); ok {
		pkey, pdef, ptype := ParseDynamicKey(expr)
		var usedef interface{}
		if indef != nil {
			usedef = indef
		} else {
			usedef = pdef
		}
		switch ptype {
		case constant.DynamicTypeLookupConfig:
			// check circle key
			if key == pkey {
				return usedef
			}
			if c.root.IsSet(pkey) {
				return c.doGet(pkey, usedef)
			} else {
				return usedef
			}

		case constant.DynamicTypeLookupEnv:
			if ev, ok := os.LookupEnv(pkey); ok {
				return ev
			} else {
				return usedef
			}

		case constant.DynamicTypeStaticValue:
			return val

		default:
			return val
		}
	}
	// check local alias
	if nil == val {
		if alias, ok := c.alias[key]; ok {
			val = c.root.Get(alias)
		}
	}
	if nil == val {
		return indef
	}
	return val
}

// ParseDynamicKey 解析动态值：配置参数：${key:defaultV}，环境变量：#{key:defaultV}
func ParseDynamicKey(pattern string) (key string, def string, typ int) {
	pattern = strings.TrimSpace(pattern)
	size := len(pattern)
	if size <= len("${}") {
		return pattern, "", constant.DynamicTypeStaticValue
	}
	dyn := "${" == pattern[:2]
	env := "#{" == pattern[:2]
	if (dyn || env) && '}' == pattern[size-1] {
		values := strings.TrimSpace(pattern[2 : size-1])
		idx := strings.IndexByte(values, ':')
		key = values
		if idx > 0 {
			key = values[:idx]
			def = values[idx+1:]
		}
		if env {
			return key, def, constant.DynamicTypeLookupEnv
		} else {
			return key, def, constant.DynamicTypeLookupConfig
		}
	}
	return pattern, "", constant.DynamicTypeStaticValue
}
