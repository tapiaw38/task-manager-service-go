package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigurationServiceUsesDefaults(t *testing.T) {
	tests := map[string]struct {
		prepare  func(t *testing.T)
		expected *ConfigurationService
	}{
		"when environment values are not defined": {
			prepare: setEmptyEnvironment,
			expected: &ConfigurationService{
				AppName:    defaultAppName,
				AppVersion: defaultAppVersion,
				ServerConfig: ServerConfig{
					Port:           defaultPort,
					GinMode:        ReleaseMode,
					AllowedOrigins: []string{defaultAllowedOrigins},
				},
				StoreConfig: StoreConfig{
					Driver:              JSONDriver,
					FilePath:            defaultStoreFilePath,
					FirestoreCollection: defaultFirestoreCollection,
				},
			},
		},
		"when environment values have surrounding spaces": {
			prepare: func(t *testing.T) {
				setEmptyEnvironment(t)
				t.Setenv("APP_NAME", " Task Service ")
				t.Setenv("APP_VERSION", " 2.0.0 ")
				t.Setenv("PORT", " 9090 ")
				t.Setenv("GIN_MODE", "DEBUG")
				t.Setenv("ALLOWED_ORIGINS", " https://web.example, https://gateway.example ")
				t.Setenv("STORE_DRIVER", "FIRESTORE")
				t.Setenv("STORE_FILE_PATH", " data/tasks.json ")
				t.Setenv("FIRESTORE_PROJECT_ID", " project-id ")
				t.Setenv("FIRESTORE_COLLECTION", " task-items ")
			},
			expected: &ConfigurationService{
				AppName:    "Task Service",
				AppVersion: "2.0.0",
				ServerConfig: ServerConfig{
					Port:           "9090",
					GinMode:        DebugMode,
					AllowedOrigins: []string{"https://web.example", "https://gateway.example"},
				},
				StoreConfig: StoreConfig{
					Driver:              FirestoreDriver,
					FilePath:            "data/tasks.json",
					FirestoreProjectID:  "project-id",
					FirestoreCollection: "task-items",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.prepare(t)

			assert.Equal(t, tt.expected, NewConfigurationService())
		})
	}
}

func setEmptyEnvironment(t *testing.T) {
	for _, key := range []string{
		"APP_NAME",
		"APP_VERSION",
		"PORT",
		"GIN_MODE",
		"ALLOWED_ORIGINS",
		"STORE_DRIVER",
		"STORE_FILE_PATH",
		"FIRESTORE_PROJECT_ID",
		"FIRESTORE_COLLECTION",
	} {
		t.Setenv(key, "")
	}
}
