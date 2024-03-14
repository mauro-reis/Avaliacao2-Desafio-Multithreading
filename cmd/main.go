package main

import (
	"avaliacao2-multithreading/configs"
	"avaliacao2-multithreading/internal/infra/webserver/handlers"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	println("Info: Configuration file loaded.")
	println("LIMITE_CONTAGEM:", config.LIMITE_CONTAGEM)

	cepHandler := handlers.NewCepHandler()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", config.TokenAuth))
	r.Use(middleware.WithValue("JwtExpiresIn", config.JWTExpiresIn))
	r.Use(middleware.WithValue("urlCep1", config.URL_CEP1))
	r.Use(middleware.WithValue("urlCep2", config.URL_CEP2))
	r.Use(middleware.WithValue("LIMITE_CONTAGEM", config.LIMITE_CONTAGEM))

	// Endpoint que ao acessado, este aciona um método que irá requisitar duas API's ao mesmo tempo.
	r.Route("/cep", func(r chi.Router) {
		r.Get("/buscacep/{cep}", cepHandler.BuscaCep)
	})
	http.ListenAndServe(":8080", r)
}
