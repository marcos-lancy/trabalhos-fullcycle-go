package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type CotacaoResponse struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

type Cotacao struct {
	ID   int    `json:"id"`
	Bid  string `json:"bid"`
	Date string `json:"date"`
}

type CotacaoDB struct {
	Cotacoes []Cotacao `json:"cotacoes"`
	mu       sync.Mutex
}

var db = &CotacaoDB{
	Cotacoes: []Cotacao{},
}

func main() {
	// Carregar dados existentes se o arquivo existir
	loadDatabase()

	// Configurar rota
	http.HandleFunc("/cotacao", handleCotacao)

	// Iniciar servidor na porta 8080
	fmt.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadDatabase() {
	data, err := os.ReadFile("cotacao.txt")
	if err == nil {
		// Tentar decodificar como JSON, se falhar, ignorar
		json.Unmarshal(data, &db.Cotacoes)
	}
}

func saveDatabase() {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	// Salvar apenas a cotação mais recente no formato simples
	if len(db.Cotacoes) > 0 {
		ultimaCotacao := db.Cotacoes[len(db.Cotacoes)-1]
		content := fmt.Sprintf("Dólar: %s", ultimaCotacao.Bid)
		err := os.WriteFile("cotacao.txt", []byte(content), 0644)
		if err != nil {
			log.Printf("Erro ao salvar arquivo: %v", err)
		}
	}
}

func handleCotacao(w http.ResponseWriter, r *http.Request) {
	// Contexto com timeout de 200ms para a API
	ctxAPI, cancelAPI := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancelAPI()

	// Fazer requisição para a API de cotação
	cotacao, err := fetchCotacao(ctxAPI)
	if err != nil {
		log.Printf("Erro ao buscar cotação: %v", err)
		http.Error(w, "Erro ao buscar cotação", http.StatusInternalServerError)
		return
	}

	// Contexto com timeout de 10ms para o banco de dados
	ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelDB()

	// Salvar no banco de dados
	err = saveCotacao(ctxDB, cotacao.USDBRL.Bid)
	if err != nil {
		log.Printf("Erro ao salvar no banco: %v", err)
		// Continua mesmo com erro no banco, pois o cliente precisa receber a cotação
	}

	// Retornar apenas o valor do bid
	response := map[string]string{
		"bid": cotacao.USDBRL.Bid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func fetchCotacao(ctx context.Context) (*CotacaoResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cotacao CotacaoResponse
	err = json.NewDecoder(resp.Body).Decode(&cotacao)
	if err != nil {
		return nil, err
	}

	return &cotacao, nil
}

func saveCotacao(ctx context.Context, bid string) error {
	// Verificar se o contexto foi cancelado
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	// Criar nova cotação
	novaCotacao := Cotacao{
		ID:   len(db.Cotacoes) + 1,
		Bid:  bid,
		Date: time.Now().Format("2006-01-02 15:04:05"),
	}

	// Adicionar à lista
	db.Cotacoes = append(db.Cotacoes, novaCotacao)

	// Salvar no arquivo
	go saveDatabase()

	return nil
}
