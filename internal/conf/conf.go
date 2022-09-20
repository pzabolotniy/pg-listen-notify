package conf

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type App struct {
	DB     *DB     `json:"db" yaml:"db" mapstructure:"db"`
	WebAPI *WebAPI `json:"webapi" yaml:"webapi" mapstructure:"webapi"`
	Events *Events `json:"events" yaml:"events" mapstructure:"events"`
}

type DB struct {
	ConnString      string        `json:"conn_string" yaml:"conn_string" mapstructure:"conn_string"`
	MigrationDir    string        `json:"migration_dir" yaml:"migration_dir" mapstructure:"migration_dir"`
	MigrationTable  string        `json:"migration_table" yaml:"migration_table" mapstructure:"migration_table"`
	MaxOpenConns    int32         `json:"max_open_conns" yaml:"max_open_conns" mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
}

type WebAPI struct {
	Listen string `json:"listen" yaml:"listen" mapstructure:"listen"`
}

type Events struct {
	ChannelName  string `json:"channel_name" yaml:"channel_name" mapstructure:"channel_name"`
	WorkersCount int    `json:"workers_count" yaml:"workers_count" mapstructure:"workers_count"`
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
