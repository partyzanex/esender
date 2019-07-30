package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/oklog/oklog/pkg/group"
	"github.com/partyzanex/esender/domain"
	"github.com/partyzanex/esender/grpc/server/endpoint"
	"github.com/partyzanex/esender/grpc/server/pb"
	"github.com/partyzanex/esender/grpc/server/transport"
	"github.com/partyzanex/esender/sender"
	"github.com/partyzanex/esender/storage"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"os/signal"
	"syscall"
)

var (
	configPath   = pflag.String("config", "config.yaml", "Configuration file path")
	connLifetime = pflag.Duration("conn_lifetime", time.Second, "postgres connection lifetime")
)

func main() {
	pflag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	config := viper.New()
	config.SetConfigType("yaml")

	if *configPath != "" {
		configFile, err := os.Open(*configPath)
		if err != nil {
			logger.Log("opening config failed:", err)
			return
		}

		err = config.ReadConfig(configFile)
		if err != nil {
			logger.Log("reading config failed:", err)
			return
		}
	}

	config.BindEnv("database.driver", "DATABASE_DRIVER")
	config.BindEnv("database.dsn", "DATABASE_DSN")
	config.BindEnv("grpc.host", "GRPC_HOST")
	config.BindEnv("grpc.port", "GRPC_PORT")

	store, err := storage.Create(storage.Config{
		Name:         config.GetString("storage.name"),
		DSN:          config.GetString("storage.dsn"),
		ConnLifetime: *connLifetime,
	})
	if err != nil {
		logger.Log("creating storage failed:", err)
		return
	}

	senders, err := sender.Create(config.Get("senders").([]interface{}))
	if err != nil {
		logger.Log("creating senders failed:", err)
		return
	}

	for _, emailSender := range senders.All() {
		go func() {
			agent := domain.NewAgent(store, emailSender)
			agent.Run()
		}()
	}

	emailSender, ok := senders.Get("")
	if !ok {
		logger.Log("no default", "sender")
		return
	}

	svc := domain.NewAgent(store, emailSender)

	grpcServer := transport.NewGRPCServer(
		endpoint.New(svc, logger), logger,
	)

	addr := fmt.Sprintf("%s:%d",
		config.GetString("grpc.host"),
		config.GetInt("grpc.port"),
	)

	var g group.Group
	{
		conn, err := net.Listen("tcp", addr)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "gRPC", "addr", addr)
			server := grpc.NewServer()
			pb.RegisterEmailServer(server, grpcServer)
			return server.Serve(conn)
		}, func(err error) {
			if err != nil {
				logger.Log(err)
			}
			conn.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())
}
