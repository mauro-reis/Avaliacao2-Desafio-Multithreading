package handlers

import (
	dto "avaliacao2-multithreading/internal/DTO"
	"context"
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

func (hand *CepHandler) GetCEPBrasilApi(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	urlCep := r.Context().Value("urlCep1").(string) + cep
	// var contador = r.Context().Value("LIMITE_CONTAGEM").(int)
	// println("contador1:", contador)

	for i := 0; i < 100; i++ {
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
		w.Header().Add("servicerequest", "brasilapi.com.br")

		json.NewEncoder(w).Encode(hand.CepBrasilAPI)
		w.Write([]byte("servicerequest: brasilapi.com.br:\n" + string(response)))
	}
}

func (hand *CepHandler) GetCEPViaCep(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancelByTimeOut := context.WithTimeout(context.Background(), time.Second*1)
	defer cancelByTimeOut()

	urlCep := r.Context().Value("urlCep2").(string) + cep + "/json/"
	// contador := r.Context().Value("LIMITE_CONTAGEM").(string)
	// println("contador2:", contador)

	http.NewRequestWithContext(ctx, "GET", urlCep, nil)

	for i := 0; i < 100; i++ {
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

		// c2 <- hand.CepViaCep

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// w.Header().Add("servicerequest", "viacep.com.br")

		json.NewEncoder(w).Encode(hand.CepViaCep)
		w.Write([]byte("servicerequest: viacep.com.br:\n" + string(response)))
	}
}

func (hand *CepHandler) BuscaCep(w http.ResponseWriter, r *http.Request) {
	go hand.GetCEPBrasilApi(w, r)
	go hand.GetCEPViaCep(w, r)

	for {
		select {
		case <-c1:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprintf(os.Stderr, "Resposta originada de brasilapi.com.br\n")

		case <-c2:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprintf(os.Stderr, "Resposta originada de viacep.com.br\n")

		case <-time.After(time.Second * 1):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("Resposta cancelada por prolongar mais de 1 segundo.\n"))

			fmt.Fprintf(os.Stderr, "Resposta cancelada por prolongar mais de 1 segundo.\n")
		}
	}
}
