## 핵심 비즈니스 로직 (Core Business Logic)

### 전체 시스템 아키텍처

이 네이버 검색 크롤러는 AWS Lambda 환경에서 실행되며, SQS 큐로부터 키워드를 받아 데스크톱/모바일 검색 결과를 수집하고 S3에 저장하는 시스템입니다.

```
[SQS Queue] → [Lambda Function] → [Naver Search] → [S3 Storage]
     ↓              ↓                    ↓              ↓
  키워드 메시지    크롤링 실행        검색결과 추출    CSV 파일 저장
```

### 주요 구성 요소 및 처리 흐름

#### 1. 메시지 처리 (`message_processor.go`)

SQS에서 메시지를 수신하고 병렬 처리합니다:

```go
package main

import (
    "log"
    "sync"
    "lambda/internal"
)

func main() {
    // SQS에서 메시지 수신
    messages, err := internal.ReceiveMessages()
    if err != nil {
        log.Fatal("메시지 수신 실패:", err)
    }
    
    // 동시 처리를 위한 고루틴 풀
    var wg sync.WaitGroup
    sem := make(chan struct{}, 10) // 최대 10개 동시 처리
    
    for _, message := range messages {
        wg.Add(1)
        sem <- struct{}{} // 세마포어로 동시성 제어
        
        go func(msg *sqs.Message) {
            defer wg.Done()
            defer func() { <-sem }()
            
            // 각 키워드에 대해 크롤링 실행
            internal.ProcessMessage(msg, &wg, sem)
        }(message)
    }
    
    wg.Wait()
    log.Println("모든 크롤링 작업 완료")
}
```

#### 2. 크롤링 로직 (`desktop_scraper.go`, `mobile_scraper.go`)

각 디바이스별로 최적화된 크롤링을 수행합니다:

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "lambda/internal"
)

func 크롤링_예제() {
    keyword := "스마트폰"
    
    // 데스크톱 크롤링
    fmt.Println("🖥️ 데스크톱 크롤링 시작...")
    desktopResults, err := internal.ScrapeDesktopResultsWithDebug(keyword)
    if err != nil {
        log.Printf("데스크톱 크롤링 실패: %v", err)
    }
    
    // 모바일 크롤링  
    fmt.Println("📱 모바일 크롤링 시작...")
    mobileResults, err := internal.ScrapeMobileResultsWithDebug(keyword)
    if err != nil {
        log.Printf("모바일 크롤링 실패: %v", err)
    }
    
    // 결과 통합 및 출력
    allResults := append(desktopResults, mobileResults...)
    
    fmt.Printf("📊 크롤링 완료: 총 %d개 결과\n", len(allResults))
    fmt.Printf("  - 데스크톱: %d개\n", len(desktopResults))
    fmt.Printf("  - 모바일: %d개\n", len(mobileResults))
    
    // JSON 형태로 출력
    jsonData, _ := json.MarshalIndent(allResults, "", "  ")
    fmt.Println(string(jsonData))
}
```

#### 3. 데이터 업로드 (`result_uploader.go`)

크롤링 결과를 CSV 형태로 압축하여 S3에 저장합니다:

```go
// 실제 업로드 프로세스 예제
func 업로드_예제() {
    results := []internal.SearchResult{
        {
            Query:       "스마트폰",
            Device:      "Mobile", 
            Rank:        1,
            SiteName:    "디베이",
            DisplayURL:  "dbay.io",
            Title:       "견적비교 플랫폼 디베이",
            Description: "700명 딜러들의 혜택경쟁으로 가격은 DOWN 혜택은 UP!",
        },
        // ... 더 많은 결과들
    }
    
    // S3 업로드 (내부적으로 CSV 변환 및 GZIP 압축 수행)
    // 파일 경로: s3://bucket/data/basic_date=20250811/hh=14/uuid.csv.gz
    internal.UploadResult(results, "스마트폰")
}
```

### 탐지 회피 메커니즘

#### 1. 랜덤 헤더 생성 (`http_client.go`)

```go
func 헤더_예제() {
    // 매 요청마다 다른 헤더 조합 생성
    headers := internal.GenerateRandomHeaders()
    
    fmt.Println("생성된 헤더:")
    for key, value := range headers {
        fmt.Printf("  %s: %s\n", key, value)
    }
    
    // 출력 예시:
    // User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36...
    // Referer: https://www.naver.com/
    // Accept-Language: ko-KR,ko;q=0.9,en-US;q=0.8,en;q=0.7
    // Origin: https://search.naver.com (50% 확률로 추가)
    // Cookie: NID=abc123; NNB=xyz456 (30% 확률로 추가)
}
```

#### 2. CSS 셀렉터 최적화

각 플랫폼별로 특화된 셀렉터를 사용합니다:

```go
// 데스크톱용 셀렉터
const (
    desktopResultSelector = "div.nad_area ul.lst_type > li"
    desktopTitleSelector  = "a.lnk_head span.lnk_tit"  
    desktopSiteSelector   = "a.site"
    desktopURLSelector    = "span.lnk_url_area > a.lnk_url"
    desktopDescSelector   = "a.link_desc"
)

// 모바일용 셀렉터 
const (
    mobileResultSelector = "div.api_subject_bx ul.lst_total > li"
    mobileTitleSelector  = "div.tit_area span.tit"
    mobileSiteSelector   = "span.site" 
    mobileURLSelector    = "span.url"
    mobileDescSelector   = "a.desc"
)
```

### 실제 사용 시나리오

#### 시나리오 1: 단일 키워드 크롤링

```go
func 단일_키워드_크롤링() {
    keyword := "갤럭시S25"
    
    // 모바일 우선 크롤링 (일반적으로 더 안정적)
    results, err := internal.ScrapeMobileResults(keyword)
    if err != nil {
        log.Printf("모바일 크롤링 실패: %v", err)
        return
    }
    
    fmt.Printf("'%s' 검색결과 %d개 수집 완료\n", keyword, len(results))
    
    // 각 결과 출력
    for i, result := range results {
        fmt.Printf("%d. [%s] %s - %s\n", 
            result.Rank, result.SiteName, result.Title, result.DisplayURL)
    }
}
```

#### 시나리오 2: 배치 처리

```go
func 배치_처리_예제() {
    keywords := []string{"갤럭시S25", "아이폰16", "픽셀9"}
    
    var allResults []internal.SearchResult
    
    for _, keyword := range keywords {
        fmt.Printf("🔍 '%s' 크롤링 중...\n", keyword)
        
        // 데스크톱과 모바일 동시 처리
        var wg sync.WaitGroup
        var desktopResults, mobileResults []internal.SearchResult
        
        wg.Add(2)
        
        // 데스크톱 크롤링 고루틴
        go func() {
            defer wg.Done()
            results, err := internal.ScrapeDesktopResults(keyword)
            if err == nil {
                desktopResults = results
            }
        }()
        
        // 모바일 크롤링 고루틴  
        go func() {
            defer wg.Done()
            results, err := internal.ScrapeMobileResults(keyword)
            if err == nil {
                mobileResults = results
            }
        }()
        
        wg.Wait()
        
        // 결과 합치기
        allResults = append(allResults, desktopResults...)
        allResults = append(allResults, mobileResults...)
        
        fmt.Printf("✅ '%s' 완료: 데스크톱 %d개, 모바일 %d개\n", 
            keyword, len(desktopResults), len(mobileResults))
    }
    
    fmt.Printf("🎉 전체 배치 처리 완료: 총 %d개 결과\n", len(allResults))
}
```

### 에러 처리 및 복구

```go
func 안정적인_크롤링() {
    keyword := "스마트폰"
    maxRetries := 3
    
    for i := 0; i < maxRetries; i++ {
        results, err := internal.ScrapeMobileResults(keyword)
        if err == nil {
            fmt.Printf("✅ 크롤링 성공: %d개 결과\n", len(results))
            return
        }
        
        fmt.Printf("⚠️ 시도 %d/%d 실패: %v\n", i+1, maxRetries, err)
        
        if i < maxRetries-1 {
            // 재시도 전 대기 (백오프 패턴)
            time.Sleep(time.Duration(i+1) * time.Second)
        }
    }
    
    fmt.Println("❌ 모든 재시도 실패")
}
```

### 테스트 실행

```bash
# 개발 환경에서 테스트
cd source/lambda

# 데스크톱 크롤링 테스트
go run test_desktop.go -keyword="갤럭시S25"

# 모바일 크롤링 테스트  
go run test_mobile.go -keyword="갤럭시S25"

# 실제 Lambda 환경 배포 후
# SQS에 메시지 전송하면 자동으로 크롤링 실행 및 S3 업로드 수행
```

이 시스템은 키워드 기반 검색 결과를 안정적으로 수집하여 비즈니스 인텔리전스나 마케팅 분석에 활용할 수 있는 구조화된 데이터를 제공합니다.

## License

This project is for educational and research purposes. Please ensure compliance with Naver's robots.txt and terms of service when using this crawler.