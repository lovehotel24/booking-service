package cmd

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
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
	rootCmd.Flags().Bool("pg-ssl", false, "postgres server ssl mode on or not")
	rootCmd.Flags().String("port", "8081", "booking service port")
	rootCmd.Flags().String("grpc-port", "50051", "booking service grpc port")
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
	dbConf := configs.NewDBConfig().
		WithHost(viper.GetString("pg-host")).
		WithPort(viper.GetString("pg-port")).
		WithUser(viper.GetString("pg-user")).
		WithPass(viper.GetString("pg-pass")).
		WithName(viper.GetString("pg-db")).
		WithSecure(viper.GetBool("pg-ssl"))

	var log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	log.SetLevel(logrus.InfoLevel)

	db, err := configs.NewDB(dbConf)
	if err != nil {
		log.WithError(err).Error("failed to connect db")
	}

	err = configs.Migrate(db)
	if err != nil {
		log.WithError(err).Error("failed to migrate db schema")
	}

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%s", viper.GetString("grpc-port")))
	if err != nil {
		fmt.Printf("Failed to create gRPC listener: %v\n", err)
		return
	}

	var opts []grpc.ServerOption

	s := grpc.NewServer(opts...)
	reflection.Register(s)

	userService := controllers.NewUserService(db, log)
	userpb.RegisterUserServiceServer(s, userService)

	app := controllers.NewApp(db, log)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Infof("gRPC server listening on :%s", viper.GetString("grpc-port"))
		if err := s.Serve(grpcListener); err != nil {
			log.WithError(err).Error("gRPC server failed to serve")
		}
	}()

	go func() {
		defer wg.Done()
		log.Infof("Fiber server listening on :%s", viper.GetString("port"))
		if err := app.Listen(fmt.Sprintf(":%s", viper.GetString("port"))); err != nil {
			log.WithError(err).Error("fiber server failed to serve")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		log.Info("Shutting down...")

		s.GracefulStop()

		if err := grpcListener.Close(); err != nil {
			log.WithError(err).Error("gRPC server failed to shutdown")
		}

		if err := app.Shutdown(); err != nil {
			log.WithError(err).Error("Fiber server failed to shutdown")
		}
	}

	wg.Wait()
	log.Info("Server gracefully stopped.")
}
