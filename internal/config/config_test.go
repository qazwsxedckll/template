package config

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		cfgFile  string
		expected Config
		err      bool
	}{
		{
			name:    "invalid",
			cfgFile: "testdata/invalid.toml",
			err:     true,
		},
		{
			name:     "empty",
			cfgFile:  "testdata/empty.toml",
			expected: DefaultConfig,
		},
		{
			name:    "valid",
			cfgFile: "testdata/config.toml",
			expected: Config{
				Log: Log{
					Level:            "warn",
					AddSource:        true,
					ToFile:           false,
					BaseName:         "testdata",
					Directory:        "testlog",
					RotateSize:       "10MB",
					RotateInterval:   "10m",
					RotateAtMidnight: false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := Load(tt.cfgFile)
			if tt.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, config)
			}
		})
	}
}

func TestLoadEnv(t *testing.T) {
	err := os.Setenv("TEST_LOG__LEVEL", "error")
	require.NoError(t, err)

	c, err := Load("testdata/config.toml")
	require.NoError(t, err)

	require.Equal(t, "error", c.Log.Level)
}

func TestWatch(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	file, err := os.CreateTemp("", "config.toml")
	require.NoError(t, err)
	defer os.Remove(file.Name())

	levelVar := slog.LevelVar{}
	err = Watch(file.Name(), &levelVar, logger)
	require.NoError(t, err)

	_, err = file.Write([]byte(`[log]
level = "error"`))
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		return levelVar.Level() == slog.LevelError
	}, 1*time.Second, 10*time.Millisecond)
}
