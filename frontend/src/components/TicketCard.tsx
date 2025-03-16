'use client'

import { useMutation } from '@apollo/client'
import { DELETE_TICKET, GET_TICKETS } from '@/graphql/tickets'

interface Ticket {
  id: string
  title: string
  description: string
  status: string
  assignee: string | null
}

interface TicketCardProps {
  ticket: Ticket
}

export default function TicketCard({ ticket }: TicketCardProps) {
  const [deleteTicket] = useMutation(DELETE_TICKET, {
    refetchQueries: [{ query: GET_TICKETS }],
  })

  const handleDelete = async () => {
    if (confirm('Are you sure you want to delete this ticket?')) {
      try {
        await deleteTicket({ variables: { id: ticket.id } })
      } catch (error) {
        console.error('Error deleting ticket:', error)
      }
    }
  }

  return (
    <div className='bg-white shadow-lg rounded-lg p-6 hover:shadow-xl transition-shadow'>
      <h3 className='text-xl font-semibold mb-2'>{ticket.title}</h3>
      <p className='text-gray-600 mb-4'>{ticket.description}</p>
      <div className='flex justify-between items-center'>
        <span
          className={`px-3 py-1 rounded-full text-sm ${
            ticket.status === 'TODO'
              ? 'bg-red-100 text-red-800'
              : ticket.status === 'PROGRESS'
              ? 'bg-yellow-100 text-yellow-800'
              : 'bg-green-100 text-green-800'
          }`}
        >
          {ticket.status}
        </span>
        <span className='text-gray-500'>{ticket.assignee || 'Unassigned'}</span>
      </div>
      <div className='mt-4 flex justify-end space-x-2'>
        <button
          onClick={handleDelete}
          className='bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600'
        >
          Delete
        </button>
      </div>
    </div>
  )
}
