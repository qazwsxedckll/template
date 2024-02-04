package config

var DefaultConfig = Config{
	Log: Log{
		Level:     "info",
		AddSource: false,
		ToFile:    false,
	},
}

type Config struct {
	Log Log `koanf:"log"`
}

type Log struct {
	Level     string `koanf:"level"`
	AddSource bool   `koanf:"add_source"`
	ToFile    bool   `koanf:"to_file"`
	BaseName  string `koanf:"base_name"`
}
