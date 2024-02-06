package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"template/internal/config"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/cobra"
)

var (
	cfgFile  string
	k        = koanf.New(".")
	c        config.Config
	logger   *slog.Logger
	levelVar = slog.LevelVar{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "template",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLogger)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "configs/config.toml", "config file (default is configs/config.toml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	f := file.Provider(cfgFile)
	err := k.Load(f, toml.Parser())
	if err != nil {
		logger.Error("error loading config", "err", err)
	}

	c = config.DefaultConfig
	err = k.Unmarshal("", &c)
	if err != nil {
		logger.Error("error unmarshalling config", "err", err)
	}

	logger.Info("config", "config", c)

	err = f.Watch(func(event interface{}, err error) {
		if err != nil {
			logger.Error("watch error", "err", err)
			return
		}

		logger.Info("config changed. Reloading ...")
		k = koanf.New(".")
		if err := k.Load(f, toml.Parser()); err != nil {
			logger.Error("error loading config", "err", err)
			return
		}
		logger.Info("config", "config", k.Raw())

		err = levelVar.UnmarshalText(k.Bytes("log.level"))
		if err != nil {
			levelVar.Set(slog.LevelInfo)
			logger.Warn("invalid log level, use info instead")
		}
	})
	if err != nil {
		logger.Error("error watching file", "err", err)
	}
}

func initLogger() {
	err := levelVar.UnmarshalText([]byte(c.Log.Level))
	if err != nil {
		logger.Info("invalid log level, use INFO instead")
		levelVar.Set(slog.LevelInfo)
	}

	var writers []io.Writer
	if c.Log.ToFile {
		file, err := newLogFile()
		if err != nil {
			panic(fmt.Errorf("cannot create log file: %w", err))
		}
		writers = append(writers, file)
	}
	writers = append(writers, os.Stdout)
	mw := io.MultiWriter(writers...)

	logger = slog.New(slog.NewJSONHandler(mw, &slog.HandlerOptions{
		AddSource: c.Log.AddSource,
		Level:     &levelVar,
	}))
}

func newLogFile() (io.Writer, error) {
	if err := os.MkdirAll("log", 0o744); err != nil {
		return nil, fmt.Errorf("cannot create log directory: %w", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("cannot get hostname: %w", err)
	}

	file, err := os.OpenFile("log"+string(os.PathSeparator)+k.String("log.base_name")+
		"."+time.Now().Format("20060102-150405")+"."+hostname+"."+fmt.Sprint(os.Getpid())+".json",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}

	return file, nil
}
