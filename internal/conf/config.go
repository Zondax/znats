package conf

import (
	"github.com/spf13/viper"
)

type Config struct {
	Foo string `json:"foo"`
}

func (c Config) SetDefaults() {
	// TODO implement me
	viper.SetDefault("foo", "bar")
}

func (c Config) Validate() error {
	// TODO implement me
	return nil
}
