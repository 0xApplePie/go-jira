package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/0xApplePie/go-jira/internal/models"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
    db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    // Test the connection
    if err := db.Ping(); err != nil {
        return nil, err
    }

    // Create tickets table if it doesn't exist
    if err := createTicketsTable(db); err != nil {
        return nil, err
    }

    return &PostgresStore{db: db}, nil
}

func createTicketsTable(db *sql.DB) error {
    query := `
        CREATE TABLE IF NOT EXISTS tickets (
            id VARCHAR(36) PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            status VARCHAR(50) NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL,
            updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
            assignee VARCHAR(255)
        );`
    
    _, err := db.Exec(query)
    return err
}

func (s *PostgresStore) Add(ticket *models.Ticket) error {
    query := `
        INSERT INTO tickets (id, title, description, status, created_at, updated_at, assignee)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`
    
    _, err := s.db.Exec(query,
        ticket.ID,
        ticket.Title,
        ticket.Description,
        ticket.Status,
        ticket.CreatedAt,
        ticket.UpdatedAt,
        ticket.Assignee,
    )
    return err
}

func (s *PostgresStore) Get(id string) (*models.Ticket, error) {
    query := `
        SELECT id, title, description, status, created_at, updated_at, assignee
        FROM tickets WHERE id = $1`
    
    ticket := &models.Ticket{}
    var assignee sql.NullString
    
    err := s.db.QueryRow(query, id).Scan(
        &ticket.ID,
        &ticket.Title,
        &ticket.Description,
        &ticket.Status,
        &ticket.CreatedAt,
        &ticket.UpdatedAt,
        &assignee,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    if assignee.Valid {
        ticket.Assignee = &assignee.String
    }
    
    return ticket, nil
}

func (s *PostgresStore) List() ([]*models.Ticket, error) {
    query := `
        SELECT id, title, description, status, created_at, updated_at, assignee
        FROM tickets`
    
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var tickets []*models.Ticket
    for rows.Next() {
        ticket := &models.Ticket{}
        var assignee sql.NullString
        
        err := rows.Scan(
            &ticket.ID,
            &ticket.Title,
            &ticket.Description,
            &ticket.Status,
            &ticket.CreatedAt,
            &ticket.UpdatedAt,
            &assignee,
        )
        if err != nil {
            return nil, err
        }
        
        if assignee.Valid {
            ticket.Assignee = &assignee.String
        }
        
        tickets = append(tickets, ticket)
    }
    
    return tickets, nil
}

func (s *PostgresStore) Update(ticket *models.Ticket) error {
    query := `
        UPDATE tickets
        SET title = $1, description = $2, status = $3, updated_at = $4, assignee = $5
        WHERE id = $6`
    
    result, err := s.db.Exec(query,
        ticket.Title,
        ticket.Description,
        ticket.Status,
        time.Now(),
        ticket.Assignee,
        ticket.ID,
    )
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return errors.New("ticket not found")
    }
    
    return nil
}

func (s *PostgresStore) Save() error {
    // No-op for PostgreSQL as changes are saved immediately
    return nil
}

func (s *PostgresStore) Close() error {
    return s.db.Close()
} 