package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/securecookie"
	"log"
	"net/http"
)

var (
	addr = flag.String("addr", ":8080", "server address")
)

func main() {
	mux := http.NewServeMux() // multiplexor
	mux.HandleFunc("/", index)
	mux.HandleFunc("/get", getCookie)
	mux.HandleFunc("/delete", deleteCookie)
	mux.HandleFunc("/set", setCookie)
	mux.HandleFunc("/sets", setSecureCookie)

	http.ListenAndServe(*addr, mux)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<a href="#" onclick="alert(document.cookie)">Click me</a>`))
}

func getCookie(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("test")

	if err != nil {
		w.Write([]byte("get cookie failed: " + err.Error()))
	} else {
		data, _ := json.MarshalIndent(c, "", "\t")
		w.Write([]byte("cookie is :\n" + string(data)))
	}
}

func deleteCookie(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:   "test",
		MaxAge: -1,
	}

	http.SetCookie(w, &c)

	w.Write([]byte("cookie has deleted\n"))
}

func setCookie(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:     "test",
		Value:    "true",
		HttpOnly: true,
		//Secure:   true,
		MaxAge: 300,
	}

	http.SetCookie(w, &c)

	w.Write([]byte("cooke has created\n"))
}

func setSecureCookie(w http.ResponseWriter, r *http.Request) {
	var (
		hashKey  = securecookie.GenerateRandomKey(16)
		blockKey = securecookie.GenerateRandomKey(16)
		s        = securecookie.New(hashKey, blockKey)
	)

	value := map[string]string{
		"foo": "bar",
	}

	encoded, err := s.Encode("cookie-name", value)
	if err != nil {
		w.Write([]byte("set secure cookie failed: " + err.Error()))
	} else {
		cookie := http.Cookie{
			Name:     "cookie-name",
			Value:    encoded,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   300,
		}

		http.SetCookie(w, &cookie)
		w.Write([]byte("set secure cookie succeed"))
	}

	// decode
	decode := make(map[string]string)
	s.Decode("cookie-name", encoded, &decode)
	log.Println(decode)
}
