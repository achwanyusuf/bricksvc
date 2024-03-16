package conf

import (
	"github.com/spf13/viper"
)

func New(file string) (Conf, error) {
	var cfg Conf
	viper.AddConfigPath(".")
	viper.SetConfigName(file)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return cfg, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
