package utils

import "os"

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

//debug
//func GetConnectUrl() (string, error) {
//	connStr := "host=" + "localhost" + " port=" + "5432" + " user=" + "reufee" + " password=" + "curiosity" +
//		" dbname=" + "curiosityDB" + " sslmode=" + "disable"
//	return connStr, nil
//}
