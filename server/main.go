package main

import (
	"database/sql"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/partyzanex/esender/server/handler"
	"github.com/partyzanex/esender/services"
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

	db, err := sql.Open("mysql", config.GetString("mysql.dsn"))
	if err != nil {
		log.Fatalf("opening sql connection failed: %s", err)
	}

	storage := &services.EmailMySqlStorage{
		DB: db,
	}

	sender := services.NewSMTPSender(services.SMTPConfig{
		Host:     config.GetString("smtp.host"),
		Port:     uint16(config.GetInt("smtp.port")),
		UserName: config.GetString("smtp.user"),
		Password: config.GetString("smtp.password"),
	})

	h := &handler.Handler{
		Sender:  sender,
		Storage: storage,
		AuthKey: config.GetString("http.auth_key"),
	}

	e := echo.New()
	e.POST("/create", h.EmailCreate)
	e.POST("/send", h.EmailSend)

	addr := config.GetString("http.host") + ":" + config.GetString("http.port")

	e.Logger.Fatal(e.Start(addr))
}
