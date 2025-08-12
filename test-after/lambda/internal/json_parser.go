package internal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type JsonParser struct{}

func NewJsonParser() *JsonParser {
	return &JsonParser{}
}

func (jp *JsonParser) ParseUserInfo(jsonStr string) (string, string, error) {
	if !gjson.Valid(jsonStr) {
		return "", "", fmt.Errorf("invalid JSON format")
	}

	name := gjson.Get(jsonStr, "user.name")
	email := gjson.Get(jsonStr, "user.email")

	if !name.Exists() || !email.Exists() {
		return "", "", fmt.Errorf("required fields missing")
	}

	return name.String(), email.String(), nil
}

func (jp *JsonParser) GetNestedValue(jsonStr, path string) (string, error) {
	if !gjson.Valid(jsonStr) {
		return "", fmt.Errorf("invalid JSON format")
	}

	result := gjson.Get(jsonStr, path)
	if !result.Exists() {
		return "", fmt.Errorf("path not found: %s", path)
	}

	return result.String(), nil
}

func ExampleUsage() {
	parser := NewJsonParser()

	sampleJSON := `{
		"user": {
			"name": "John Doe",
			"email": "john@example.com",
			"profile": {
				"age": 30,
				"city": "Seoul"
			}
		},
		"timestamp": "2025-08-12T10:00:00Z"
	}`

	name, email, err := parser.ParseUserInfo(sampleJSON)
	if err != nil {
		log.Printf("Error parsing user info: %v", err)
		return
	}

	fmt.Printf("User: %s (%s)\n", name, email)

	city, err := parser.GetNestedValue(sampleJSON, "user.profile.city")
	if err != nil {
		log.Printf("Error getting city: %v", err)
		return
	}

	fmt.Printf("City: %s\n", city)
}

func (jp *JsonParser) CreateHTTPHandler() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/parse", func(c *gin.Context) {
		var requestBody struct {
			JSON string `json:"json"`
			Path string `json:"path"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if requestBody.Path == "" {
			name, email, err := jp.ParseUserInfo(requestBody.JSON)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"name": name, "email": email})
		} else {
			value, err := jp.GetNestedValue(requestBody.JSON, requestBody.Path)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"value": value})
		}
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	return r
}

// 주석 추가요~
