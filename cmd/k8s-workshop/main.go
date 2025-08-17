package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/blacksails/k8s-workshop/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	publicFS, err := fs.Sub(web.FS(), "public")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	middlewares := []func(http.Handler) http.Handler{
		middleware.Logger,
	}

	user := os.Getenv("K8S_WORKSHOP_USER")
	pass := os.Getenv("K8S_WORKSHOP_PASS")
	if user != "" && pass != "" {
		middlewares = append(middlewares, middleware.BasicAuth("Kubernetes Workshop", map[string]string{
			user: pass,
		}))
	}

	mux := chi.NewMux()
	mux.Use(middlewares...)
	mux.Handle("/*", http.FileServerFS(publicFS))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
