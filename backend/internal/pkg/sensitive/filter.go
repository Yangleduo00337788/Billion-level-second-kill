package sensitive

import (
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

type SensitiveWord struct {
	ID        uint   `gorm:"primarykey"`
	Word      string `gorm:"size:100;uniqueIndex"`
	Level     int    `gorm:"default:1"` // 1=forbidden, 2=replace
	ReplaceTo string `gorm:"size:100"`
}

func (SensitiveWord) TableName() string {
	return "sensitive_words"
}

type Filter struct {
	db       *gorm.DB
	words    []SensitiveWord
	mu       sync.RWMutex
	lastLoad time.Time
}

func NewFilter(db *gorm.DB) *Filter {
	f := &Filter{db: db}
	f.loadWords()
	return f
}

func (f *Filter) loadWords() {
	var words []SensitiveWord
	f.db.Find(&words)
	f.mu.Lock()
	f.words = words
	f.lastLoad = time.Now()
	f.mu.Unlock()
}

func (f *Filter) reloadIfNeeded() {
	f.mu.RLock()
	needReload := time.Since(f.lastLoad) > 5*time.Minute
	f.mu.RUnlock()
	if needReload {
		f.loadWords()
	}
}

// CheckContent checks content against sensitive words
// Returns: filtered content, hasForbidden (true if contains forbidden words), forbidden word list
func (f *Filter) CheckContent(content string) (string, bool, []string) {
	f.reloadIfNeeded()
	f.mu.RLock()
	words := f.words
	f.mu.RUnlock()

	result := content
	hasForbidden := false
	var forbiddenWords []string

	for _, sw := range words {
		lowerContent := strings.ToLower(result)
		lowerWord := strings.ToLower(sw.Word)

		if strings.Contains(lowerContent, lowerWord) {
			if sw.Level == 1 {
				// Forbidden - reject
				hasForbidden = true
				forbiddenWords = append(forbiddenWords, sw.Word)
			} else if sw.Level == 2 {
				// Replace
				replaceTo := sw.ReplaceTo
				if replaceTo == "" {
					replaceTo = strings.Repeat("*", len(sw.Word))
				}
				result = strings.ReplaceAll(result, sw.Word, replaceTo)
				result = strings.ReplaceAll(result, strings.ToLower(sw.Word), replaceTo)
			}
		}
	}

	return result, hasForbidden, forbiddenWords
}
