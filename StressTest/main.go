package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type StressTestConfig struct {
	URL         string
	Requests    int
	Concurrency int
}

type StressTestResult struct {
	TotalTime      time.Duration
	TotalRequests  int64
	StatusCodes    map[int]int64
	Errors         int64
	StartTime      time.Time
	EndTime        time.Time
	RequestsPerSec float64
	mu             sync.Mutex
}

func main() {
	config := parseFlags()

	if err := validateConfig(config); err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Iniciando teste de carga...\n")
	fmt.Printf("URL: %s\n", config.URL)
	fmt.Printf("Requests: %d\n", config.Requests)
	fmt.Printf("Concorrência: %d\n\n", config.Concurrency)

	result := runStressTest(config)
	printReport(result)
}

func parseFlags() StressTestConfig {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 0, "Número total de requests")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")

	flag.Parse()

	return StressTestConfig{
		URL:         *url,
		Requests:    *requests,
		Concurrency: *concurrency,
	}
}

func validateConfig(config StressTestConfig) error {
	if config.URL == "" {
		return fmt.Errorf("--url é obrigatório")
	}
	if config.Requests <= 0 {
		return fmt.Errorf("--requests deve ser maior que 0")
	}
	if config.Concurrency <= 0 {
		return fmt.Errorf("--concurrency deve ser maior que 0")
	}
	return nil
}

func runStressTest(config StressTestConfig) StressTestResult {
	result := StressTestResult{
		StatusCodes: make(map[int]int64),
		StartTime:   time.Now(),
	}

	var wg sync.WaitGroup
	requestChannel := make(chan struct{}, config.Concurrency)
	var totalRequests int64

	// Criar workers
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go worker(&wg, requestChannel, config.URL, &result, &totalRequests)
	}

	// Enviar requests após iniciar workers
	go func() {
		for i := 0; i < config.Requests; i++ {
			requestChannel <- struct{}{}
		}
		close(requestChannel)
	}()

	wg.Wait()

	result.EndTime = time.Now()
	result.TotalTime = result.EndTime.Sub(result.StartTime)
	result.TotalRequests = atomic.LoadInt64(&totalRequests)
	if result.TotalTime.Seconds() > 0 {
		result.RequestsPerSec = float64(result.TotalRequests) / result.TotalTime.Seconds()
	}

	return result
}

func worker(wg *sync.WaitGroup, requestChannel chan struct{}, url string, result *StressTestResult, totalRequests *int64) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for range requestChannel {
		resp, err := client.Get(url)

		atomic.AddInt64(totalRequests, 1)

		if err != nil {
			atomic.AddInt64(&result.Errors, 1)
			continue
		}

		statusCode := resp.StatusCode
		result.mu.Lock()
		result.StatusCodes[statusCode]++
		result.mu.Unlock()

		resp.Body.Close()
	}
}

func printReport(result StressTestResult) {
	separator := strings.Repeat("=", 60)
	dash := strings.Repeat("-", 60)

	fmt.Println("\n" + separator)
	fmt.Println("RELATÓRIO DE TESTE DE CARGA")
	fmt.Println(separator)
	fmt.Printf("Tempo total: %v\n", result.TotalTime)
	fmt.Printf("Total de requests: %d\n", result.TotalRequests)
	fmt.Printf("Requests por segundo: %.2f\n", result.RequestsPerSec)
	fmt.Printf("Erros: %d\n", result.Errors)
	fmt.Println("\nDistribuição de códigos HTTP:")
	fmt.Println(dash)

	// Exibir status 200 em destaque
	if count, ok := result.StatusCodes[200]; ok {
		fmt.Printf("  Status 200: %d\n", count)
	}

	// Exibir outros códigos de status
	for statusCode := 100; statusCode <= 599; statusCode++ {
		if count, ok := result.StatusCodes[statusCode]; ok && statusCode != 200 {
			fmt.Printf("  Status %d: %d\n", statusCode, count)
		}
	}

	if result.Errors > 0 {
		fmt.Printf("  Erros de conexão: %d\n", result.Errors)
	}

	fmt.Println(strings.Repeat("=", 60))
}
