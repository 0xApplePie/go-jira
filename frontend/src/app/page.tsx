import CreateTicketForm from '@/components/CreateTicketForm'
import TicketList from '@/components/TicketList'

export default function Home() {
  return (
    <main className='container mx-auto py-8'>
      <h1 className='text-3xl font-bold text-center mb-8'>
        Ticket Management System
      </h1>
      <CreateTicketForm />
      <div className='mt-12'>
        <h2 className='text-2xl font-semibold mb-4'>All Tickets</h2>
        <TicketList />
      </div>
    </main>
  )
}
