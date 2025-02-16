package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"reflect"
	"strconv"
)

const (
	JwtSecret                = "JWT_SECRET"
	MailjetAPIKey            = "MAILJET_API_KEY"
	MailjetSecretKey         = "MAILJET_SECRET_KEY"
	FrontendBaseUrl          = "FRONTEND_BASE_URL"
	SlackWebhookURL          = "SLACK_WEBHOOK_URL"
	EnableSlackNotifications = "ENABLE_SLACK_NOTIFICATIONS"
	DocumentBasePath         = "DOCUMENT_BASE_PATH"
	ObjectStorageEndpoint    = "S3_ENDPOINT"
	ObjectStorageRegion      = "S3_REGION"
	ObjectStorageBucket      = "S3_BUCKET_NAME"
	ObjectStorageAccessKey   = "S3_ACCESS_KEY_ID"
	ObjectStorageSecretKey   = "S3_SECRET_KEY"
)

type Config struct {
	JwtSecret                string
	MailjetAPIKey            string
	MailjetSecretKey         string
	FrontendBaseUrl          string
	SlackWebhookURL          string
	DocumentBasePath         string
	ObjectStorageEndpoint    string
	ObjectStorageRegion      string
	ObjectStorageBucket      string
	ObjectStorageAccessKey   string
	ObjectStorageSecretKey   string
	EnableSlackNotifications bool
}

func New() *Config {
	if err := godotenv.Load(); err != nil {
		log.Panic("No .env file found")
	}
	return &Config{
		JwtSecret:                getEnv[string](JwtSecret, "secret"),
		MailjetAPIKey:            getEnv[string](MailjetAPIKey, ""),
		MailjetSecretKey:         getEnv[string](MailjetSecretKey, ""),
		FrontendBaseUrl:          getEnv[string](FrontendBaseUrl, ""),
		SlackWebhookURL:          getEnv[string](SlackWebhookURL, ""),
		EnableSlackNotifications: getEnv[bool](EnableSlackNotifications, false),
		DocumentBasePath:         getEnv[string](DocumentBasePath, "test"),
		ObjectStorageEndpoint:    getEnv[string](ObjectStorageEndpoint, "http://localhost:9000"),
		ObjectStorageRegion:      getEnv[string](ObjectStorageRegion, "us-east-1"),
		ObjectStorageBucket:      getEnv[string](ObjectStorageBucket, "kariyerklubu"),
		ObjectStorageAccessKey:   getEnv[string](ObjectStorageAccessKey, "minio"),
		ObjectStorageSecretKey:   getEnv[string](ObjectStorageSecretKey, "minio123"),
	}
}

func (c *Config) GetEnvironment(key string) string {
	return getEnv(key, "")
}

func getEnv[T any](key string, defaultValue T) T {
	if value, ok := os.LookupEnv(key); ok {
		return convertToType(value, defaultValue)
	}

	return defaultValue
}

func convertToType[T any](value string, defaultValue T) T {
	var result any

	switch reflect.TypeOf(defaultValue).Kind() {
	case reflect.String:
		result = value
	case reflect.Int:
		if intValue, err := strconv.Atoi(value); err == nil {
			result = intValue
		} else {
			return defaultValue
		}
	case reflect.Float64:
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			result = floatValue
		} else {
			return defaultValue
		}
	case reflect.Bool:
		if boolValue, err := strconv.ParseBool(value); err == nil {
			result = boolValue
		} else {
			return defaultValue
		}
	default:
		return defaultValue
	}

	return result.(T)
}
