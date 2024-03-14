package dto

// Classe DTO para inputs/outputs durante as comunicações.

type CEPInputDTO struct {
	Cep string `json:"cep"`
}

type CEPBrasilAPIOutputDTO struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type CEPViaCepOutputDTO struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type GetJWTInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetJWTOutPut struct {
	AccessToken string `json:"access_token"`
}
