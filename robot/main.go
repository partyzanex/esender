package main

import (
	"database/sql"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/partyzanex/esender/domain"
	"github.com/partyzanex/esender/services"
	"time"
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

	emails, err := storage.List(&domain.Filter{
		Status: domain.StatusCreated,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, email := range emails {
		email := email
		err := sender.Send(email)
		if err != nil {
			email.Status = domain.StatusError
			textError := err.Error()
			email.Error = &textError
		} else {
			email.Status = domain.StatusSent
			now := time.Now()
			email.DTSent = &now
		}

		_, err = storage.Update(email)
		if err != nil {
			log.Fatal(err)
		}
	}
}
