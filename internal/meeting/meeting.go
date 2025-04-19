package meeting

import (
	"sync"
	"time"
)

type Participant struct {
	ID       int64
	Username string
	JoinedAt time.Time
}

type Service struct {
	mu           sync.RWMutex
	participants map[int64]Participant
}

func NewService() *Service {
	return &Service{
		participants: make(map[int64]Participant),
	}
}

func (s *Service) AddParticipant(id int64, username string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.participants[id] = Participant{
		ID:       id,
		Username: username,
		JoinedAt: time.Now(),
	}
}

func (s *Service) GetParticipants() []Participant {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	result := make([]Participant, 0, len(s.participants))
	for _, p := range s.participants {
		result = append(result, p)
	}
	return result
}
