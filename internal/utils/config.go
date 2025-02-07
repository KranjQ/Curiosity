package utils

import (
	"os"
)

func GetConnectUrl() (string, error) {
	user := os.Getenv("POSTGRES_USER")
	dbname := os.Getenv("POSTGRES_NAME")
	port := os.Getenv("POSTGRES_PORT")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("DB_HOST")
	sslmode := os.Getenv("POSTGRES_SSLMODE")
	connStr := "host=" + host + " port=" + port + " user=" + user + " password=" + password +
		" dbname=" + dbname + " sslmode=" + sslmode
	return connStr, nil
}
