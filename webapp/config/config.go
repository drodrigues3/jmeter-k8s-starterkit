package config

import "github.com/ilyakaznacheev/cleanenv"

type Server struct {
	Host string `env:"HOST" env-default:"127.0.0.1"`
	Port string `env:"PORT" env-default:"8080"`
}

type InfluxConn struct {
	Host     string `env:"INFLUX_HOST" env-default:"localhost"`
	Port     string `env:"INFLUX_PORT" env-default:"8086"`
	DB       string `env:"INFLUX_DB" env-default:"telegraf"`
	User     string `env:"INFLUX_USER" env-default:"user"`
	Password string `env:"INFLUXDB_USER_PASSWORD"`
	Token    string `env:"INFLUX_TOKEN"`
}

type Config struct {
	Server     Server
	InfluxConn InfluxConn
}

func LoadConfiguration() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
