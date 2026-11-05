import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { NotesPage } from './pages/Notes.tsx'
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import 'csh-material-bootstrap/color-modes.js';
import { BrowserRouter, Routes, Route } from "react-router";


import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ForbiddenPage } from './pages/Forbidden.tsx';
import { NotFoundPage } from './pages/NotFound.tsx';
import { HomePage } from './pages/Home.tsx';
import { AdminPage } from './pages/Admin.tsx';
const queryClient = new QueryClient();

createRoot(document.getElementById('root')!).render(
  	<StrictMode>
		<QueryClientProvider client={queryClient}>
				<BrowserRouter basename="/csh">
				<Routes>
						<Route index element={<HomePage />}/>
						<Route path="/notes" element={<NotesPage/>}/>
						<Route path="/admin" element={<AdminPage/>}/>
						<Route path="/forbidden/:user" element={<ForbiddenPage/>}/>								
						<Route path="*" element={<NotFoundPage/>}/>
				</Routes>
			</BrowserRouter>
    	</QueryClientProvider>
  	</StrictMode>,
)
