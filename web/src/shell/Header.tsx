import { Link } from 'react-router-dom';
import { useAuth } from '../contexts';
import type { VoidFunction } from '../types';

interface HeaderProps {
  onToggleSidebar: VoidFunction;
}

export const Header: React.FC<HeaderProps> = ({ onToggleSidebar }) => {
  const { user, isAuthenticated, logout } = useAuth();

  const handleLogout = async (): Promise<void> => {
    try {
      await logout();
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  return (
    <header className="header">
      <div className="header-left">
        <button onClick={onToggleSidebar} className="menu-button" aria-label="Toggle sidebar">
          ☰
        </button>
        <Link to="/" className="header-title">
          AWS Go Server
        </Link>
      </div>

      <nav className="header-nav">
        {isAuthenticated ? (
          <div className="user-menu">
            <span className="user-name">
              {user?.name || user?.username || user?.email}
            </span>
            <button onClick={handleLogout} className="logout-button">
              Logout
            </button>
          </div>
        ) : (
          <div className="auth-links">
            <Link to="/login">Login</Link>
            <Link to="/register">Register</Link>
          </div>
        )}
      </nav>
    </header>
  );
};
