package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ClientResponse struct {
	Bid string `json:"bid"`
}

func main() {
	// Contexto com timeout de 300ms para receber resultado do servidor
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Fazer requisição para o servidor
	cotacao, err := fetchCotacaoFromServer(ctx)
	if err != nil {
		log.Printf("Erro ao buscar cotação do servidor: %v", err)
		return
	}

	// Exibir cotação
	fmt.Printf("Cotação atual do Dólar: %s\n", cotacao.Bid)
	fmt.Printf("Cotação salva no arquivo cotacao.txt pelo servidor\n")
}

func fetchCotacaoFromServer(ctx context.Context) (*ClientResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cotacao ClientResponse
	err = json.NewDecoder(resp.Body).Decode(&cotacao)
	if err != nil {
		return nil, err
	}

	return &cotacao, nil
}
