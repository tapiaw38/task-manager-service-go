package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfigServiceKeepsTheFirstValue(t *testing.T) {
	configService = nil
	t.Cleanup(func() { configService = nil })

	first := &ConfigurationService{AppName: "First"}
	second := &ConfigurationService{AppName: "Second"}

	InitConfigService(first)
	InitConfigService(second)

	assert.Equal(t, "First", GetConfigService().AppName)
}

func TestGetConfigServiceReturnsACopy(t *testing.T) {
	configService = nil
	t.Cleanup(func() { configService = nil })

	InitConfigService(&ConfigurationService{
		AppName:      "Task Manager Service",
		ServerConfig: ServerConfig{Port: "8080"},
	})

	copied := GetConfigService()
	copied.AppName = "Mutated outside"

	assert.Equal(t, "Task Manager Service", GetConfigService().AppName)
	assert.Equal(t, "8080", GetConfigService().ServerConfig.Port)
}
