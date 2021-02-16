package ctrl

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

type Sess struct {
	s *sessions.Session
}

func CreateSess(r *http.Request) *Sess {
	var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	session, _ := store.Get(r, "session-name")
	return &Sess{session}
}

func (session *Sess) SaveSess(w http.ResponseWriter, r *http.Request) {
	err := session.s.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (session *Sess) SetVal(id string, data interface{}) {
	session.s.Values[id] = data
}

func (session *Sess) GetVal(id string) interface{} {
	return session.s.Values[id]
}
