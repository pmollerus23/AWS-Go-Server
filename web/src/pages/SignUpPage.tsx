import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Input, Button } from '../ui';
import { authApi } from '../api';
import type { SignUpData, ConfirmSignUpData } from '../types';

export const SignUpPage: React.FC = () => {
  const navigate = useNavigate();
  const [step, setStep] = useState<'signup' | 'confirm'>('signup');
  const [signUpData, setSignUpData] = useState<SignUpData>({
    email: '',
    password: '',
    name: '',
  });
  const [confirmData, setConfirmData] = useState<ConfirmSignUpData>({
    email: '',
    code: '',
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [message, setMessage] = useState<string | null>(null);

  const handleSignUpSubmit = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    setError(null);
    setIsLoading(true);

    try {
      const response = await authApi.signUp(signUpData);
      setMessage(response.message);
      setConfirmData(prev => ({ ...prev, email: signUpData.email }));
      setStep('confirm');
    } catch (err: any) {
      console.error('Sign up failed:', err);
      setError(err?.message || 'Sign up failed. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleConfirmSubmit = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    setError(null);
    setIsLoading(true);

    try {
      const response = await authApi.confirmSignUp(confirmData);
      setMessage(response.message);
      // Navigate to login after successful confirmation
      setTimeout(() => {
        navigate('/login', { replace: true });
      }, 2000);
    } catch (err: any) {
      console.error('Confirmation failed:', err);
      setError(err?.message || 'Confirmation failed. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSignUpChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    setSignUpData(prev => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  };

  const handleConfirmChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    setConfirmData(prev => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  };

  return (
    <div className="signup-page">
      <div className="signup-container">
        {step === 'signup' ? (
          <>
            <h1>Sign Up for AWS Go Server</h1>

            <form onSubmit={handleSignUpSubmit} className="signup-form">
              <Input
                type="email"
                name="email"
                label="Email"
                value={signUpData.email}
                onChange={handleSignUpChange}
                required
                autoComplete="email"
                fullWidth
              />

              <Input
                type="text"
                name="name"
                label="Full Name"
                value={signUpData.name}
                onChange={handleSignUpChange}
                autoComplete="name"
                fullWidth
              />

              <Input
                type="password"
                name="password"
                label="Password"
                value={signUpData.password}
                onChange={handleSignUpChange}
                required
                autoComplete="new-password"
                helperText="Must be at least 8 characters"
                fullWidth
              />

              {error && <p className="error-message">{error}</p>}

              <Button type="submit" disabled={isLoading} fullWidth>
                {isLoading ? 'Signing up...' : 'Sign Up'}
              </Button>

              <div className="form-footer">
                <p>
                  Already have an account? <Link to="/login">Login</Link>
                </p>
              </div>
            </form>
          </>
        ) : (
          <>
            <h1>Confirm Your Email</h1>
            {message && <p className="success-message">{message}</p>}

            <form onSubmit={handleConfirmSubmit} className="confirm-form">
              <Input
                type="text"
                name="code"
                label="Verification Code"
                value={confirmData.code}
                onChange={handleConfirmChange}
                required
                helperText="Enter the code sent to your email"
                fullWidth
              />

              {error && <p className="error-message">{error}</p>}

              <Button type="submit" disabled={isLoading} fullWidth>
                {isLoading ? 'Confirming...' : 'Confirm Email'}
              </Button>

              <div className="form-footer">
                <p>
                  <button
                    type="button"
                    onClick={() => setStep('signup')}
                    className="btn-link"
                  >
                    Back to Sign Up
                  </button>
                </p>
              </div>
            </form>
          </>
        )}
      </div>
    </div>
  );
};
