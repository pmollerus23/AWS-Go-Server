import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './contexts';
import { ErrorBoundary } from './components';
import { Shell } from './shell';
import { HomePage, LoginPage, SignUpPage, ProfilePage } from './pages';

function App() {
  return (
    <ErrorBoundary>
      <BrowserRouter>
        <AuthProvider>
          <Routes>
            {/* Public routes without shell */}
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<SignUpPage />} />

            {/* Routes with shell */}
            <Route path="/" element={<Shell><HomePage /></Shell>} />
            <Route path="/profile" element={<Shell><ProfilePage /></Shell>} />

            {/* Catch all - redirect to home */}
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </AuthProvider>
      </BrowserRouter>
    </ErrorBoundary>
  );
}

export default App;
