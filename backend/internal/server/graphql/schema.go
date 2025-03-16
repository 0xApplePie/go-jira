package graphql

import (
	"time"

	"errors"

	"github.com/0xApplePie/go-jira/internal/models"
	"github.com/0xApplePie/go-jira/internal/store"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

var ticketType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Ticket",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"assignee": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

func NewSchema(store store.TicketStore) (graphql.Schema, error) {
	queryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"ticket": &graphql.Field{
					Type: ticketType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id := p.Args["id"].(string)
						return store.Get(id)
					},
				},
				"tickets": &graphql.Field{
					Type: graphql.NewList(ticketType),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return store.List()
					},
				},
			},
		},
	)

	mutationType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"createTicket": &graphql.Field{
					Type: ticketType,
					Args: graphql.FieldConfigArgument{
						"title": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"description": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"assignee": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						now := time.Now()
						ticket := &models.Ticket{
							ID:          uuid.New().String(),
							Title:       p.Args["title"].(string),
							Description: p.Args["description"].(string),
							Status:      models.StatusTodo,
							CreatedAt:   now,
							UpdatedAt:   now,
						}
						if assignee, ok := p.Args["assignee"].(string); ok {
							ticket.Assignee = &assignee
						}
						
						if err := store.Add(ticket); err != nil {
							return nil, err
						}
						return ticket, nil
					},
				},
				"updateTicket": &graphql.Field{
					Type: ticketType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"title": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"description": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"status": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"assignee": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id := p.Args["id"].(string)
						
						ticket, err := store.Get(id)
						if err != nil {
							return nil, err
						}
						if ticket == nil {
							return nil, errors.New("ticket not found")
						}

						if title, ok := p.Args["title"].(string); ok {
							ticket.Title = title
						}
						if description, ok := p.Args["description"].(string); ok {
							ticket.Description = description
						}
						if status, ok := p.Args["status"].(string); ok {
							newStatus, err := models.ParseStatus(status)
							if err != nil {
								return nil, err
							}
							ticket.Status = newStatus
						}
						if assignee, ok := p.Args["assignee"].(string); ok {
							ticket.Assignee = &assignee
						}
						
						ticket.UpdatedAt = time.Now()
						
						if err := store.Update(ticket); err != nil {
							return nil, err
						}
						
						return ticket, nil
					},
				},
				"deleteTicket": &graphql.Field{
					Type: graphql.Boolean,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id := p.Args["id"].(string)
						return store.Delete(id)
					},
				},
			},
		},
	)

	return graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)
}
