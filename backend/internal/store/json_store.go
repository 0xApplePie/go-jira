package store

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"github.com/0xApplePie/go-jira/internal/models"
)

type JSONStore struct {
	filepath string
	tickets  map[string]*models.Ticket
	mu       sync.RWMutex
}

func NewJSONStore(filepath string) (*JSONStore, error) {
	store := &JSONStore{
		filepath: filepath,
		tickets:  make(map[string]*models.Ticket),
	}
	
	// Create directory if it doesn't exist
	dir := filepath[:len(filepath)-len("/tickets.json")]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// Try to load existing data
	if _, err := os.Stat(filepath); err == nil {
		data, err := os.ReadFile(filepath)
		if err != nil {
			return nil, err
		}
		
		// Only try to unmarshal if the file is not empty
		if len(data) > 0 {
			if err := json.Unmarshal(data, &store.tickets); err != nil {
				return nil, err
			}
		}
	}
	
	return store, nil
}

func (s *JSONStore) Add(ticket *models.Ticket) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.tickets[ticket.ID] = ticket
	return s.Save()
}

func (s *JSONStore) Get(id string) (*models.Ticket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	ticket, exists := s.tickets[id]
	if !exists {
		return nil, nil
	}
	return ticket, nil
}

func (s *JSONStore) List() ([]*models.Ticket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	tickets := make([]*models.Ticket, 0, len(s.tickets))
	for _, ticket := range s.tickets {
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

func (s *JSONStore) Update(ticket *models.Ticket) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.tickets[ticket.ID]; !exists {
		return errors.New("ticket not found")
	}
	
	s.tickets[ticket.ID] = ticket
	return s.Save()
}

func (s *JSONStore) Save() error {
	data, err := json.MarshalIndent(s.tickets, "", "    ")
	if err != nil {
		return err
	}
 
	return os.WriteFile(s.filepath, data, 0644)
} 