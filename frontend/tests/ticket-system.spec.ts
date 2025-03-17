import { test, expect } from '@playwright/test'

test.describe('Ticket System', () => {
  let tickets = []

  test.beforeEach(async ({ page }) => {
    tickets = [] // Reset tickets state for each test

    // Mock GraphQL responses
    await page.route('**/graphql', async (route) => {
      const request = route.request()
      const postData = JSON.parse(request.postData() || '{}')

      if (postData.query.includes('createTicket')) {
        const newTicket = {
          id: 'test-123',
          title: 'Test Ticket',
          description: 'Test Description',
          status: 'TODO',
          assignee: 'Test User',
        }
        tickets.push(newTicket)
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            data: { createTicket: newTicket },
          }),
        })
      } else if (postData.query.includes('deleteTicket')) {
        tickets = [] // Clear tickets after delete
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            data: { deleteTicket: true },
          }),
        })
      } else {
        // Query for tickets list
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            data: { tickets },
          }),
        })
      }
    })

    await page.goto('http://localhost:3000')
    // Wait for initial page load
    await page.waitForSelector('form')
  })

  test('should create and display a new ticket', async ({ page }) => {
    // Fill in the create ticket form
    await page.fill('input[placeholder="Title"]', 'Test Ticket')
    await page.fill('textarea[placeholder="Description"]', 'Test Description')
    await page.fill('input[placeholder="Assignee"]', 'Test User')

    // Submit the form
    await page.click('button:text("Create Ticket")')

    // Wait for the network request to complete
    await page.waitForResponse(
      (response) =>
        response.url().includes('/graphql') &&
        response.request().postData()?.includes('createTicket')
    )

    // Wait for and verify the new ticket appears
    await expect(page.getByText('Test Ticket')).toBeVisible()
  })

  test('should delete a ticket', async ({ page }) => {
    // Create a ticket first
    await page.fill('input[placeholder="Title"]', 'Test Ticket')
    await page.fill('textarea[placeholder="Description"]', 'Test Description')
    await page.fill('input[placeholder="Assignee"]', 'Test User')
    await page.click('button:text("Create Ticket")')

    // Wait for the ticket to appear
    await expect(page.getByText('Test Ticket')).toBeVisible()

    // Set up dialog handler before clicking delete
    page.on('dialog', (dialog) => dialog.accept())

    // Click delete button
    await page.click('button:text("Delete")')

    // Wait for the delete request to complete
    await page.waitForResponse(
      (response) =>
        response.url().includes('/graphql') &&
        response.request().postData()?.includes('deleteTicket')
    )

    // Verify the ticket is no longer visible
    await expect(page.getByText('Test Ticket')).not.toBeVisible()
  })
})
