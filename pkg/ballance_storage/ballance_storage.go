package ballance_storage

import (
	"sync"
)

type BallanceStorage interface {
	GetBallance(address string) int64
	AddBallance(address string, value int64) (int64, error)
	SubBallance(address string, value int64) (int64, error)
	Transfer(sender string, reciver string, value int64) error
	Confirm() error
	Reject() error
}

type BallanceStorageMemory struct {
	ballancePool map[string]int64
	txPool       map[string]int64
	mu           sync.Mutex
}

func NewMemoryStorage() BallanceStorage {
	var storage = BallanceStorageMemory{
		ballancePool: make(map[string]int64),
		txPool:       make(map[string]int64),
	}

	return &storage
}

func (s *BallanceStorageMemory) Confirm() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range s.txPool {
		if _, ok := s.ballancePool[k]; !ok {
			s.ballancePool[k] = 0
		}
		s.ballancePool[k] += v
	}
	for k := range s.txPool {
		delete(s.txPool, k)
	}
	return nil
}

func (s *BallanceStorageMemory) Reject() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k := range s.txPool {
		delete(s.txPool, k)
	}
	return nil
}

func (p *BallanceStorageMemory) GetBallance(address string) int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	ballance, ok := p.ballancePool[address]
	if !ok {
		ballance = 0
	}
	tempTx, ok := p.txPool[address]
	if !ok {
		tempTx = 0
	}

	return ballance + tempTx
}

func (s *BallanceStorageMemory) addBallance(address string, value int64) int64 {
	if _, ok := s.txPool[address]; !ok {
		s.txPool[address] = 0
	}
	s.txPool[address] += value
	return s.txPool[address]
}

func (s *BallanceStorageMemory) subBallance(address string, value int64) int64 {
	if _, ok := s.txPool[address]; !ok {
		s.txPool[address] = 0
	}
	s.txPool[address] -= value
	return s.txPool[address]
}

func (s *BallanceStorageMemory) AddBallance(address string, value int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addBallance(address, value), nil
}

func (s *BallanceStorageMemory) SubBallance(address string, value int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.subBallance(address, value), nil
}

func (s *BallanceStorageMemory) Transfer(sender string, reciver string, value int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subBallance(sender, value)
	s.addBallance(reciver, value)
	return nil
}
