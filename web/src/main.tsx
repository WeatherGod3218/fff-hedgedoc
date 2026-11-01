import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { NotesPage } from './pages/Notes.tsx'
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import 'csh-material-bootstrap/color-modes.js';


import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
const queryClient = new QueryClient();

createRoot(document.getElementById('root')!).render(
  	<StrictMode>
		<QueryClientProvider client={queryClient}>
      		<NotesPage />
    	</QueryClientProvider>
  	</StrictMode>,
)
