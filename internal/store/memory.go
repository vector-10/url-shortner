package store

import (
	"sync"
	"errors"	
	"github.com/vector-10/url-shortner/internal/models"
)

 type MemoryStore struct {
	records map[string]*models.URLRecord
	mu sync.RWMutex
 }

 func NewMemoryStore() *MemoryStore {
   return &MemoryStore{
	records: make(map[string]*models.URLRecord),
   }
 }

 func(m *MemoryStore) Save(record *models.URLRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.records[record.Slug] = record
	return nil
}


func (m *MemoryStore) GetBySlug(slug string) (*models.URLRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	record, ok := m.records[slug]
	if !ok {
		return nil, errors.New("slug not found")
	}
	return record, nil
}

func (m *MemoryStore) IncrementClicks(slug string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	record, ok := m.records[slug]
	if !ok {
		return errors.New("slug not found")
	}
	record.Clicks++
	return nil
}

func (m *MemoryStore) Delete(slug string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.records[slug]
	if !ok {
		return errors.New("slug not found")		
	}
	delete(m.records, slug)
	return nil
}