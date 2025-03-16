import { gql } from '@apollo/client'

export const GET_TICKETS = gql`
  query GetTickets {
    tickets {
      id
      title
      description
      status
      assignee
      createdAt
      updatedAt
    }
  }
`

export const GET_TICKET = gql`
  query GetTicket($id: String!) {
    ticket(id: $id) {
      id
      title
      description
      status
      assignee
      createdAt
      updatedAt
    }
  }
`

export const CREATE_TICKET = gql`
  mutation CreateTicket(
    $title: String!
    $description: String!
    $assignee: String
  ) {
    createTicket(
      title: $title
      description: $description
      assignee: $assignee
    ) {
      id
      title
      description
      status
      assignee
    }
  }
`

export const UPDATE_TICKET = gql`
  mutation UpdateTicket(
    $id: String!
    $title: String
    $description: String
    $status: String
    $assignee: String
  ) {
    updateTicket(
      id: $id
      title: $title
      description: $description
      status: $status
      assignee: $assignee
    ) {
      id
      title
      description
      status
      assignee
      updatedAt
    }
  }
`

export const DELETE_TICKET = gql`
  mutation DeleteTicket($id: String!) {
    deleteTicket(id: $id)
  }
`
