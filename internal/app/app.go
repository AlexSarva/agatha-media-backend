package app

import (
	"AlexSarva/media/models"
	"AlexSarva/media/storage"
	"AlexSarva/media/storage/storagepg"
	"errors"
	"fmt"
)

// Database interface for different types of databases
type Database struct {
	Repo storage.Repo
}

// NewStorage generate new instance of database
func NewStorage(dbName string, cfg models.Config) (*Database, error) {
	if dbName == "PG" {
		DB := storagepg.NewPostgresDBConnection(cfg.DatabasePG)
		fmt.Println("Using ClickHouse Database")
		return &Database{
			Repo: DB,
		}, nil
	} else {
		return &Database{}, errors.New("u must use database config")
	}

}
