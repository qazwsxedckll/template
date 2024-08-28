package config

var DefaultConfig = Config{
	Log: Log{
		Level:            "info",
		AddSource:        false,
		ToFile:           false,
		BaseName:         "main",
		Directory:        "log",
		RotateSize:       "1GB",
		RotateInterval:   "24h",
		RotateAtMidnight: true,
	},
}

type Config struct {
	Log Log `koanf:"log"`
}

type Log struct {
	Level            string `koanf:"level"`
	AddSource        bool   `koanf:"add_source"`
	ToFile           bool   `koanf:"to_file"`
	BaseName         string `koanf:"base_name"`
	Directory        string `koanf:"directory"`
	RotateSize       string `koanf:"rotate_size"`
	RotateInterval   string `koanf:"rotate_interval"`
	RotateAtMidnight bool   `koanf:"rotate_at_midnight"`
}
