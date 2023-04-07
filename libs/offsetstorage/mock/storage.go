package mock

import (
	"context"
	"sync"
)

type MockOffsetStorage struct {
	state map[int32]int64
	mtx   *sync.RWMutex
}

func New() *MockOffsetStorage {
	return &MockOffsetStorage{
		state: make(map[int32]int64),
		mtx:   new(sync.RWMutex),
	}
}

func (s *MockOffsetStorage) GetOffset(ctx context.Context, partition int32) (int64, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.state[partition], nil
}

func (s *MockOffsetStorage) SetOffset(ctx context.Context, partition int32, value int64) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.state[partition] = value
	return nil
}
