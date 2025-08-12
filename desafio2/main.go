package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Endereco struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
}

type Resultado struct {
	API      string
	Endereco Endereco
	Err      error
}

func buscaBrasilAPI(ctx context.Context, cep string, ch chan<- Resultado) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/01153000%s", cep)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- Resultado{API: "BrasilAPI", Err: err}
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var e Endereco
	if err := json.Unmarshal(body, &e); err != nil {
		ch <- Resultado{API: "BrasilAPI", Err: err}
		return
	}
	ch <- Resultado{API: "BrasilAPI", Endereco: e}
}

func buscaViaCep(ctx context.Context, cep string, ch chan<- Resultado) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- Resultado{API: "ViaCEP", Err: err}
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var e Endereco
	if err := json.Unmarshal(body, &e); err != nil {
		ch <- Resultado{API: "ViaCEP", Err: err}
		return
	}
	ch <- Resultado{API: "ViaCEP", Endereco: e}
}

func main() {
	cep := "89035300"
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	resultCh := make(chan Resultado)

	go buscaBrasilAPI(ctx, cep, resultCh)
	go buscaViaCep(ctx, cep, resultCh)

	select {
	case res := <-resultCh:
		if res.Err != nil {
			fmt.Printf("Erro na API %s: %v\n", res.API, res.Err)
			return
		}
		fmt.Printf("Resposta da API: %s\n", res.API)
		fmt.Printf("Endereço: %+v\n", res.Endereco)
	case <-ctx.Done():
		fmt.Println("Erro: Timeout de 1 segundo atingido.")
	}
}
