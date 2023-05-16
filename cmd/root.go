package cmd

import (
	"os"
	"path/filepath"

	"fmt"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xls2ins",
	Short: "Generate MySql Insert SQL from xls",
	Long: `Generate MySql Insert SQL from defined in xls.

Define table cell in application config file(.xls2ins)

## Sample .xls2ins

Entries:
  - BookNameRegExp: test.*.xlsx
    Tables:
      - SheetName: サンプル
        StartRow: 3
        SqlFormat: "insert into sample values('%s', '%s', '%s', '%s');"
        SqlArgCols:
          - C
          - D
          - F
          - G
      - SheetName: サンプル2
        StartRow: 4
        SqlFormat: "insert into sample2 values('%s', '%s', '%s', '%s');"
        SqlArgCols:
          - D
          - E
          - G
          - H
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Error("Specify xls file")
			os.Exit(1)
		}
		err := Gen(config, args)
		if err != nil {
			log.Error("Failed to gen ddl", err)
			os.Exit(1)
		}
		log.Info("done")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgFile string
var config Config

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.xls2ins)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	err := initConfigInner()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfigInner() error {
	level, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err == nil {
		log.SetLevel(level)
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		// Search config in home directory with name ".xls2ins" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		// Add config path of executable directory
		ex, err := os.Executable()
		if err != nil {
			return err
		}
		exPath := filepath.Dir(ex)
		viper.AddConfigPath(exPath)

		// viper.AddConfigPath("..")
		viper.SetConfigName(".xls2ins")
	}
	viper.SetConfigType("yml")

	viper.AutomaticEnv() // read in environment variables that match

	log.Debug(">>>>>> os.Args", os.Args)
	// If a config file is found, read it in.
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	log.Debugf("> Loading config from %s.\n", viper.ConfigFileUsed())
	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}
	log.Debugf("> Loaded config: %#+v", config)
	return nil
}
