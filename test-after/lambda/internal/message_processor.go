package internal

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	//	"time"
)

var (
	sqsClient = sqs.New(session.Must(session.NewSession(&aws.Config{Region: aws.String("ap-northeast-2")})))
	queueURL  = "https://sqs.ap-northeast-2.amazonaws.com/289023186990/skale-hourly-keyword-queue"
)

func ReceiveMessages() ([]*sqs.Message, error) {
	resp, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(2),
		VisibilityTimeout:   aws.Int64(5),
	})
	if err != nil {
		return nil, err
	}
	return resp.Messages, nil
}

func deleteMessage(receiptHandle *string) {
	_, err := sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: receiptHandle,
	})
	if err != nil {
		log.Printf("😡 [ERROR] Failed to delete message: %v", err)
	}
}

func ProcessMessage(message *sqs.Message, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()
	defer func() { <-sem }()

	var body SearchRequest
	err := json.Unmarshal([]byte(*message.Body), &body)
	if err != nil {
		log.Printf("😡 [ERROR] Failed to parse message: %v", err)
		return
	}

	desktopResults, err := ScrapeDesktopResults(body.Keyword)
	if err != nil {
		log.Printf("😡 [ERROR] Crawling error in PC for %s: %v", body.Keyword, err)
		return
	}

	deleteMessage(message.ReceiptHandle)

	if len(desktopResults) > 0 {
		uploadResult(desktopResults, body.Keyword)
		log.Printf("✅ Crawling completed: %s", body.Keyword)
	}
}
