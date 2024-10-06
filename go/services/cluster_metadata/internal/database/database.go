package database

import (
	"os"

	"github.com/Magpie-Monitor/magpie-monitor/pkg/envs"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/mongodb"
)

func NewMongoDbConnectionDetails() *mongodb.MongoDbConnectionDetails {
	envs.ValidateEnvs("Failed to connect to metadatadb", []string{
		"METADATADB_USER",
		"METADATADB_PASSWORD",
		"METADATADB_HOST",
		"METADATADB_PORT",
	})

	return &mongodb.MongoDbConnectionDetails{
		User:     os.Getenv("METADATADB_USER"),
		Password: os.Getenv("METADATADB_PASSWORD"),
		Host:     os.Getenv("METADATADB_HOST"),
		Port:     os.Getenv("METADATADB_PORT"),
	}
}
