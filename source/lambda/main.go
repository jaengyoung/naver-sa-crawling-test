package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"lambda/internal"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	lambda.Start(handler)
}

func processRound(roundNum int) (int, error) {
	messages, err := internal.ReceiveMessages()
	if err != nil {
		log.Printf("❌ Failed to receive messages: %v\n", err)
		return 0, err
	}

	if len(messages) == 0 {
		return 0, nil
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	log.Printf("🔄 Round %d: Processing %d keywords with 10 goroutines", roundNum, len(messages))

	for _, msg := range messages {
		localMsg := msg
		wg.Add(1)
		sem <- struct{}{}
		go internal.ProcessMessage(localMsg, &wg, sem)
	}

	wg.Wait()
	log.Printf("✅ Round %d completed: %d keywords processed", roundNum, len(messages))
	return len(messages), nil
}

func handler(_ context.Context) (string, error) {
	// Registry Cache 테스트를 위한 수정
	const totalRounds = 5
	totalProcessed := 0

	log.Printf("🚀 Starting lambda execution: %d rounds, 10 keywords per round", totalRounds)

	for round := 1; round <= totalRounds; round++ {
		processed, err := processRound(round)
		if err != nil {
			return "", fmt.Errorf("error in round %d: %v", round, err)
		}

		if processed == 0 {
			log.Printf("🏁 Round %d: No messages to process, stopping early", round)
			break
		}

		totalProcessed += processed
		log.Printf("📊 Round %d summary: %d keywords processed (Total: %d)", round, processed, totalProcessed)
	}

	return fmt.Sprintf("Lambda completed. Total rounds: %d, Total keywords processed: %d", totalRounds, totalProcessed), nil
}
