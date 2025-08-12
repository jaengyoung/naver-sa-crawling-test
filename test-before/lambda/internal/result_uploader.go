package internal

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

var (
	s3Client = s3.New(session.Must(session.NewSession(&aws.Config{Region: aws.String("ap-northeast-2")})))
	bucket   = "skale-crawling-manager"
)

func uploadResult(result []SearchResult, keyword string) {
	if len(result) == 0 {
		log.Println("🚫 No results to upload. Skipping S3 upload.")
		return
	}

	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic(err)
	}

	now := time.Now().In(loc)

	basicDate := now.Format("20060102")
	hour := now.Hour()
	uuid := uuid.New().String()[:18]

	key := fmt.Sprintf("data/basic_date=%s/hh=%d/%s.csv.gz", basicDate, hour, uuid)

	buffer := new(bytes.Buffer)
	gzWriter := gzip.NewWriter(buffer)
	csvWriter := csv.NewWriter(gzWriter)

	// CSV 헤더 작성
	header := []string{"query", "device", "rank", "site_name", "display_url", "title", "description"}
	if err := csvWriter.Write(header); err != nil {
		log.Printf("❌ Failed to write CSV header: %v", err)
		return
	}

	// 데이터 작성
	for _, item := range result {
		// Device values are already optimized in scrapers (PC/Mobile)
		record := []string{
			item.Query,
			item.Device,
			strconv.Itoa(item.Rank),
			item.SiteName,
			item.DisplayURL,
			item.Title,
			item.Description,
		}
		if err := csvWriter.Write(record); err != nil {
			log.Printf("❌ Failed to write CSV record: %v", err)
			return
		}
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		log.Printf("❌ CSV writer error: %v", err)
		return
	}

	if err := gzWriter.Close(); err != nil {
		log.Printf("❌ Failed to close gzip writer: %v", err)
		return
	}

	reader := bytes.NewReader(buffer.Bytes())
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          reader,
		ContentLength: aws.Int64(int64(reader.Len())),
		ContentType:   aws.String("application/gzip"),
	})

	if err != nil {
		log.Printf("❌ Failed to upload CSV.GZ to S3: %v", err)
		return
	}

	log.Printf("✅ Successfully uploaded %d records to S3: %s", len(result), key)
}
