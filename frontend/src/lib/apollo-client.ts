import { ApolloClient, InMemoryCache } from '@apollo/client'

const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export const client = new ApolloClient({
  uri: `${apiUrl}/graphql`,
  cache: new InMemoryCache(),
})
