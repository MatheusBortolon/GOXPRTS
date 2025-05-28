package main

import (
	"context"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	req.Header.Set("Accept", "application/json")
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)

	tmp := template.New("CotacaoTemplate")
	tmp, _ = tmp.Parse("Dólar: {{.Bid}}")

	file, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = tmp.Execute(file, cotacao)
	if err != nil {
		panic(err)
	}
}
