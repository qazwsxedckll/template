package config

import (
	"testing"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	k := koanf.New(".")
	f := file.Provider("testdata/config.toml")
	err := k.Load(f, toml.Parser())
	require.NoError(t, err)

	var c Config
	err = k.Unmarshal("", &c)
	require.NoError(t, err)

	require.Equal(t, "debug", c.Log.Level)
	require.Equal(t, true, c.Log.AddSource)
	require.Equal(t, true, c.Log.ToFile)
}
