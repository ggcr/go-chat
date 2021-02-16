package ctrl

import (
	"fmt"
	"time"
)

type msg struct {
	User_id      interface{}
	Username     interface{}
	Date         string
	Body         interface{}
	Sess_user_id interface{}
}

func (c *Conn) StoreMSG(user_id interface{}, msgBody interface{}, date string) {
	sqlStatement := `INSERT INTO messages (user_id, body, date) VALUES ($1, $2, $3) RETURNING id`

	r, err := c.Sql.Exec(sqlStatement, user_id, msgBody, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		panic(err)
	}

	fmt.Println(r)
}

func (c *Conn) GetChatMsgs(s *Sess) []*msg {
	data := make([]*msg, 0)

	rows, err := c.Sql.Query("SELECT user_id, body, date FROM messages")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user_id interface{}
		var msgBody interface{}
		var Date time.Time
		err = rows.Scan(&user_id, &msgBody, &Date)
		if err != nil {
			panic(err)
		}
		// get username
		name := c.GetUserName(user_id)
		data = append(data, &msg{user_id, name, Date.Format("2006-01-02 15:04:05"), msgBody, s.GetVal("id")})
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return data

}
