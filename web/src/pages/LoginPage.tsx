import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts';
import { Input, Button } from '../ui';
import type { LoginCredentials } from '../types';

export const LoginPage: React.FC = () => {
  const navigate = useNavigate();
  const { login, isLoading, error } = useAuth();
  const [credentials, setCredentials] = useState<LoginCredentials>({
    email: '',
    password: '',
  });
  const [loginError, setLoginError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    setLoginError(null);

    try {
      await login(credentials);
      // Navigate to home page on success
      navigate('/', { replace: true });
    } catch (err: any) {
      console.error('Login failed:', err);
      setLoginError(err?.message || 'Login failed. Please try again.');
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    setCredentials(prev => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  };

  return (
    <div className="login-page">
      <div className="login-container">
        <h1>Login to AWS Go Server</h1>

        <form onSubmit={handleSubmit} className="login-form">
          <Input
            type="email"
            name="email"
            label="Email"
            value={credentials.email}
            onChange={handleChange}
            required
            autoComplete="email"
            fullWidth
          />

          <Input
            type="password"
            name="password"
            label="Password"
            value={credentials.password}
            onChange={handleChange}
            required
            autoComplete="current-password"
            fullWidth
          />

          {(error || loginError) && (
            <p className="error-message">{error || loginError}</p>
          )}

          <Button type="submit" disabled={isLoading} fullWidth>
            {isLoading ? 'Logging in...' : 'Login'}
          </Button>

          <div className="form-footer">
            <p>
              Don't have an account? <Link to="/register">Sign Up</Link>
            </p>
            <p>
              <Link to="/forgot-password">Forgot Password?</Link>
            </p>
          </div>
        </form>
      </div>
    </div>
  );
};
