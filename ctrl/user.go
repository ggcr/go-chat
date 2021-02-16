package ctrl

import (
	"database/sql"
	"fmt"
)

func (c *Conn) AddUser(username string) int {
	sqlStatement := `INSERT INTO users (username) VALUES ($1) RETURNING id`
	id := 0
	err := c.Sql.QueryRow(sqlStatement, username).Scan(&id)
	if err != nil {
		panic(err)
	}
	return id
}

func (c *Conn) GetUserId(username interface{}) int {
	var id int
	sqlStatement := `SELECT id FROM users WHERE username=$1`

	switch err := c.Sql.QueryRow(sqlStatement, username).Scan(&id); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		return id
	default:
		panic(err)
	}
	return id
}

func (c *Conn) GetUserName(user_id interface{}) interface{} {
	var username interface{}
	sqlStatement := `SELECT username FROM users WHERE id=$1`

	switch err := c.Sql.QueryRow(sqlStatement, user_id).Scan(&username); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		return username
	default:
		panic(err)
	}
	return username
}
