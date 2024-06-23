package config

import (
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

const (
	ProductionEnv = "production"

	DatabaseTimeout    = 5 * time.Second
	ProductCachingTime = 1 * time.Minute
	AddressCachingTime = 1 * time.Minute
	UsersCachingTime   = 1 * time.Minute
)

var AuthIgnoreMethods = []string{
	"/user.UserService/Login",
	"/user.UserService/Register",
}

type Schema struct {
	Environment            string `env:"environment"`
	HttpPort               int    `env:"http_port"`
	GrpcPort               int    `env:"grpc_port"`
	AuthSecret             string `env:"auth_secret"`
	DatabaseURI            string `env:"database_uri"`
	RedisURI               string `env:"redis_uri"`
	RedisPassword          string `env:"redis_password"`
	RedisDB                int    `env:"redis_db"`
	GOOGLE_CLIENT_ID       string `env:"google_client_id"`
	GOOGLE_CLIENT_SECRET   string `env:"google_client_secret"`
	GOOGLE_REDIRECT_URL    string `env:"google_redirect_url"`
	FACEBOOK_CLIENT_ID     string `env:"facebook_client_id"`
	FACEBOOK_CLIENT_SECRET string `env:"facebook_client_secret"`
	FACEBOOK_REDIRECT_URL  string `env:"facebook_redirect_url"`
}

var (
	cfg Schema
)

func LoadConfig() *Schema {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)

	err := godotenv.Load(filepath.Join(currentDir, "config.yaml"))
	if err != nil {
		log.Printf("Error on load configuration file, error: %v", err)
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error on parsing configuration file, error: %v", err)
	}

	return &cfg
}

func GetConfig() *Schema {
	return &cfg
}
