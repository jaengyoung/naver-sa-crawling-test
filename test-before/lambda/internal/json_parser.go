package internal

import (
	"fmt"
	"log"
	
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