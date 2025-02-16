package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env        string         `yaml:"env" env-default:"local"`
	PgConf     PostgresConfig `yaml:"postgres"`
	AmqpConf   AmqpConfig     `yaml:"amqp"`
	OtlpConfig OtlpConfig     `yaml:"otlp_config"`
	//MigrationsPath string
	//TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type PostgresConfig struct {
	//"postgres://vse:polPOL765@localhost:5432/tg-queue"
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"user_name"`
	UserPass string `yaml:"user_pass"`
	DbName   string `yaml:"db_name"`
}

type AmqpConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	UserName     string `yaml:"user_name"`
	UserPass     string `yaml:"user_pass"`
	QueueName    string `yaml:"queue"`
	ExchangeName string `yaml:"exchange"`
	RoutingKey   string `yaml:"routing_key"`
}

type OtlpConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	ServiceName string `yaml:"service_name"`
}

func (r AmqpConfig) GetAmqpUri() string {
	//"amqp://guest:guest@localhost:5672/"
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", r.UserName, r.UserPass, r.Host, r.Port)
}

func (r PostgresConfig) GetDbUri() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", r.UserName, r.UserPass, r.Host, r.Port, r.DbName)
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
