package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/lovehotel24/booking-service/pkg/configs"
	"github.com/lovehotel24/booking-service/pkg/controllers"
	"github.com/lovehotel24/booking-service/pkg/grpc/userpb"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "booking-service",
	Short: "booking service for love hotel24",
	Run:   runCommand,
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

func runCommand(cmd *cobra.Command, args []string) {
	dbConf := &configs.DBConfig{
		Host:       viper.GetString("pg-host"),
		Port:       viper.GetString("pg-port"),
		User:       viper.GetString("pg-user"),
		Pass:       viper.GetString("pg-pass"),
		DBName:     viper.GetString("pg-db"),
		SSLMode:    viper.GetString("pg-ssl"),
		AdminPhone: viper.GetString("adm-ph"),
		AdminPass:  viper.GetString("adm-pass"),
		UserPhone:  viper.GetString("usr-ph"),
		UserPass:   viper.GetString("usr-pass"),
	}

	configs.Connect(dbConf)

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	s := grpc.NewServer(opts...)
	reflection.Register(s)

	userpb.RegisterUserServiceServer(s, &controllers.UserService{})

	go func() {

		fmt.Println("Starting server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}

	}()

	// Wait for control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()

}
