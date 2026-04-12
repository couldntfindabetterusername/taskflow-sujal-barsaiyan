import { CssBaseline, ThemeProvider, createTheme, Box, Container, Typography } from '@mui/material'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'

const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
  },
})

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 5 * 60 * 1000,
    },
  },
})

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <BrowserRouter>
          <Container maxWidth="lg">
            <Box sx={{ my: 4 }}>
              <Typography variant="h2" component="h1" gutterBottom align="center">
                TaskFlow
              </Typography>
              <Typography variant="h5" component="h2" gutterBottom align="center" color="text.secondary">
                Modern Task Management System
              </Typography>
              <Box sx={{ mt: 4, p: 3, bgcolor: 'background.paper', borderRadius: 2, boxShadow: 1 }}>
                <Typography variant="body1" align="center">
                  Application is ready for development
                </Typography>
              </Box>
            </Box>
          </Container>
        </BrowserRouter>
      </ThemeProvider>
    </QueryClientProvider>
  )
}

export default App
