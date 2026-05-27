package ai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
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

// AIConfig 动态配置，从数据库加载
type AIConfig struct {
	Provider    string  `json:"provider"`
	APIKey      string  `json:"api_key"`
	BaseURL     string  `json:"base_url"`
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	RateLimit   int     `json:"rate_limit"`
}

// ===== OpenAI Provider =====

type OpenAIProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewOpenAIProvider(apiKey, baseURL string) *OpenAIProvider {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &OpenAIProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
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
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
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

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("stream API error (status %d): %s", resp.StatusCode, string(body))
			return
		}

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

// ===== Service (动态配置版) =====

type Service struct {
	db       *gorm.DB
	mu       sync.RWMutex
	provider AIProvider
	cfg      AIConfig

	rateMu   sync.Mutex
	reqCount int
	lastTime time.Time
}

// NewService 创建 AI 服务，从数据库加载配置
func NewService(db *gorm.DB) *Service {
	s := &Service{
		db:       db,
		lastTime: time.Now(),
	}
	s.loadConfig()
	return s
}

// loadConfig 从 system_configs 表加载 AI 配置
func (s *Service) loadConfig() {
	getConfig := func(key, defaultVal string) string {
		var sc struct{ Value string }
		if err := s.db.Table("system_configs").Where("`key` = ?", key).Select("value").Scan(&sc).Error; err != nil || sc.Value == "" {
			return defaultVal
		}
		return sc.Value
	}

	maxTokens, _ := strconv.Atoi(getConfig("ai_max_tokens", "2048"))
	if maxTokens <= 0 {
		maxTokens = 2048
	}
	temperature, _ := strconv.ParseFloat(getConfig("ai_temperature", "0.7"), 64)
	if temperature <= 0 {
		temperature = 0.7
	}
	rateLimit, _ := strconv.Atoi(getConfig("ai_rate_limit", "60"))
	if rateLimit <= 0 {
		rateLimit = 60
	}

	cfg := AIConfig{
		Provider:    getConfig("ai_provider", "openai"),
		APIKey:      getConfig("ai_api_key", ""),
		BaseURL:     getConfig("ai_base_url", ""),
		Model:       getConfig("ai_model", "gpt-4o"),
		MaxTokens:   maxTokens,
		Temperature: temperature,
		RateLimit:   rateLimit,
	}

	s.mu.Lock()
	s.cfg = cfg
	s.provider = s.createProvider(cfg)
	s.mu.Unlock()

	log.Printf("[AI] config loaded: provider=%s model=%s base_url=%s has_key=%v",
		cfg.Provider, cfg.Model, cfg.BaseURL, cfg.APIKey != "")
}

func (s *Service) createProvider(cfg AIConfig) AIProvider {
	switch cfg.Provider {
	case "deepseek":
		baseURL := cfg.BaseURL
		if baseURL == "" {
			baseURL = "https://api.deepseek.com/v1"
		}
		return NewOpenAIProvider(cfg.APIKey, baseURL)
	default: // openai 兼容
		return NewOpenAIProvider(cfg.APIKey, cfg.BaseURL)
	}
}

// ReloadConfig 供 admin 调用，重新加载配置
func (s *Service) ReloadConfig() {
	s.loadConfig()
}

// GetConfig 返回当前配置（隐藏 api_key 敏感信息）
func (s *Service) GetConfig() AIConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cfg := s.cfg
	if cfg.APIKey != "" {
		masked := cfg.APIKey
		if len(masked) > 8 {
			masked = masked[:4] + "****" + masked[len(masked)-4:]
		}
		cfg.APIKey = masked
	}
	return cfg
}

func (s *Service) rateLimit() {
	s.rateMu.Lock()
	defer s.rateMu.Unlock()

	now := time.Now()
	elapsed := now.Sub(s.lastTime)
	if elapsed >= time.Minute {
		s.reqCount = 0
		s.lastTime = now
	}

	s.mu.RLock()
	limit := s.cfg.RateLimit
	s.mu.RUnlock()

	if s.reqCount >= limit {
		time.Sleep(time.Minute - elapsed)
		s.reqCount = 0
		s.lastTime = time.Now()
	}

	s.reqCount++
}

func (s *Service) getProvider() AIProvider {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.provider
}

func (s *Service) getModel() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg.Model
}

func (s *Service) ChatCompletion(messages []Message) (*AIResponse, error) {
	s.rateLimit()
	return s.getProvider().ChatCompletion(messages, s.getModel())
}

func (s *Service) StreamChatCompletion(messages []Message) (<-chan string, error) {
	s.rateLimit()
	return s.getProvider().StreamChatCompletion(messages, s.getModel())
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
