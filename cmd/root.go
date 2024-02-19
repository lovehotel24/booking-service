package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "booking-service",
	Short: "booking service for love hotel24",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	rootCmd.Flags().String("pg-user", "postgres", "user name for postgres database")
	rootCmd.Flags().String("pg-pass", "postgres", "password for postgres database")
	rootCmd.Flags().String("pg-host", "localhost", "postgres server address")
	rootCmd.Flags().String("pg-port", "5432", "postgres server port")
	rootCmd.Flags().String("pg-db", "postgres", "postgres database name")
	rootCmd.Flags().String("pg-ssl", "disable", "postgres server ssl mode on or not")
	rootCmd.Flags().String("port", "8080", "booking service port")
	rootCmd.Flags().String("grpc-port", "8081", "booking service grpc port")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("book")
	viper.BindPFlag("pg-user", rootCmd.Flags().Lookup("pg-user"))
	viper.BindPFlag("pg-pass", rootCmd.Flags().Lookup("pg-pass"))
	viper.BindPFlag("pg-host", rootCmd.Flags().Lookup("pg-host"))
	viper.BindPFlag("pg-port", rootCmd.Flags().Lookup("pg-port"))
	viper.BindPFlag("pg-db", rootCmd.Flags().Lookup("pg-db"))
	viper.BindPFlag("pg-ssl", rootCmd.Flags().Lookup("pg-ssl"))
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("grpc-port", rootCmd.Flags().Lookup("grpc-port"))
	viper.BindEnv("gin_mode", "GIN_MODE")
	viper.AutomaticEnv()
}
