package conf

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var vipers *Conf

var DIR = "./conf/gateway"

type OnConfigChange func(name string, op uint32)

type Conf struct {
	*viper.Viper
}

func (c *Conf) Use(key string) *Conf {
	if c.Sub(key) == nil {
		return nil
	}
	return &Conf{
		c.Sub(key),
	}
}

func (c *Conf) MustString(key string, defaultValue string) string {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}
	value, err := cast.ToStringE(v)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *Conf) MustInt(key string, defaultValue int) int {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}
	value, err := cast.ToIntE(v)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *Conf) MustStringSlice(key string, defaultValue []string) []string {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}
	value, err := cast.ToStringSliceE(v)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *Conf) MustFloat64(key string, defaultValue float64) float64 {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}
	value, err := cast.ToFloat64E(v)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *Conf) MustInt64(key string, defaultValue int64) int64 {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}
	value, err := cast.ToInt64E(v)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *Conf) MustStringMap(key string, defaultValue map[string]interface{}) map[string]interface{} {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}
	value, err := cast.ToStringMapE(v)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *Conf) MustStringMapString(key string, defaultValue map[string]string) map[string]string {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}
	value, err := cast.ToStringMapStringE(v)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *Conf) MustUint64(key string, defaultValue uint64) uint64 {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}
	value, err := cast.ToUint64E(v)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *Conf) MustBool(key string, defaultValue bool) bool {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}
	value, err := cast.ToBoolE(v)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *Conf) MustDuration(key string, defaultValue time.Duration) time.Duration {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}
	value, err := cast.ToDurationE(v)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *Conf) NeedUse(key string) *Conf {
	if c.Sub(key) == nil {
		panic("conf Use error," + key + " is needed")
	}
	return c.Use(key)
}

//使用配置对象
func Use(key string) *Conf {
	if vipers.Sub(key) == nil {
		return nil
	}
	return &Conf{
		vipers.Sub(key),
	}
}

//需要使用，如果不存在则抛panic
func NeedUse(key string) *Conf {
	if vipers.Sub(key) == nil {
		panic("conf Use error," + key + " is needed")
	}
	return Use(key)
}

//返回该conf当前级别下所有的key
func Keys(c *Conf) []string {
	cMap := c.AllSettings()
	var keys []string
	for k := range cMap {
		keys = append(keys, k)
	}
	return keys
}

var NotifyFunc []OnConfigChange

func Register(f OnConfigChange) {
	NotifyFunc = append(NotifyFunc, f)
}

func SetConfDir(dir string) {
	DIR = dir
}

func Set(key string, value interface{}) {
	vipers.Set(key, value)
}

func Init() {
	vipers = &Conf{
		viper.New(),
	}
	vipers.SetConfigType("yaml")
	vipers.SetConfigName("gateway")
	vipers.AddConfigPath(DIR)
	err := vipers.ReadInConfig() // Find and read the conf file
	if err != nil {              // Handle errors reading the conf file
		panic("Fatal error conf file: " + err.Error())
	}
	vipers.WatchConfig()
	vipers.OnConfigChange(func(e fsnotify.Event) {
		for _, f := range NotifyFunc {
			f(e.Name, uint32(e.Op))
		}
	})
}

func InitConfigAbpath(dir, file, suffix string) {
	vipers = &Conf{
		viper.New(),
	}
	vipers.SetConfigType(suffix)
	vipers.SetConfigName(file)
	vipers.AddConfigPath(dir)
	err := vipers.ReadInConfig() // Find and read the conf file
	if err != nil {              // Handle errors reading the conf file
		panic("Fatal error conf file: " + err.Error())
	}
	vipers.WatchConfig()
	vipers.OnConfigChange(func(e fsnotify.Event) {
		for _, f := range NotifyFunc {
			f(e.Name, uint32(e.Op))
		}
	})
}
