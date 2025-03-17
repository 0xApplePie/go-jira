import { render, screen, fireEvent } from '@testing-library/react'
import { MockedProvider } from '@apollo/client/testing'
import TicketCard from '../TicketCard'
import { DELETE_TICKET } from '@/graphql/tickets'
import '@testing-library/jest-dom'

const mockTicket = {
  id: '1',
  title: 'Test Ticket',
  description: 'Test Description',
  status: 'TODO',
  assignee: 'Test User',
}

const mocks = [
  {
    request: {
      query: DELETE_TICKET,
      variables: { id: '1' },
    },
    result: {
      data: {
        deleteTicket: true,
      },
    },
  },
]

const renderWithApollo = (component: React.ReactNode) => {
  return render(
    <MockedProvider mocks={mocks} addTypename={false}>
      {component}
    </MockedProvider>
  )
}

describe('TicketCard', () => {
  it('renders ticket information', () => {
    renderWithApollo(<TicketCard ticket={mockTicket} />)

    expect(screen.getByText('Test Ticket')).toBeInTheDocument()
    expect(screen.getByText('Test Description')).toBeInTheDocument()
    expect(screen.getByText('Test User')).toBeInTheDocument()
    expect(screen.getByText('TODO')).toBeInTheDocument()
  })

  it('handles delete action', async () => {
    const mockConfirm = jest.fn(() => true)
    window.confirm = mockConfirm

    renderWithApollo(<TicketCard ticket={mockTicket} />)

    fireEvent.click(screen.getByText('Delete'))
    expect(mockConfirm).toHaveBeenCalled()
  })
})
