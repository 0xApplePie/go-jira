'use client'

import { useQuery } from '@apollo/client'
import { GET_TICKETS } from '@/graphql/tickets'
import TicketCard from './TicketCard'

interface Ticket {
  id: string
  title: string
  description: string
  status: string
  assignee: string | null
}
export default function TicketList() {
  const { loading, error, data } = useQuery(GET_TICKETS)

  if (loading) return <div className='text-center p-4'>Loading...</div>
  if (error)
    return <div className='text-red-500 p-4'>Error: {error.message}</div>

  return (
    <div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 p-4'>
      {data.tickets.map((ticket: Ticket) => (
        <TicketCard key={ticket.id} ticket={ticket} />
      ))}
    </div>
  )
}
