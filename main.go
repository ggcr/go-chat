package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/ggcr/gchat/ctrl"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	user_id interface{}
}

var upgrader = websocket.Upgrader{}
var tpl *template.Template
var connections = make(map[*connection]bool)

func init() {
	tpl = template.Must(template.ParseGlob("./tpl/*.html"))
}

func main() {
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/chat", ChatHandler).Methods("GET", "POST")

	r.HandleFunc("/ws", WsEndpoint)

	log.Fatalln(http.ListenAndServe(":8080", r))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Fatalln(err)
	}
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	// Init Db conn
	conn := ctrl.NewConnection()
	defer conn.CloseConnection()

	// Create Session
	Session := ctrl.CreateSess(r)

	if r.Method == "POST" {
		// Get Post Data
		if err := r.ParseForm(); err != nil {
			log.Fatalln(err)
		}

		if username := r.PostFormValue("fname"); username != "" { // User POST
			// Add user to db
			id := conn.AddUser(username)

			Session.SetVal("username", username)
			Session.SetVal("id", id)
			Session.SaveSess(w, r)

			data := conn.GetChatMsgs(Session)

			err := tpl.ExecuteTemplate(w, "chat.html", data)
			if err != nil {
				log.Fatalln(err)
			}
		} else { // Chat POST
			msgBody := r.PostFormValue("msg-body")

			// Store msg in DB
			conn.StoreMSG(Session.GetVal("id"), msgBody, time.Now().Format("01-02-2006"), Session)

			// Build chat msgs
			data := conn.GetChatMsgs(Session)

			err := tpl.ExecuteTemplate(w, "chat.html", data)
			if err != nil {
				log.Fatalln(err)
			}
		}

	} else {
		data := conn.GetChatMsgs(Session)
		err := tpl.ExecuteTemplate(w, "chat.html", data)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func reader(conn connection, s *ctrl.Sess) {
	// Init Db conn
	db := ctrl.NewConnection()
	defer db.CloseConnection()

	for {
		_, p, err := conn.ws.ReadMessage()
		if err != nil {
			panic(err)
		}
		// build msg

		// store db
		m := db.StoreMSG(s.GetVal("id"), string(p), time.Now().Format("01-02-2006"), s)

		// broadcast to all clients
		for k, v := range connections {
			if v {
				m.Sess_user_id = k.user_id
				b, err := json.Marshal(m)
				if err != nil {
					panic(err)
				}
				if err := k.ws.WriteJSON(string(b)); err != nil {
					panic(err)
				}
			}
		}
	}
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	// Create Session
	Session := ctrl.CreateSess(r)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	// connections map[*connection]bool
	fmt.Println("HELO")
	c := connection{ws, Session.GetVal("id")}
	connections[&c] = true
	fmt.Println(connections)
	reader(c, Session)
	_ = ws
}
