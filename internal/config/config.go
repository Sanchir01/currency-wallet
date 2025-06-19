package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string      `yaml:"env"`
	Domain      string      `yaml:"domain"`
	HTTPServer  HttpServer  `yaml:"http_server"`
	GrpcClients GRPCClients `yaml:"grpc_clients"`
	RedisDB     Redis       `yaml:"redis"`
	DB          DataBase    `yaml:"database"`
	Prometheus  Prometheus  `yaml:"prometheus"`
}
type Kafka struct {
	Outbox Producer `yaml:"outbox"`
}

type Producer struct {
	Retries int      `yaml:"retries"`
	Topic   []string `yaml:"topic"`
	Broke   []string `yaml:"brokers"`
}
type Redis struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Retries  int    `yaml:"retries"`
	DBNumber int    `yaml:"dbnumber"`
}
type Prometheus struct {
	Timeout     time.Duration `yaml:"timeout"  env-default:"4s"`
	Host        string        `yaml:"host"  env-default:"localhost"`
	Port        string        `yaml:"port"  env-default:"8081"`
	Debug       bool          `yaml:"debug"  env-default:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout"  env-default:"60s"`
}
type DataBase struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	User        string `yaml:"user"`
	Database    string `yaml:"dbname"`
	SSL         string `yaml:"ssl"`
	MaxAttempts int    `yaml:"max_attempts"`
}
type GrpcOrder struct {
	Timeout int `yaml:"timeout"`
	Port    int `yaml:"port"`
}
type GRPCExchanger struct {
	Host    string        `yaml:"host"`
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
	Retries int           `yaml:"retries"`
}
type GRPCClients struct {
	GRPCExchanger GRPCExchanger `yaml:"grpc_exchanger"`
}
type S3Store struct {
	Key        string `yaml:"key"`
	BucketName string `yaml:"bucketname"`
	Region     string `yaml:"region"`
	URL        string `yaml:"url"`
}
type HttpServer struct {
	Timeout     time.Duration `yaml:"timeout"  env-default:"4s"`
	Host        string        `yaml:"host"  env-default:"localhost"`
	Port        string        `yaml:"port"  env-default:"5000"`
	Debug       bool          `yaml:"debug"  env-default:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout"  env-default:"60s"`
}

func InitConfig() *Config {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env.dev"
	}
	fmt.Println("env name", envFile)
	if err := godotenv.Load(envFile); err != nil {
		slog.Error("ошибка при инициализации переменных окружения", err.Error())
	}
	configPath := os.Getenv("CONFIG_PATH")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist:%s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	return &cfg
}
