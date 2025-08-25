package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Estrutura para resposta da BrasilAPI
type BrasilAPIResponse struct {
	CEP          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
	Location     struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"location"`
}

// Estrutura para resposta da ViaCEP
type ViaCEPResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	SIAFI       string `json:"siafi"`
}

// Estrutura para resultado unificado
type AddressResult struct {
	API          string
	CEP          string
	Street       string
	Neighborhood string
	City         string
	State        string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Uso: go run main.go <CEP>")
		fmt.Println("Exemplo: go run main.go 01153000")
		os.Exit(1)
	}

	cep := os.Args[1]

	// Contexto com timeout de 1 segundo
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Canal para receber o resultado mais rápido
	resultChan := make(chan AddressResult, 1)
	errorChan := make(chan error, 1)

	// Goroutine para BrasilAPI
	go func() {
		result, err := fetchBrasilAPI(ctx, cep)
		if err != nil {
			errorChan <- fmt.Errorf("BrasilAPI: %v", err)
			return
		}
		select {
		case resultChan <- result:
		default:
		}
	}()

	// Goroutine para ViaCEP
	go func() {
		result, err := fetchViaCEP(ctx, cep)
		if err != nil {
			errorChan <- fmt.Errorf("ViaCEP: %v", err)
			return
		}
		select {
		case resultChan <- result:
		default:
		}
	}()

	// Aguarda o primeiro resultado ou timeout
	select {
	case result := <-resultChan:
		displayResult(result)
	case <-errorChan:
		// Se uma API falhou, aguarda a outra
		select {
		case result := <-resultChan:
			displayResult(result)
		case <-ctx.Done():
			fmt.Println("Erro: Timeout - nenhuma API respondeu em 1 segundo")
		}
	case <-ctx.Done():
		fmt.Println("Erro: Timeout - nenhuma API respondeu em 1 segundo")
	}
}

func fetchBrasilAPI(ctx context.Context, cep string) (AddressResult, error) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return AddressResult{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AddressResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AddressResult{}, err
	}

	if resp.StatusCode != 200 {
		return AddressResult{}, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var brasilResp BrasilAPIResponse
	if err := json.Unmarshal(body, &brasilResp); err != nil {
		return AddressResult{}, err
	}

	return AddressResult{
		API:          "BrasilAPI",
		CEP:          brasilResp.CEP,
		Street:       brasilResp.Street,
		Neighborhood: brasilResp.Neighborhood,
		City:         brasilResp.City,
		State:        brasilResp.State,
	}, nil
}

func fetchViaCEP(ctx context.Context, cep string) (AddressResult, error) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return AddressResult{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AddressResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AddressResult{}, err
	}

	if resp.StatusCode != 200 {
		return AddressResult{}, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var viaResp ViaCEPResponse
	if err := json.Unmarshal(body, &viaResp); err != nil {
		return AddressResult{}, err
	}

	return AddressResult{
		API:          "ViaCEP",
		CEP:          viaResp.CEP,
		Street:       viaResp.Logradouro,
		Neighborhood: viaResp.Bairro,
		City:         viaResp.Localidade,
		State:        viaResp.UF,
	}, nil
}

func displayResult(result AddressResult) {
	fmt.Println("=== Resultado da Consulta de CEP ===")
	fmt.Printf("API: %s\n", result.API)
	fmt.Printf("CEP: %s\n", result.CEP)
	fmt.Printf("Logradouro: %s\n", result.Street)
	fmt.Printf("Bairro: %s\n", result.Neighborhood)
	fmt.Printf("Cidade: %s\n", result.City)
	fmt.Printf("Estado: %s\n", result.State)
}
