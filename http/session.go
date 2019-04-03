package main

import (
	"flag"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var addr = flag.String("addr", ":8080", "server address")

func setSession(w http.ResponseWriter, r *http.Request) {
	// 得到一个session
	session, _ := store.Get(r, "session-name")
	// 设置session的一些值
	session.Values["foo"] = "bar"
	session.Values[42] = 43
	// 在返回之前保存它
	session.Save(r, w)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/session", setSession)

	http.ListenAndServe(*addr, mux)
}
