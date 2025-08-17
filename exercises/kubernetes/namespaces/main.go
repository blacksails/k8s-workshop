package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("echosvc expects exactly one argument")
		os.Exit(1)
	}
	msg := os.Args[1]

	fmt.Println("Listening on port :8080")
	fmt.Printf("Configured to reply the following:\n%s\n", msg)

	err := http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(msg))
	}))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
