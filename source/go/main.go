package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Body string `json:"body"`
}

type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

type Result struct {
	Language       string  `json:"language"`
	Goroutines     int     `json:"goroutines"`
	CountPerGoroutine int  `json:"count_per_goroutine"`
	DurationMs     float64 `json:"duration_ms"`
	Status         string  `json:"status"`
}

type ErrorResult struct {
	Error  string `json:"error"`
	Status string `json:"status"`
}

func countWorker(goroutineId int, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for i := 1; i <= 100; i++ {
		fmt.Printf("Goroutine %d: %d\n", goroutineId, i)
	}
}

func handleRequest(ctx context.Context, request Request) (Response, error) {
	startTime := time.Now()
	
	var wg sync.WaitGroup
	
	// Start 10 goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go countWorker(i, &wg)
	}
	
	// Wait for all goroutines to complete
	wg.Wait()
	
	endTime := time.Now()
	duration := endTime.Sub(startTime).Seconds() * 1000 // Convert to milliseconds
	
	result := Result{
		Language:          "Go",
		Goroutines:        10,
		CountPerGoroutine: 100,
		DurationMs:        duration,
		Status:            "completed",
	}
	
	resultJSON, err := json.Marshal(result)
	if err != nil {
		errorResult := ErrorResult{
			Error:  err.Error(),
			Status: "failed",
		}
		errorJSON, _ := json.Marshal(errorResult)
		return Response{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(errorJSON),
		}, nil
	}
	
	return Response{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(resultJSON),
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}