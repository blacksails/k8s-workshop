package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	handler := middleware.Logger(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		}),
	)

	fmt.Println("Listening on port :8080")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
