package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const GEMINI_URL = "https://generativelanguage.googleapis.com/v1/models/gemini-1.5-flash:generateContent"

type TextPart struct {
	Text string `json:"text"`
}

type Content struct {
	Parts []TextPart `json:"parts"`
}

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Candidate struct {
	Content Content `json:"content"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

func gpt(content string) (string, error) {
	godotenv.Load()

	apiKey := os.Getenv("GEMINI_API_KEY")

	requestBody := GeminiRequest{
		Contents: []Content{
			{Parts: []TextPart{{Text: content}}},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", GEMINI_URL+"?key="+apiKey, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", errors.New(string(body))
	}

	var geminiResp GeminiResponse
	err = json.Unmarshal(body, &geminiResp)
	if err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}
	return "", errors.New("something went wrong")
}
