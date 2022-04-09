package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/onur1/middleware"
	"github.com/onur1/stache/template"
	"github.com/rs/cors"
)

func layoutGet(m *middleware.Middleware) {
	t, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Printf("error: %v", err)
		m.Status(http.StatusInternalServerError)
	}

	ret := make([]interface{}, len(t))
	for i, v := range t {
		ret[i] = v.Serialize()
	}

	if err := m.SendJSON(ret, http.StatusOK); err != nil {
		log.Printf("error: send: %v", err)
	}
}

func main() {
	r := mux.NewRouter()
	ch := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
	})
	r.Methods(http.MethodGet).Path("/").Handler(execMiddleware(layoutGet))
	server := &http.Server{
		Addr:    ":7882",
		Handler: ch.Handler(r),
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

type handlerFunc = func(*middleware.Middleware)

func execMiddleware(f handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(middleware.NewMiddleware(w, r))
	}
}
