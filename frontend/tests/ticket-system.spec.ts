import { test, expect } from '@playwright/test'

test.describe('Ticket System', () => {
  test('should create and delete a ticket', async ({ page }) => {
    // Navigate to the app
    await page.goto('/')

    // Debug: Log what we see on the page
    console.log('Current URL:', page.url())

    // Create a ticket with more explicit waiting
    const titleInput = await page.waitForSelector('input[type="text"]')
    await titleInput.type('Integration Test Ticket')

    const descriptionInput = await page.waitForSelector('textarea')
    await descriptionInput.type('Testing with real backend')

    const assigneeInput = await page.waitForSelector(
      'input[type="text"]:nth-of-type(2)'
    )
    await assigneeInput.type('Tester')

    // Find and click the create button
    const createButton = await page.waitForSelector(
      'button:has-text("Create Ticket")'
    )
    await createButton.click()

    // Wait for the new ticket to appear
    const ticketTitle = await page.waitForSelector(
      'text=Integration Test Ticket'
    )
    await expect(ticketTitle).toBeVisible()

    // Set up dialog handler BEFORE finding the delete button
    page.once('dialog', async (dialog) => {
      console.log('Dialog appeared:', dialog.message())
      await dialog.accept()
    })

    // Find and click the delete button with more explicit waiting
    const deleteButton = await page.waitForSelector('button:has-text("Delete")')
    await deleteButton.click()

    // Wait for the ticket to be removed
    await expect(page.locator('text=Integration Test Ticket')).not.toBeVisible()
  })
})
