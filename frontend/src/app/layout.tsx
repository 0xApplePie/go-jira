import './globals.css'
import { Metadata } from 'next'
import Providers from '@/app/providers'

export const metadata: Metadata = {
  title: 'Jira Clone',
  description: 'Jira Clone',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang='en'>
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  )
}
