package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/Ssawa/destiny/storage"
	"github.com/Ssawa/destiny/utils"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultHome = "~/.destiny"

var expandedDefaultHome string
var cfgFile string

var useVerbose bool

var rootTags []string
var rootExcludes []string
var rootExclusive bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "destiny",
	Short: "Manage your adage",
	Long: `Display, store, and manage a database of your favorite personal
adages.

Like 'fortune' but more complicated.

When executed without a subcommand it will print a random adage from the
database.
`,

	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize random seed
		rand.Seed(time.Now().Unix())

		// Open the database in Read Only mode (so that we don't lock the file)
		utils.Verbose.Println("Opening database...")
		db, err := utils.OpenReadOnly(viper.GetString("database"))
		if err != nil {
			return err
		}

		var adage *storage.Adage
		if len(rootTags) == 0 {
			adage, err = storage.GetAdageFromAll(db)
		} else {
			adage, err = storage.GetAdageFromCategories(db, rootTags, rootExcludes, rootExclusive)
		}

		if err != nil {
			return err
		}

		if adage != nil {
			fmt.Println(adage.Body)
		}
		return nil
	},

	// Don't show usage when the Run function returns an error
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogger, initConfig)

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.destiny.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&useVerbose, "verbose", "v", false, "Print verbose output")

	RootCmd.Flags().StringArrayVarP(&rootTags, "tag", "t", nil, "A tag to search for in the query")
	RootCmd.Flags().StringArrayVarP(&rootExcludes, "exclude", "x", nil, "What tags should be excluded from the result set?")
	RootCmd.Flags().BoolVarP(&rootExclusive, "exclusive", "e", false, "Should we use an AND query for tags? (defaults to OR)")

	// Configuration options
	var err error
	expandedDefaultHome, err = homedir.Expand(defaultHome)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.SetDefault("database", filepath.Join(expandedDefaultHome, "/destiny.db"))
}

// initLogger initializes the verbose logger
func initLogger() {
	if useVerbose {
		utils.Verbose = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".destiny" (without extension).
		viper.AddConfigPath(expandedDefaultHome)
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
