package handlers

import (
	dto "avaliacao2-multithreading/internal/DTO"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
)

type CepHandler struct {
	CepBrasilAPI dto.CEPBrasilAPIOutputDTO
	CepViaCep    dto.CEPViaCepOutputDTO
}

func NewCepHandler() *CepHandler {
	return &CepHandler{}
}

var c1 = make(chan dto.CEPBrasilAPIOutputDTO)
var c2 = make(chan dto.CEPViaCepOutputDTO)
var message string

func (hand *CepHandler) GetCEPBrasilApi(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	urlCep := r.Context().Value("urlCep1").(string) + cep

	request, err := http.Get(urlCep)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("BuscaCEPBrasilApi. Erro na request: %v\n", err)
	}
	defer request.Body.Close()

	response, err := io.ReadAll(request.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "BuscaCEPBrasilApi. Erro na response: %v\n", err)
	}

	err = json.Unmarshal(response, &hand.CepBrasilAPI)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "BuscaCEPBrasilApi. Erro ao deserializar: %v\n", err)
	}

	c1 <- hand.CepBrasilAPI
	message = string(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (hand *CepHandler) GetCEPViaCep(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	urlCep := r.Context().Value("urlCep2").(string) + cep + "/json/"

	request, err := http.Get(urlCep)
	if err != nil {
		fmt.Printf("BuscaCEPViaCep. Erro na request: %v\n", err)
	}
	defer request.Body.Close()

	response, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "BuscaCEPViaCep. Erro na response: %v\n", err)
	}

	err = json.Unmarshal(response, &hand.CepViaCep)
	if err != nil {
		fmt.Fprintf(os.Stderr, "BuscaCEPViaCep. Erro ao deserializar: %v\n", err)
	}

	c2 <- hand.CepViaCep
	message = string(response)

	w.Header().Set("Content-Type", "application/json")
}

func (hand *CepHandler) BuscaCep(w http.ResponseWriter, r *http.Request) {
	go hand.GetCEPBrasilApi(w, r)
	go hand.GetCEPViaCep(w, r)

	select {
	case <-c1:
		w.Write([]byte("Service request: brasilapi.com.br:\n" + message))

	case <-c2:
		w.Write([]byte("Service request: viacep.com.br:\n" + message))

	case <-time.After(time.Second * 1):
		w.WriteHeader(http.StatusRequestTimeout)

		fmt.Fprintf(os.Stderr, "Resposta cancelada por prolongar mais de 1 segundo.\n")
	}
}
