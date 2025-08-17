package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Listening on :8080")
	err := http.ListenAndServe(
		":8080",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("<html><body><h1>🐹 Hello world!</h1></body></html>"))
		}),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
