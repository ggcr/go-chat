package ctrl

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "chat_app"
)

type Conn struct {
	Sql *sql.DB
}

func NewConnection() *Conn {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return &Conn{db}
}

func (conn *Conn) CloseConnection() {
	err := conn.Sql.Close()
	if err != nil {
		panic(err)
	}
}
