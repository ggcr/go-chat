package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/ggcr/gchat/ctrl"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type msg struct {
	user_id  interface{}
	Username interface{}
	Date     string
	Body     interface{}
}

var upgrader = websocket.Upgrader{}
var tpl *template.Template

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
			fmt.Println("USER POST")
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
			fmt.Println("CHAT POST")
			msgBody := r.PostFormValue("msg-body")

			// Store msg in DB
			conn.StoreMSG(Session.GetVal("id"), msgBody, time.Now().Format("01-02-2006"))

			// Build chat msgs
			data := conn.GetChatMsgs(Session)
			for _, v := range data {
				fmt.Println(v)
			}

			err := tpl.ExecuteTemplate(w, "chat.html", data)
			if err != nil {
				log.Fatalln(err)
			}
		}

	} else {
		fmt.Println("ELSE")
		data := conn.GetChatMsgs(Session)
		err := tpl.ExecuteTemplate(w, "chat.html", data)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func reader(conn *websocket.Conn) {
	for {
		mt, p, err := conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		log.Println(string(p))

		if err := conn.WriteMessage(mt, p); err != nil {
			panic(err)
		}
	}
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	reader(ws)
	_ = ws
}
