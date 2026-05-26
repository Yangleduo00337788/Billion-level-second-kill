package ai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"inference-engine/internal/config"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIResponse struct {
	Content string `json:"content"`
	Usage   Usage  `json:"usage"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type AIProvider interface {
	ChatCompletion(messages []Message, model string) (*AIResponse, error)
	StreamChatCompletion(messages []Message, model string) (<-chan string, error)
}

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	Stream      bool      `json:"stream"`
}

type chatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage Usage `json:"usage"`
}

type OpenAIProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey:  apiKey,
		baseURL: "https://api.openai.com/v1",
		client: &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: false,
			},
		},
	}
}

func (p *OpenAIProvider) ChatCompletion(messages []Message, model string) (*AIResponse, error) {
	reqBody := chatRequest{
		Model:       model,
		Messages:    messages,
		MaxTokens:   2048,
		Temperature: 0.7,
		Stream:      false,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", p.baseURL+"/chat/completions", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	return &AIResponse{
		Content: chatResp.Choices[0].Message.Content,
		Usage:   chatResp.Usage,
	}, nil
}

func (p *OpenAIProvider) StreamChatCompletion(messages []Message, model string) (<-chan string, error) {
	reqBody := chatRequest{
		Model:       model,
		Messages:    messages,
		MaxTokens:   2048,
		Temperature: 0.7,
		Stream:      true,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", p.baseURL+"/chat/completions", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	ch := make(chan string, 100)
	go func() {
		defer resp.Body.Close()
		defer close(ch)

		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				return
			}
			var streamResp struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
					FinishReason *string `json:"finish_reason"`
				} `json:"choices"`
			}
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue
			}
			if len(streamResp.Choices) > 0 {
				if streamResp.Choices[0].Delta.Content != "" {
					ch <- streamResp.Choices[0].Delta.Content
				}
				if streamResp.Choices[0].FinishReason != nil {
					return
				}
			}
		}
	}()

	return ch, nil
}

type DeepSeekProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewDeepSeekProvider(apiKey string) *DeepSeekProvider {
	return &DeepSeekProvider{
		apiKey:  apiKey,
		baseURL: "https://api.deepseek.com/v1",
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (p *DeepSeekProvider) ChatCompletion(messages []Message, model string) (*AIResponse, error) {
	reqBody := chatRequest{
		Model:       model,
		Messages:    messages,
		MaxTokens:   2048,
		Temperature: 0.7,
		Stream:      false,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", p.baseURL+"/chat/completions", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	return &AIResponse{
		Content: chatResp.Choices[0].Message.Content,
		Usage:   chatResp.Usage,
	}, nil
}

func (p *DeepSeekProvider) StreamChatCompletion(messages []Message, model string) (<-chan string, error) {
	reqBody := chatRequest{
		Model:       model,
		Messages:    messages,
		MaxTokens:   2048,
		Temperature: 0.7,
		Stream:      true,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", p.baseURL+"/chat/completions", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	ch := make(chan string, 100)
	go func() {
		defer resp.Body.Close()
		defer close(ch)

		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				return
			}
			var streamResp struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
					FinishReason *string `json:"finish_reason"`
				} `json:"choices"`
			}
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue
			}
			if len(streamResp.Choices) > 0 {
				if streamResp.Choices[0].Delta.Content != "" {
					ch <- streamResp.Choices[0].Delta.Content
				}
				if streamResp.Choices[0].FinishReason != nil {
					return
				}
			}
		}
	}()

	return ch, nil
}

type Service struct {
	provider AIProvider
	cfg      *config.AIConfig
	mu       sync.Mutex
	reqCount int
	lastTime time.Time
}

func NewService(cfg *config.AIConfig) *Service {
	var provider AIProvider
	switch cfg.Provider {
	case "deepseek":
		provider = NewDeepSeekProvider(cfg.APIKey)
	default:
		provider = NewOpenAIProvider(cfg.APIKey)
	}

	return &Service{
		provider: provider,
		cfg:      cfg,
		lastTime: time.Now(),
	}
}

func (s *Service) rateLimit() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(s.lastTime)
	if elapsed >= time.Minute {
		s.reqCount = 0
		s.lastTime = now
	}

	if s.reqCount >= s.cfg.RateLimit {
		time.Sleep(time.Minute - elapsed)
		s.reqCount = 0
		s.lastTime = time.Now()
	}

	s.reqCount++
}

func (s *Service) ChatCompletion(messages []Message) (*AIResponse, error) {
	s.rateLimit()
	return s.provider.ChatCompletion(messages, s.cfg.Model)
}

func (s *Service) StreamChatCompletion(messages []Message) (<-chan string, error) {
	s.rateLimit()
	return s.provider.StreamChatCompletion(messages, s.cfg.Model)
}

func (s *Service) GenerateTitle(content string) string {
	if content == "" {
		return ""
	}

	messages := []Message{
		{Role: "system", Content: "You are a helpful assistant that generates concise titles. Respond with only the title, no extra text."},
		{Role: "user", Content: fmt.Sprintf("Generate a concise title for this content: %s", truncate(content, 500))},
	}

	resp, err := s.ChatCompletion(messages)
	if err != nil {
		log.Printf("GenerateTitle error: %v", err)
		return ""
	}

	return strings.TrimSpace(resp.Content)
}

func (s *Service) GenerateSummary(content string) string {
	if content == "" {
		return ""
	}

	messages := []Message{
		{Role: "system", Content: "You are a helpful assistant that generates summaries. Respond with 1-2 sentences."},
		{Role: "user", Content: fmt.Sprintf("Generate a brief summary for this content: %s", truncate(content, 1000))},
	}

	resp, err := s.ChatCompletion(messages)
	if err != nil {
		log.Printf("GenerateSummary error: %v", err)
		return ""
	}

	return strings.TrimSpace(resp.Content)
}

func (s *Service) SuggestTags(content string) []string {
	if content == "" {
		return nil
	}

	messages := []Message{
		{Role: "system", Content: "You are a helpful assistant that suggests relevant tags. Return tags as a comma-separated list, no more than 5 tags."},
		{Role: "user", Content: fmt.Sprintf("Suggest relevant tags for this content: %s", truncate(content, 1000))},
	}

	resp, err := s.ChatCompletion(messages)
	if err != nil {
		log.Printf("SuggestTags error: %v", err)
		return nil
	}

	tags := strings.Split(resp.Content, ",")
	for i := range tags {
		tags[i] = strings.TrimSpace(tags[i])
	}
	return tags
}

func (s *Service) Chat(messages []Message) (*AIResponse, error) {
	return s.ChatCompletion(messages)
}

func (s *Service) StreamChat(messages []Message) (<-chan string, error) {
	return s.StreamChatCompletion(messages)
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen {
		return string(runes[:maxLen])
	}
	return s
}
