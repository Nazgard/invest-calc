package infrastructure

import (
	"context"
	"fmt"
	"sync"

	"invest-calc/internal/domain"
)

type MemoryStore struct {
	mu           sync.RWMutex
	calculations map[string]domain.CalculationResult
	counter      int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		calculations: make(map[string]domain.CalculationResult),
	}
}

func (s *MemoryStore) SaveCalculation(ctx context.Context, result domain.CalculationResult) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	id := fmt.Sprintf("calc-%d", s.counter)
	result.ID = id
	s.calculations[id] = result
	return id, nil
}

func (s *MemoryStore) GetCalculation(ctx context.Context, id string) (*domain.CalculationResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result, ok := s.calculations[id]
	if !ok {
		return nil, fmt.Errorf("calculation not found: %s", id)
	}
	return &result, nil
}

func (s *MemoryStore) ListCalculations(ctx context.Context) ([]domain.CalculationResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]domain.CalculationResult, 0, len(s.calculations))
	for _, r := range s.calculations {
		results = append(results, r)
	}
	return results, nil
}
