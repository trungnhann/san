package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	Host        string `mapstructure:"HOST"`
	Port        int    `mapstructure:"PORT"`

	DBUsername    string `mapstructure:"DB_USERNAME"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBHostname    string `mapstructure:"DB_HOSTNAME"`
	DBPort        int    `mapstructure:"DB_PORT"`
	DBName        string `mapstructure:"DB_DBNAME"`
	DBNameTest    string `mapstructure:"DB_DBNAME_TEST"`
	MigrationPath string `mapstructure:"MIGRATION_PATH"`
	DBRecreate    bool   `mapstructure:"DB_RECREATE"`
	DBUrl         string

	StorageEndpoint  string `mapstructure:"STORAGE_ENDPOINT"`
	StorageAccessKey string `mapstructure:"STORAGE_ACCESS_KEY"`
	StorageSecretKey string `mapstructure:"STORAGE_SECRET_KEY"`
	StorageBucket    string `mapstructure:"STORAGE_BUCKET_NAME"`
	StorageRegion    string `mapstructure:"STORAGE_REGION"`
	StorageUseSSL    bool   `mapstructure:"STORAGE_USE_SSL"`

	JWTSecret             string `mapstructure:"JWT_SECRET"`
	JWTExpirationHours    int    `mapstructure:"JWT_EXPIRATION_HOURS"`
	RefreshExpirationDays int    `mapstructure:"REFRESH_EXPIRATION_DAYS"`

	EmailSenderName     string `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress  string `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD"`
	EmailSenderHost     string `mapstructure:"EMAIL_SENDER_HOST"`
	EmailSenderPort     int    `mapstructure:"EMAIL_SENDER_PORT"`
	EmailSenderUsername string `mapstructure:"EMAIL_SENDER_USERNAME"`

	RabbitMQSource string `mapstructure:"RABBITMQ_SOURCE"`
	RedisAddress   string `mapstructure:"REDIS_ADDRESS"`
}

func LoadConfig(name string, path string) (config Config) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("config: %v", err)
		return
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("config: %v", err)
		return
	}
	return
}
