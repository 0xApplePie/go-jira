'use client'

import { useState } from 'react'
import { useMutation } from '@apollo/client'
import { CREATE_TICKET, GET_TICKETS } from '@/graphql/tickets'

export default function CreateTicketForm() {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [assignee, setAssignee] = useState('')

  const [createTicket] = useMutation(CREATE_TICKET, {
    refetchQueries: [{ query: GET_TICKETS }],
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await createTicket({
        variables: {
          title,
          description,
          assignee: assignee || null,
        },
      })
      setTitle('')
      setDescription('')
      setAssignee('')
    } catch (error) {
      console.error('Error creating ticket:', error)
    }
  }

  return (
    <form
      onSubmit={handleSubmit}
      className='max-w-md mx-auto p-6 bg-white rounded-lg shadow-lg'
    >
      <h2 className='text-2xl font-bold mb-6'>Create New Ticket</h2>
      <div className='space-y-4'>
        <div>
          <label className='block text-sm font-medium text-gray-700'>
            Title
          </label>
          <input
            type='text'
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className='mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500'
            required
          />
        </div>
        <div>
          <label className='block text-sm font-medium text-gray-700'>
            Description
          </label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            className='mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500'
            rows={3}
            required
          />
        </div>
        <div>
          <label className='block text-sm font-medium text-gray-700'>
            Assignee
          </label>
          <input
            type='text'
            value={assignee}
            onChange={(e) => setAssignee(e.target.value)}
            className='mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500'
          />
        </div>
        <button
          type='submit'
          className='w-full bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2'
        >
          Create Ticket
        </button>
      </div>
    </form>
  )
}
