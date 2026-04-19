package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Question holds the fields returned by the LeetCode GraphQL API.
type Question struct {
	QuestionFrontendID string        `json:"questionFrontendId"`
	Title              string        `json:"title"`
	TitleSlug          string        `json:"titleSlug"`
	Content            string        `json:"content"`
	Difficulty         string        `json:"difficulty"`
	ExampleTestcases   string        `json:"exampleTestcases"`
	CodeSnippets       []CodeSnippet `json:"codeSnippets"`
	MetaData           string        `json:"metaData"`
	SampleTestCase     string        `json:"sampleTestCase"`
	IsPaidOnly         bool          `json:"isPaidOnly"`
}

// CodeSnippet is a language-specific starter code template.
type CodeSnippet struct {
	Lang     string `json:"lang"`
	LangSlug string `json:"langSlug"`
	Code     string `json:"code"`
}

const gqlQuery = `query questionData($titleSlug: String!) {
  question(titleSlug: $titleSlug) {
    questionFrontendId
    title
    titleSlug
    content
    difficulty
    exampleTestcases
    codeSnippets { lang langSlug code }
    metaData
    sampleTestCase
    isPaidOnly
  }
}`

var httpClient = &http.Client{Timeout: 15 * time.Second}

// fetchQuestion calls the LeetCode GraphQL API and returns the question data.
// Returns (nil, nil) if the question does not exist.
func fetchQuestion(slug string) (*Question, error) {
	payload, err := json.Marshal(map[string]any{
		"query":         gqlQuery,
		"variables":     map[string]string{"titleSlug": slug},
		"operationName": "questionData",
	})
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://leetcode.com/graphql/", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", fmt.Sprintf("https://leetcode.com/problems/%s/", slug))
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Origin", "https://leetcode.com")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result struct {
		Data struct {
			Question *Question `json:"question"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql error: %s", result.Errors[0].Message)
	}

	return result.Data.Question, nil
}
