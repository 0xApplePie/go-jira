package graphql

import (
	"github.com/0xApplePie/go-jira/internal/models"
	"github.com/0xApplePie/go-jira/internal/store"
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
						ticket := &models.Ticket{
							Title:       p.Args["title"].(string),
							Description: p.Args["description"].(string),
							Status:      models.StatusTodo,
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
