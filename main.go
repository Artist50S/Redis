package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func main() {
	// Initialize a connection pool and assign
	// it to the pool global variable.
	pool = &redis.Pool{
		MaxIdle:     10, // что такое Idle
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) { // что такое redis.Conn
			return redis.Dial("tcp", "0.0.0.0:4")
		},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/album", showAlbum)
	log.Println("listening on : http://127.0.0.1:4001")
	http.ListenAndServe(":4001", mux)
}
func showAlbum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(405), 405)
	}
	//Retrieve the id from the request URL query string. If there is
	// no id key in the query string then Get() will return an empty
	// string. We check for this, returning a 400 Bad Request response
	// if it's missing.
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
	}
	strconv.Atoi(id)
	bk, err := FindAlbum(id)
	if err == ErrNoAlbum {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Fprintf(w, "%s by %s $%.2f [%d likes]\n", bk.Title, bk.Artist, bk.Price, bk.Likes)
}
