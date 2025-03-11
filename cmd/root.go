package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"template/internal/config"

	"github.com/alecthomas/units"
	"github.com/qazwsxedckll/logh"
	"github.com/spf13/cobra"
)

var (
	cfgFile  string
	c        config.Config
	logger   *slog.Logger
	levelVar = slog.LevelVar{}
	Version  string
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
	Version: Version,
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "configs/config.toml", "config file (default is configs/config.toml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initConfig()
	initLogger()
}

func initConfig() {
	var err error
	c, err = config.Load(cfgFile)
	if err != nil {
		panic(fmt.Sprintf("error loading config: %v", err))
	}
}

func initLogger() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := levelVar.UnmarshalText([]byte(c.Log.Level))
	if err != nil {
		logger.Warn("invalid log level, use INFO instead")
		levelVar.Set(slog.LevelInfo)
	}

	var handler slog.Handler
	if c.Log.ToFile {
		size, err := units.ParseStrictBytes(c.Log.RotateSize)
		if err != nil {
			logger.Warn("invalid log rotate size, use 100MB instead", "err", err, "size", c.Log.RotateSize)
			size = int64(100 * units.MB)
		}
		interval, err := time.ParseDuration(c.Log.RotateInterval)
		if err != nil {
			logger.Warn("invalid log rotate interval, use 24h instead", "err", err, "interval", c.Log.RotateInterval)
			interval = 24 * time.Hour
		}

		options := []logh.Option{
			logh.WithRotateInterval(interval),
		}
		if c.Log.RotateAtMidnight {
			options = append(options, logh.WithRotateAtMidnight())
		}

		handler, err = logh.NewRotateJSONHandler(c.Log.Directory, c.Log.BaseName, int(size), &slog.HandlerOptions{
			AddSource: c.Log.AddSource,
			Level:     &levelVar,
		}, options...)
		if err != nil {
			panic(fmt.Sprintf("error creating rotate handler: %v", err))
		}
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: c.Log.AddSource,
			Level:     &levelVar,
		})
	}

	logger = slog.New(handler)

	err = config.Watch(cfgFile, &levelVar, logger)
	if err != nil {
		panic(fmt.Sprintf("error watching config: %v", err))
	}
}
