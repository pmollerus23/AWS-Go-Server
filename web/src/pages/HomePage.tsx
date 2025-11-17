import { Link } from 'react-router-dom';
import { useAuth } from '../contexts';
import { Card } from '../ui';

export const HomePage: React.FC = () => {
  const { user, isAuthenticated } = useAuth();

  if (!isAuthenticated) {
    return (
      <div className="home-page">
        <h1>Welcome to AWS Go Server</h1>
        <p>Please log in to access your dashboard.</p>
        <Link to="/login" className="btn btn-primary">
          Login
        </Link>
      </div>
    );
  }

  return (
    <div className="home-page">
      <h1>Welcome back, {user?.name || user?.username || user?.email}!</h1>
      <p>Manage your items and AWS resources from this dashboard.</p>

      <div className="page-grid">
        <Card elevation="medium">
          <h3>📦 Items Management</h3>
          <p>Create and manage items in your application.</p>
          <Link to="/items" className="btn btn-primary">
            Go to Items
          </Link>
        </Card>

        <Card elevation="medium">
          <h3>☁️ AWS Resources</h3>
          <p>View and manage your AWS S3 buckets and DynamoDB tables.</p>
          <Link to="/aws" className="btn btn-primary">
            Go to AWS Resources
          </Link>
        </Card>

        <Card elevation="medium">
          <h3>👤 Profile</h3>
          <p>View and update your profile information.</p>
          <Link to="/profile" className="btn btn-primary">
            Go to Profile
          </Link>
        </Card>
      </div>

      <section className="user-info">
        <h2>Quick Info</h2>
        <Card>
          <div className="user-details">
            <p>
              <strong>Email:</strong> {user?.email}
            </p>
            <p>
              <strong>Username:</strong> {user?.username}
            </p>
            {user?.name && (
              <p>
                <strong>Name:</strong> {user.name}
              </p>
            )}
            {user?.groups && user.groups.length > 0 && (
              <p>
                <strong>Groups:</strong> {user.groups.join(', ')}
              </p>
            )}
          </div>
        </Card>
      </section>
    </div>
  );
};
