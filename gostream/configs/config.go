package configs

import "os"

type Config struct {
	ServerPort string
	JWTSecret  string
}

func LoadConfig() *Config {
	return &Config{
		ServerPort: getEnvOrDefault("SERVER_PORT", "8080"),
		JWTSecret:  getEnvOrDefault("JWT_SECRET", "your-secret-key"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
