package config

const (
	ReleaseMode GinMode = "release"
	DebugMode   GinMode = "debug"

	JSONDriver      StoreDriver = "json"
	FirestoreDriver StoreDriver = "firestore"
)

type (
	GinMode string

	StoreDriver string

	ConfigurationService struct {
		AppName      string
		AppVersion   string
		ServerConfig ServerConfig
		StoreConfig  StoreConfig
	}

	ServerConfig struct {
		GinMode        GinMode
		Port           string
		AllowedOrigins []string
	}

	StoreConfig struct {
		Driver              StoreDriver
		FilePath            string
		FirestoreProjectID  string
		FirestoreCollection string
	}
)

var configService *ConfigurationService

func InitConfigService(config *ConfigurationService) {
	if configService == nil {
		configService = config
	}
}

func GetConfigService() ConfigurationService {
	return *configService
}
