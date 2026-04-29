package db

import (
	"database/sql"
	"fmt"
	"rabbitmq/utils"

	_ "github.com/lib/pq"
)

type DBConn struct {
	DB *sql.DB
}

func NewDbConn() *DBConn {
	user := utils.GetEnvString("PSQL_USER", "user123")
	password := utils.GetSecretString(
		utils.ReadSecret("/run/secrets/sql_pass"),
		"PSQL_PASS",
		"pass123",
	)
	dbName := utils.GetEnvString("PSQL_DBNAME", "inventorydb")
	hostName := utils.GetEnvString("PSQL_HOSTNAME", "127.0.0.1")
	port := utils.GetEnvString("PSQL_HOSTPORT", "5432")

	connstr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		user,
		dbName,
		password,
		hostName,
		port,
	)

	conn, err := sql.Open("postgres", connstr)
	utils.FailOnError(err, "unable to open a db connection")

	err = conn.Ping()
	utils.FailOnError(err, "unable to ping to the database")
	return &DBConn{
		DB: conn,
	}
}
