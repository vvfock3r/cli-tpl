package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"cli-tpl/cmd/config"
	"cli-tpl/pkg/logger"
	viperutil "cli-tpl/pkg/viper"
)

const VERSION = "v0.0.1"

var (
	// cli flags
	version    bool
	configFile string
	logLevel   string
	logFormat  string
	logOutput  string
)

const (
	// viper bind keys
	BindLogLevel  = "global.log.level"
	BindLogFormat = "global.log.format"
	BindLogOutput = "global.log.output"
)

var rootCmd = &cobra.Command{
	Use:           "cli-tpl",
	Short:         "Simple Command-Line Interface Template\nFor details, please refer to https://github.com/vvfock3r/cli-tpl",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		// -v/--version
		if version {
			fmt.Println(VERSION)
			os.Exit(0)
		}

		// -c / --configFile
		viper.SupportedExts = []string{"yaml"}
		if configFile != "" {
			viper.SetConfigFile(configFile)
			err := viper.ReadInConfig()
			if err != nil {
				return err
			}
		}

		// --log-xxx
		err = setDefaultLogger()
		if err != nil {
			return err
		}

		// viper watch config
		if configFile != "" {
			viperutil.RegisterWatchFunc("global.log", setDefaultLogger)
			viperutil.StartWatchConfig()
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		for i := 0; i < 10000; i++ {
			logger.Info("root command run")
			time.Sleep(time.Second)
		}
	},
}

func init() {
	// -h / --help
	// -v / --version
	rootCmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		if command.Long != "" {
			fmt.Printf("%s\n\n", command.Long)
		} else {
			fmt.Printf("%s\n\n", command.Short)
		}
		fmt.Printf("%s", command.UsageString())
		os.Exit(0)
	})
	rootCmd.SetHelpCommand(&cobra.Command{Use: "no-help", Hidden: true})
	rootCmd.PersistentFlags().BoolP("help", "h", false, "help message")
	rootCmd.PersistentFlags().BoolVarP(&version, "version", "v", false, "version message")

	// -c / --configFile
	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "c", "", "config file")

	// --log-xxx
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "", "info", "log level")
	rootCmd.PersistentFlags().StringVarP(&logFormat, "log-format", "", "console", "log format")
	rootCmd.PersistentFlags().StringVarP(&logOutput, "log-output", "", "stdout", "log output")

	// viper bind
	var exitIfHasError = func(err error) {
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}
	exitIfHasError(viper.BindPFlag(BindLogLevel, rootCmd.PersistentFlags().Lookup("log-level")))
	exitIfHasError(viper.BindPFlag(BindLogFormat, rootCmd.PersistentFlags().Lookup("log-format")))
	exitIfHasError(viper.BindPFlag(BindLogOutput, rootCmd.PersistentFlags().Lookup("log-output")))

	// sub commands
	rootCmd.AddCommand(config.ConfigCmd)
}

// setDefaultLogger 替换默认的Logger
func setDefaultLogger() error {
	// --log-level
	level, err := logger.NewAtomicLevel(viper.GetString(BindLogLevel))
	if err != nil {
		return err
	}

	// --log-format
	encoder, err := logger.NewEndocer(viper.GetString(BindLogFormat))
	if err != nil {
		return err
	}

	// --log-output
	syncers, err := logger.NewOutput(strings.Split(viper.GetString(BindLogOutput), ","))
	if err != nil {
		return err
	}

	// Replace the default Logger
	logger.SetDefaultLogger(logger.NewLogger(level, encoder, syncers))

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
