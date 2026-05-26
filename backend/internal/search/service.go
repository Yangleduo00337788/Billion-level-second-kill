package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"inference-engine/internal/article"
	"inference-engine/internal/config"
	"inference-engine/internal/prompt"
	"inference-engine/internal/user"

	"github.com/elastic/go-elasticsearch/v8"
)

type ArticleDoc struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Summary   string `json:"summary"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type Service struct {
	client *elasticsearch.Client
	index  string
}

func NewService(cfg *config.ESConfig) *Service {
	address := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)

	esCfg := elasticsearch.Config{
		Addresses: []string{address},
	}
	if cfg.Username != "" && cfg.Password != "" {
		esCfg.Username = cfg.Username
		esCfg.Password = cfg.Password
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		log.Printf("failed to create elasticsearch client: %v", err)
		return &Service{index: "inference_engine"}
	}

	return &Service{
		client: client,
		index:  "inference_engine",
	}
}

func (s *Service) IndexArticle(a *article.Article) error {
	if s.client == nil {
		return nil
	}

	doc := ArticleDoc{
		ID:        a.ID,
		Title:     a.Title,
		Content:   a.Content,
		Summary:   a.Summary,
		UserID:    a.UserID,
		Username:  a.User.Username,
		Status:    a.Status,
		CreatedAt: a.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal article: %w", err)
	}

	_, err = s.client.Index(
		s.index,
		bytes.NewReader(data),
		s.client.Index.WithDocumentID(fmt.Sprintf("article_%d", a.ID)),
		s.client.Index.WithRefresh("wait_for"),
	)
	return err
}

func (s *Service) SearchArticles(query string, page, pageSize int) ([]ArticleDoc, int64, error) {
	if s.client == nil {
		return nil, 0, nil
	}

	from := (page - 1) * pageSize

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":  query,
							"fields": []string{"title^3", "content", "summary"},
						},
					},
					{
						"term": map[string]interface{}{
							"status": "published",
						},
					},
				},
			},
		},
		"from": from,
		"size": pageSize,
	}

	data, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, 0, err
	}

	resp, err := s.client.Search(
		s.client.Search.WithBody(bytes.NewReader(data)),
		s.client.Search.WithContext(context.Background()),
	)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source ArticleDoc `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	var docs []ArticleDoc
	for _, hit := range result.Hits.Hits {
		docs = append(docs, hit.Source)
	}

	return docs, result.Hits.Total.Value, nil
}

func (s *Service) SearchUsers(query string) ([]user.User, error) {
	return nil, nil
}

func (s *Service) SearchPrompts(query string, page, pageSize int) ([]prompt.Prompt, int64, error) {
	return nil, 0, nil
}

func (s *Service) DeleteIndex(id string) error {
	if s.client == nil {
		return nil
	}

	_, err := s.client.Delete(s.index, id)
	return err
}

type SearchResult struct {
	Articles []ArticleDoc    `json:"articles"`
	Users    []user.User     `json:"users"`
	Prompts  []prompt.Prompt `json:"prompts"`
}

func (s *Service) SearchAll(query string, page, pageSize int) (*SearchResult, error) {
	articles, _, _ := s.SearchArticles(query, page, pageSize)
	users, _ := s.SearchUsers(query)
	prompts, _, _ := s.SearchPrompts(query, page, pageSize)
	return &SearchResult{Articles: articles, Users: users, Prompts: prompts}, nil
}

func (s *Service) IsAvailable() bool {
	if s.client == nil {
		return false
	}

	resp, err := s.client.Ping()
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}
