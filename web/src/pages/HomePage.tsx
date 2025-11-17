import { Link } from 'react-router-dom';
import { useCallback } from 'react';
import { useAuth } from '../contexts';
import { useQuery } from '../hooks';
import { itemsApi } from '../api';
import type { Item } from '../types';

export const HomePage: React.FC = () => {
  const { user, isAuthenticated } = useAuth();

  // Memoize the query function to prevent infinite loops
  const fetchItems = useCallback(async () => {
    return itemsApi.getAll();
  }, []);

  // Fetch items using useQuery hook
  const {
    data: items,
    isLoading,
    isError,
    error,
    refetch,
  } = useQuery<Item[]>(fetchItems, {
    enabled: isAuthenticated,
  });

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

      <section className="dashboard-section">
        <h2>Your Items</h2>

        {isLoading && <p>Loading items...</p>}

        {isError && (
          <div className="error-container">
            <p>Error loading items: {error?.message}</p>
            <button onClick={refetch} className="btn btn-secondary">
              Retry
            </button>
          </div>
        )}

        {!isLoading && !isError && (
          <>
            {items && items.length > 0 ? (
              <div className="items-list">
                {items.map((item) => (
                  <div key={item.id} className="item-card">
                    <h3>{item.name}</h3>
                    <p>{item.description}</p>
                  </div>
                ))}
              </div>
            ) : (
              <p>No items found. Create your first item!</p>
            )}
          </>
        )}
      </section>

      <section className="user-info">
        <h2>User Information</h2>
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
      </section>
    </div>
  );
};
