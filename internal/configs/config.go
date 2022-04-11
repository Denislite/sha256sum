package configs

import (
	"github.com/spf13/viper"
	"sha256sum/internal/model"
)

func ParseConfigFile(dir string) (*model.PostgresConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg *model.PostgresConfig
	err = viper.Unmarshal(&cfg)

	return cfg, nil
}
