package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

type App struct {
	DB     *DB     `json:"db" yaml:"db" mapstructure:"db"`
	WebAPI *WebAPI `json:"webapi" yaml:"webapi" mapstructure:"webapi"`
}

type DB struct {
	ConnString     string `json:"conn_string" yaml:"conn_string" mapstructure:"conn_string"`
	MigrationDir   string `json:"migration_dir" yaml:"migration_dir" mapstructure:"migration_dir"`
	MigrationTable string `json:"migration_table" yaml:"migration_table" mapstructure:"migration_table"`
}

type WebAPI struct {
	Listen string `json:"listen" yaml:"listen" mapstructure:"listen"`
}

func GetConfig() (*App, error) {
	viper.SetConfigName("config") // hardcoded config name
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // hardcoded configfile path
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("unable to read config from file: %w", err)
	}
	viper.AutomaticEnv()

	config := new(App)
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %w", err)
	}

	return config, nil
}
