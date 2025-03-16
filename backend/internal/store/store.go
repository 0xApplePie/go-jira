package store

import (
	"github.com/0xApplePie/go-jira/internal/models"
)

type TicketStore interface {
	Add(ticket *models.Ticket) error
	Get(id string) (*models.Ticket, error)
	List() ([]*models.Ticket, error)
	Update(ticket *models.Ticket) error
	Save() error
} 