package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

func GetConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg Config
	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	setFromEnv(&cfg)
	return &cfg, nil
}

func setFromEnv(cfg *Config) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	cfg.Repository.Mongo.Host = os.Getenv("MONGO_HOST")
	cfg.Repository.Mongo.Port = os.Getenv("MONGO_PORT")
	cfg.Repository.Mongo.DBName = os.Getenv("MONGO_DB_NAME")
	cfg.Repository.Mongo.User = os.Getenv("MONGO_USER")
	cfg.Repository.Mongo.Password = os.Getenv("MONGO_PASS")

	cfg.Repository.Redis.Addr = os.Getenv("REDIS_ADDR")

	cfg.Auth.PasswordSalt = os.Getenv("PASSWORD_SALT")
	cfg.Auth.JWT.SigningKey = os.Getenv("JWT_SIGNING_KEY")
	cfg.Repository.Mongo.URI = fmt.Sprintf("mongodb://%s:%s@%s:%s", cfg.Repository.Mongo.User, cfg.Repository.Mongo.Password, cfg.Repository.Mongo.Host, cfg.Repository.Mongo.Port)
}

type Config struct {
	App        *AppCfg     `yaml:"app"`
	Repository *Repository `yaml:"repository"`
	Auth       *AuthConfig `yaml:"auth"`
}

type AuthConfig struct {
	JWT          JWTConfig
	PasswordSalt string
}

type JWTConfig struct {
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl"`
	SigningKey      string
}

type Mongo struct {
	URI             string
	DBName          string
	UsersCollection string `yaml:"users_collection"`
	BooksCollection string `yaml:"books_collection"`
	User            string
	Password        string
	Host            string
	Port            string
}

type Redis struct {
	Addr string
	Ttl  time.Duration `yaml:"ttl"`
}

type CorsCfg struct {
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	ExposeHeaders    []string `yaml:"expose_headers"` // FIX if this service use outside k8s cluster
	AllowCredentials bool     `yaml:"allow_credentials"`
}

type AppCfg struct {
	Port int           `yaml:"port"`
	Cors *CorsCfg      `yaml:"cors"`
	RTO  time.Duration `yaml:"rto"`
	WTO  time.Duration `yaml:"wto"`
}

type Repository struct {
	Mongo Mongo `yaml:"mongo"`
	Redis Redis `yaml:"redis"`
}

type HTTPClientConf struct {
	ProxyURL string        `yaml:"proxy_url"`
	Timeout  time.Duration `yaml:"timeout"`
}
