package handlers

import (
	dto "avaliacao2-multithreading/internal/DTO"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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
		fmt.Fprint(os.Stderr, "BuscaCEPBrasilApi. Erro na request: %v\n", err)
	}
	defer request.Body.Close()

	response, err := io.ReadAll(request.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "BuscaCEPBrasilApi. Erro na response: %v\n", err)
	}

	fmt.Println("servicerequest brasilapi.com.br:", string(response))

	err = json.Unmarshal(response, &hand.CepBrasilAPI)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "BuscaCEPBrasilApi. Erro ao deserializar: %v\n", err)
	}

	c1 <- hand.CepBrasilAPI

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(hand.CepBrasilAPI)
	w.Write([]byte("servicerequest: brasilapi.com.br:\n" + string(response)))
	//}
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
		fmt.Fprint(os.Stderr, "BuscaCEPViaCep. Erro na request: %v\n", err)
	}
	defer request.Body.Close()

	response, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "BuscaCEPViaCep. Erro na response: %v\n", err)
	}

	fmt.Println("servicerequest viacep.com.br:", string(response))

	err = json.Unmarshal(response, &hand.CepViaCep)
	if err != nil {
		fmt.Fprintf(os.Stderr, "BuscaCEPViaCep. Erro ao deserializar: %v\n", err)
	}

	c2 <- hand.CepViaCep

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(hand.CepViaCep)
	w.Write([]byte("servicerequest: viacep.com.br:\n" + string(response)))
}

func (hand *CepHandler) BuscaCep(w http.ResponseWriter, r *http.Request) {
	contador := r.Context().Value("LIMITE_CONTAGEM").(string)
	intcontador, _ := strconv.Atoi(contador)

	for i := 0; i < int(intcontador); i++ {
		go hand.GetCEPBrasilApi(w, r)
		go hand.GetCEPViaCep(w, r)

		select {
		case <-c1:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

		case <-c2:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

		case <-time.After(time.Second * 1):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("Resposta cancelada por prolongar mais de 1 segundo.\n"))

			fmt.Fprintf(os.Stderr, string(i)+"Resposta cancelada por prolongar mais de 1 segundo.\n")
		}
	}
}
