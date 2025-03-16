package main

import (
    "fmt"
    "os"
    "time"

    "github.com/google/uuid"
    "github.com/spf13/cobra"
    "github.com/0xApplePie/go-jira/internal/models"
    "github.com/0xApplePie/go-jira/internal/store"
)

var rootCmd = &cobra.Command{
    Use:   "jira",
    Short: "A Jira-like CLI ticket management system",
    Long:  `A simple CLI for managing tickets, similar to Jira. Create, list, and manage tickets easily.`,
}

var (
    ticketStore store.TicketStore
    dataFile    = "data/tickets.json"
)

func init() {
    var err error
    ticketStore, err = store.NewJSONStore(dataFile)
    if err != nil {
        fmt.Printf("Error initializing store: %v\n", err)
        os.Exit(1)
    }

    // Create Command
    var createCmd = &cobra.Command{
        Use:   "create",
        Short: "Create a new ticket",
        Run:   createTicket,
    }
    createCmd.Flags().StringP("title", "t", "", "Ticket title")
    createCmd.Flags().StringP("description", "d", "", "Ticket description")
    createCmd.Flags().StringP("assignee", "a", "", "Ticket assignee")
    createCmd.MarkFlagRequired("title")
    createCmd.MarkFlagRequired("description")

    // List Command
    var listCmd = &cobra.Command{
        Use:   "list",
        Short: "List all tickets",
        Run:   listTickets,
    }

    // View Command
    var viewCmd = &cobra.Command{
        Use:   "view [id]",
        Short: "View a specific ticket",
        Args:  cobra.ExactArgs(1),
        Run:   viewTicket,
    }

    // Update Command
    var updateCmd = &cobra.Command{
        Use:   "update [id]",
        Short: "Update a ticket",
        Args:  cobra.ExactArgs(1),
        Run:   updateTicket,
    }
    updateCmd.Flags().StringP("title", "t", "", "New ticket title")
    updateCmd.Flags().StringP("description", "d", "", "New ticket description")
    updateCmd.Flags().StringP("status", "s", "", "New status (TODO, PROGRESS, DONE)")
    updateCmd.Flags().StringP("assignee", "a", "", "New assignee")

    rootCmd.AddCommand(createCmd, listCmd, viewCmd, updateCmd)
}

func createTicket(cmd *cobra.Command, args []string) {
    title, _ := cmd.Flags().GetString("title")
    description, _ := cmd.Flags().GetString("description")
    assignee, _ := cmd.Flags().GetString("assignee")

    var assigneePtr *string
    if assignee != "" {
        assigneePtr = &assignee
    }

    now := time.Now()
    ticket := &models.Ticket{
        ID:          uuid.New().String(),
        Title:       title,
        Description: description,
        Status:      models.StatusTodo,
        CreatedAt:   now,
        UpdatedAt:   now,
        Assignee:    assigneePtr,
    }

    if err := ticketStore.Add(ticket); err != nil {
        fmt.Printf("Error creating ticket: %v\n", err)
        return
    }

    fmt.Printf("Ticket created successfully with ID: %s\n", ticket.ID)
}

func listTickets(cmd *cobra.Command, args []string) {
    tickets, err := ticketStore.List()
    if err != nil {
        fmt.Printf("Error listing tickets: %v\n", err)
        return
    }

    if len(tickets) == 0 {
        fmt.Println("No tickets found.")
        return
    }

    fmt.Println("ID | Title | Status | Assignee")
    fmt.Println("------------------------")
    for _, t := range tickets {
        assignee := "Unassigned"
        if t.Assignee != nil {
            assignee = *t.Assignee
        }
        fmt.Printf("%s | %s | %s | %s\n", t.ID, t.Title, t.Status, assignee)
    }
}

func viewTicket(cmd *cobra.Command, args []string) {
    ticket, err := ticketStore.Get(args[0])
    if err != nil {
        fmt.Printf("Error retrieving ticket: %v\n", err)
        return
    }
    if ticket == nil {
        fmt.Println("Ticket not found.")
        return
    }

    fmt.Printf("ID: %s\n", ticket.ID)
    fmt.Printf("Title: %s\n", ticket.Title)
    fmt.Printf("Description: %s\n", ticket.Description)
    fmt.Printf("Status: %s\n", ticket.Status)
    fmt.Printf("Created: %s\n", ticket.CreatedAt.Format(time.RFC3339))
    fmt.Printf("Updated: %s\n", ticket.UpdatedAt.Format(time.RFC3339))
    if ticket.Assignee != nil {
        fmt.Printf("Assignee: %s\n", *ticket.Assignee)
    } else {
        fmt.Println("Assignee: Unassigned")
    }
}

func updateTicket(cmd *cobra.Command, args []string) {
    ticket, err := ticketStore.Get(args[0])
    if err != nil {
        fmt.Printf("Error retrieving ticket: %v\n", err)
        return
    }
    if ticket == nil {
        fmt.Println("Ticket not found.")
        return
    }

    if title, _ := cmd.Flags().GetString("title"); title != "" {
        ticket.Title = title
    }
    if desc, _ := cmd.Flags().GetString("description"); desc != "" {
        ticket.Description = desc
    }
    if status, _ := cmd.Flags().GetString("status"); status != "" {
        newStatus, err := models.ParseStatus(status)
        if err != nil {
            fmt.Printf("Invalid status: %v\n", err)
            return
        }
        ticket.Status = newStatus
    }
    if assignee, _ := cmd.Flags().GetString("assignee"); assignee != "" {
        ticket.Assignee = &assignee
    }

    ticket.UpdatedAt = time.Now()

    if err := ticketStore.Update(ticket); err != nil {
        fmt.Printf("Error updating ticket: %v\n", err)
        return
    }

    fmt.Println("Ticket updated successfully.")
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}