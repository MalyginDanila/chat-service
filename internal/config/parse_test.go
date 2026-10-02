package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MalyginDanila/chat-service/internal/config"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.toml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestParseAndValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		cfg, err := config.ParseAndValidate(writeConfig(t, `
[global]
env = "dev"
[log]
level = "debug"
[servers.debug]
addr = ":8079"
`))
		require.NoError(t, err)
		assert.Equal(t, "debug", cfg.Log.Level)
		assert.Equal(t, ":8079", cfg.Servers.Debug.Addr)
	})

	t.Run("invalid log level", func(t *testing.T) {
		_, err := config.ParseAndValidate(writeConfig(t, `
[global]
env = "dev"
[log]
level = "trace"
[servers.debug]
addr = ":8079"
`))
		require.Error(t, err)
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := config.ParseAndValidate("unknown.toml")
		require.Error(t, err)
	})
}
