package main

import (
	"AlexSarva/media/admin"
	"AlexSarva/media/internal/app"
	"AlexSarva/media/models"
	"AlexSarva/media/server"
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

func main() {
	var cfg models.Config
	// Priority on flags
	// Load config from env
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Rewrite from start parameters
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "host:port to listen on")
	flag.StringVar(&cfg.DatabasePG, "dbpg", cfg.DatabasePG, "postgresql database config")
	flag.StringVar(&cfg.DatabaseClick, "dbclick", cfg.DatabaseClick, "clickhouse database config")
	flag.Parse()
	log.Printf("%+v\n", cfg)
	log.Printf("ServerAddress: %v", cfg.ServerAddress)
	workDB, dbErr := app.NewStorage("PG", cfg)
	if dbErr != nil {
		log.Fatal(dbErr.Error() + "говно")
	}
	adminPG := admin.NewAdminDBConnection(cfg.DatabasePG)
	ping := workDB.Repo.Ping()
	log.Println(ping)
	MainApp := server.NewServer(&cfg, workDB, adminPG)
	if runErr := MainApp.Run(); runErr != nil {
		log.Printf("%s", runErr.Error())
	}
}
