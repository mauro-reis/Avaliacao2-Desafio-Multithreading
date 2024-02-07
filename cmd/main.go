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

	//TO DO: Verificar a possibilidade de criar um banco de dados MySql via docker e integrá-lo.

	cepHandler := handlers.NewCepHandler()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", config.TokenAuth))
	r.Use(middleware.WithValue("JwtExpiresIn", config.JWTExpiresIn))
	r.Use(middleware.WithValue("urlCep1", config.URL_CEP1))
	r.Use(middleware.WithValue("urlCep2", config.URL_CEP2))
	r.Use(middleware.WithValue("LIMITE_CONTAGEM", config.LIMITE_CONTAGEM))

	r.Route("/cep", func(r chi.Router) {
		// r.Get("/BrasilApi/{cep}", cepHandler.GetCEPBrasilApi)
		r.Get("/{cep}", cepHandler.GetCEPBrasilApi)
		r.Get("/ViaCep/{cep}", cepHandler.GetCEPViaCep)
		r.Get("/buscacep/{cep}", cepHandler.BuscaCep)
	})
	// http.HandleFunc("/CepV1", BuscaCEPBrasilApi)
	// http.HandleFunc("/CepV2", BuscaCEPViaCep)
	http.ListenAndServe(":8080", r)
}
