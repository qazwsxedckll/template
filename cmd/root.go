package cmd

import (
	"log"
	"log/slog"
	"os"

	"template/internal/config"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	k       = koanf.New(".")
	c       config.Config
	logger  *slog.Logger
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "configs/config.toml", "config file (default is configs/config.toml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	f := file.Provider(cfgFile)
	err := k.Load(f, toml.Parser())
	if err != nil {
		log.Printf("error loading config: %v", err)
	}

	levelVar := slog.LevelVar{}
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: k.Bool("log.add_source"),
		Level:     &levelVar,
	}))

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
			logger.Info("invalid log level, use info instead")
		}
	})
	if err != nil {
		log.Printf("error watching file: %v", err)
	}

	c = config.DefaultConfig
	err = k.Unmarshal("", &c)
	if err != nil {
		log.Printf("error unmarshalling config: %v", err)
	}
}
