package statistics

import (
	"sync"
	"time"
)

type AttendanceRecord struct {
	UserID    int64
	Username  string
	MeetingAt time.Time
	Present   bool
}

type Service struct {
	mu      sync.RWMutex
	records []AttendanceRecord
}

func NewService() *Service {
	return &Service{
		records: make([]AttendanceRecord, 0),
	}
}

func (s *Service) RecordAttendance(userID int64, username string, meetingTime time.Time, present bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records = append(s.records, AttendanceRecord{
		UserID:    userID,
		Username:  username,
		MeetingAt: meetingTime,
		Present:   present,
	})
}

func (s *Service) GetUserStats(userID int64) (total, attended int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, record := range s.records {
		if record.UserID == userID {
			total++
			if record.Present {
				attended++
			}
		}
	}
	return
}

func (s *Service) GetGeneralStats() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := make(map[string]int)
	for _, record := range s.records {
		if record.Present {
			stats[record.Username]++
		}
	}
	return stats
}
