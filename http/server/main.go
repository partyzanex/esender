package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/partyzanex/esender/http/server/handler"
	"github.com/partyzanex/esender/storage"
	"github.com/partyzanex/esender/sender"
	"github.com/partyzanex/esender/domain"
)

var (
	configPath = pflag.String("config", "config.yaml", "Configuration file path")
)

func main() {
	pflag.Parse()

	config := viper.New()
	config.SetConfigType("yaml")

	if *configPath != "" {
		configFile, err := os.Open(*configPath)
		if err != nil {
			log.Fatalf("opening config failed: %s", err)
		}

		err = config.ReadConfig(configFile)
		if err != nil {
			log.Fatalf("reading config failed: %s", err)
		}
	}

	store, err := storage.Create(storage.Config{
		Name: config.GetString("storage.name"),
		DSN:  config.GetString("storage.dsn"),
	})
	if err != nil {
		log.Fatalf("creating storage failed: %s", err)
	}

	senders, err := sender.Create(config.Get("senders").([]interface{}))
	if err != nil {
		log.Fatalf("creating senders failed: %s", err)
	}

	for _, emailSender := range senders.All() {
		go func() {
			agent := domain.NewAgent(store, emailSender)
			agent.Run()
		}()
	}

	h := &handler.Handler{
		Senders: senders,
		Storage: store,
	}

	e := echo.New()

	e.GET("/emails", h.CreateEmail)
	e.GET("/emails/:id", h.CreateEmail)
	e.POST("/emails", h.CreateEmail)
	e.PUT("/emails", h.CreateEmail)
	e.POST("/emails/send", h.SendEmail)

	addr := config.GetString("http.host") + ":" + config.GetString("http.port")

	e.Logger.Fatal(e.Start(addr))
}
