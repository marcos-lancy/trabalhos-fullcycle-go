package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

type TestConfig struct {
	URL         string
	Requests    int
	Concurrency int
}

type TestResult struct {
	StatusCode int
	Duration   time.Duration
	Error      error
}

type TestReport struct {
	TotalRequests     int                    `json:"total_requests"`
	TotalTime         time.Duration          `json:"total_time"`
	SuccessfulRequests int                   `json:"successful_requests"`
	StatusDistribution map[int]int           `json:"status_distribution"`
	AverageResponseTime time.Duration        `json:"average_response_time"`
	MinResponseTime    time.Duration         `json:"min_response_time"`
	MaxResponseTime    time.Duration         `json:"max_response_time"`
}

var (
	url         string
	requests    int
	concurrency int
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "stress-test",
		Short: "Ferramenta de teste de carga para serviços web",
		Long:  "Uma ferramenta CLI em Go para realizar testes de carga em serviços web",
		Run:   runStressTest,
	}

	rootCmd.Flags().StringVar(&url, "url", "", "URL do serviço a ser testado (obrigatório)")
	rootCmd.Flags().IntVar(&requests, "requests", 0, "Número total de requests (obrigatório)")
	rootCmd.Flags().IntVar(&concurrency, "concurrency", 0, "Número de chamadas simultâneas (obrigatório)")

	rootCmd.MarkFlagRequired("url")
	rootCmd.MarkFlagRequired("requests")
	rootCmd.MarkFlagRequired("concurrency")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao executar o comando: %v\n", err)
		os.Exit(1)
	}
}

func runStressTest(cmd *cobra.Command, args []string) {
	config := TestConfig{
		URL:         url,
		Requests:    requests,
		Concurrency: concurrency,
	}

	fmt.Printf("Iniciando teste de carga...\n")
	fmt.Printf("URL: %s\n", config.URL)
	fmt.Printf("Requests: %d\n", config.Requests)
	fmt.Printf("Concorrência: %d\n", config.Concurrency)
	fmt.Println("----------------------------------------")

	report := executeLoadTest(config)
	printReport(report)
}

func executeLoadTest(config TestConfig) TestReport {
	startTime := time.Now()
	
	results := make(chan TestResult, config.Requests)
	semaphore := make(chan struct{}, config.Concurrency)
	var wg sync.WaitGroup
	
	for i := 0; i < config.Requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			result := makeRequest(config.URL)
			results <- result
		}()
	}
	
	go func() {
		wg.Wait()
		close(results)
	}()
	
	var testResults []TestResult
	for result := range results {
		testResults = append(testResults, result)
	}
	
	totalTime := time.Since(startTime)
	report := processResults(testResults, totalTime)
	
	return report
}

func makeRequest(url string) TestResult {
	startTime := time.Now()
	
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return TestResult{
			StatusCode: 0,
			Duration:   time.Since(startTime),
			Error:      err,
		}
	}
	
	req.Header.Set("User-Agent", "Stress-Test-Tool/1.0")
	
	resp, err := client.Do(req)
	if err != nil {
		return TestResult{
			StatusCode: 0,
			Duration:   time.Since(startTime),
			Error:      err,
		}
	}
	defer resp.Body.Close()
	
	_, _ = io.Copy(io.Discard, resp.Body)
	
	return TestResult{
		StatusCode: resp.StatusCode,
		Duration:   time.Since(startTime),
		Error:      nil,
	}
}

func processResults(results []TestResult, totalTime time.Duration) TestReport {
	report := TestReport{
		TotalRequests:      len(results),
		TotalTime:          totalTime,
		SuccessfulRequests: 0,
		StatusDistribution: make(map[int]int),
		MinResponseTime:    time.Duration(^uint64(0) >> 1), // Max duration
		MaxResponseTime:    0,
	}
	
	var totalResponseTime time.Duration
	
	for _, result := range results {
		if result.StatusCode == 200 {
			report.SuccessfulRequests++
		}
		
		report.StatusDistribution[result.StatusCode]++
		
		if result.Duration > 0 {
			totalResponseTime += result.Duration
			
			if result.Duration < report.MinResponseTime {
				report.MinResponseTime = result.Duration
			}
			
			if result.Duration > report.MaxResponseTime {
				report.MaxResponseTime = result.Duration
			}
		}
	}
	
	if len(results) > 0 {
		report.AverageResponseTime = totalResponseTime / time.Duration(len(results))
	}
	
	if report.MinResponseTime == time.Duration(^uint64(0)>>1) {
		report.MinResponseTime = 0
	}
	
	return report
}

func printReport(report TestReport) {
	fmt.Println("\n========================================")
	fmt.Println("RELATÓRIO DO TESTE DE CARGA")
	fmt.Println("========================================")
	
	fmt.Printf("Tempo total de execução: %v\n", report.TotalTime)
	fmt.Printf("Total de requests realizados: %d\n", report.TotalRequests)
	fmt.Printf("Requests com status 200: %d\n", report.SuccessfulRequests)
	
	fmt.Println("\nDistribuição de códigos de status:")
	for status, count := range report.StatusDistribution {
		if status == 0 {
			fmt.Printf("  Erros de conexão: %d\n", count)
		} else {
			fmt.Printf("  Status %d: %d\n", status, count)
		}
	}
	
	fmt.Println("\nEstatísticas de tempo de resposta:")
	fmt.Printf("  Tempo médio: %v\n", report.AverageResponseTime)
	fmt.Printf("  Tempo mínimo: %v\n", report.MinResponseTime)
	fmt.Printf("  Tempo máximo: %v\n", report.MaxResponseTime)
	
	if report.TotalTime.Seconds() > 0 {
		rps := float64(report.TotalRequests) / report.TotalTime.Seconds()
		fmt.Printf("  Requests por segundo: %.2f\n", rps)
	}
	
	fmt.Println("========================================")
}
