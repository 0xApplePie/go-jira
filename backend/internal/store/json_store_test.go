package store

import (
	"os"
	"testing"
	"time"

	"github.com/0xApplePie/go-jira/internal/models"
)

func TestJSONStore(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "tickets-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Initialize store
	store, err := NewJSONStore(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Test adding a ticket
	ticket := &models.Ticket{
		ID:          "test-123",
		Title:       "Test Ticket",
		Description: "Test Description",
		Status:      models.StatusTodo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Assignee:    nil,
	}

	t.Run("Add and Get Ticket", func(t *testing.T) {
		err := store.Add(ticket)
		if err != nil {
			t.Errorf("Failed to add ticket: %v", err)
		}

		got, err := store.Get(ticket.ID)
		if err != nil {
			t.Errorf("Failed to get ticket: %v", err)
		}
		if got == nil {
			t.Error("Expected to get ticket, got nil")
		}
		if got.ID != ticket.ID {
			t.Errorf("Got ticket ID %s, want %s", got.ID, ticket.ID)
		}
	})

	t.Run("List Tickets", func(t *testing.T) {
		tickets, err := store.List()
		if err != nil {
			t.Errorf("Failed to list tickets: %v", err)
		}
		if len(tickets) != 1 {
			t.Errorf("Expected 1 ticket, got %d", len(tickets))
		}
	})

	t.Run("Update Ticket", func(t *testing.T) {
		ticket.Status = models.StatusProgress
		err := store.Update(ticket)
		if err != nil {
			t.Errorf("Failed to update ticket: %v", err)
		}

		got, err := store.Get(ticket.ID)
		if err != nil {
			t.Errorf("Failed to get updated ticket: %v", err)
		}
		if got.Status != models.StatusProgress {
			t.Errorf("Got status %v, want %v", got.Status, models.StatusProgress)
		}
	})

	t.Run("Update Non-existent Ticket", func(t *testing.T) {
		nonExistentTicket := &models.Ticket{
			ID: "non-existent",
		}
		err := store.Update(nonExistentTicket)
		if err == nil {
			t.Error("Expected error when updating non-existent ticket, got nil")
		}
	})
} 