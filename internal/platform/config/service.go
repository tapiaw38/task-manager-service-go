package config

import (
	"os"
	"strings"
)

const (
	defaultAppName             = "Task Manager Service"
	defaultAppVersion          = "1.0.0"
	defaultPort                = "8080"
	defaultStoreFilePath       = "data/bd.json"
	defaultAllowedOrigins      = "http://localhost:5173"
	defaultFirestoreCollection = "tasks"
)

func NewConfigurationService() *ConfigurationService {
	return &ConfigurationService{
		AppName:    getEnv("APP_NAME", defaultAppName),
		AppVersion: getEnv("APP_VERSION", defaultAppVersion),
		ServerConfig: ServerConfig{
			GinMode:        ginMode(getEnv("GIN_MODE", string(ReleaseMode))),
			Port:           getEnv("PORT", defaultPort),
			AllowedOrigins: splitAndTrim(getEnv("ALLOWED_ORIGINS", defaultAllowedOrigins)),
		},
		StoreConfig: StoreConfig{
			Driver:              storeDriver(getEnv("STORE_DRIVER", string(JSONDriver))),
			FilePath:            getEnv("STORE_FILE_PATH", defaultStoreFilePath),
			FirestoreProjectID:  getEnv("FIRESTORE_PROJECT_ID", ""),
			FirestoreCollection: getEnv("FIRESTORE_COLLECTION", defaultFirestoreCollection),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}

	return fallback
}

func ginMode(value string) GinMode {
	if strings.EqualFold(value, string(DebugMode)) {
		return DebugMode
	}

	return ReleaseMode
}

func storeDriver(value string) StoreDriver {
	if strings.EqualFold(value, string(FirestoreDriver)) {
		return FirestoreDriver
	}

	return JSONDriver
}

func splitAndTrim(value string) []string {
	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))

	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}

	return origins
}
