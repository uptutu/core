package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadConfig(t *testing.T) {
	InitConfig("../../config.yml")
	cfg := Get()
	assert.Equal(t, "core", cfg.Server.AppID)
}
