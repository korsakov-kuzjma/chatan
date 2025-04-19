package questionnaire

import (
	"sync"
)

type Question struct {
	Text    string
	Options []string
}

type Answer struct {
	UserID   int64
	Question string
	Response string
}

type Service struct {
	mu        sync.RWMutex
	questions []Question
	answers   []Answer
}

func NewService() *Service {
	return &Service{
		questions: []Question{
			{
				Text: "Как давно вы в программе?",
				Options: []string{
					"Менее месяца",
					"1-6 месяцев",
					"6-12 месяцев",
					"Более года",
				},
			},
			{
				Text: "Как часто вы посещаете встречи?",
				Options: []string{
					"Регулярно",
					"Иногда",
					"Редко",
				},
			},
		},
		answers: make([]Answer, 0),
	}
}

func (s *Service) GetQuestions() []Question {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.questions
}

func (s *Service) AddAnswer(userID int64, question, response string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.answers = append(s.answers, Answer{
		UserID:   userID,
		Question: question,
		Response: response,
	})
}

func (s *Service) GetAnswers(question string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []string
	for _, a := range s.answers {
		if a.Question == question {
			result = append(result, a.Response)
		}
	}
	return result
}
