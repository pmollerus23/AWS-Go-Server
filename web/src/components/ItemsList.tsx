import { useCallback, useEffect } from 'react';
import { Card, Button } from '../ui';
import { useQuery } from '../hooks';
import { itemsApi } from '../api';
import type { Item } from '../types';

interface ItemsListProps {
  refreshTrigger?: number;
}

export const ItemsList: React.FC<ItemsListProps> = ({ refreshTrigger }) => {
  const fetchItems = useCallback(async () => {
    return itemsApi.getAll();
  }, []);

  const { data: items, isLoading, isError, error, refetch } = useQuery<Item[]>(
    fetchItems,
    { enabled: true }
  );

  // Trigger refetch when refreshTrigger changes
  useEffect(() => {
    if (refreshTrigger !== undefined && refreshTrigger > 0) {
      refetch();
    }
  }, [refreshTrigger, refetch]);

  if (isLoading) {
    return (
      <Card>
        <h3>Items</h3>
        <p>Loading items...</p>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card>
        <h3>Items</h3>
        <div className="error-message">
          Error: {error?.message || 'Failed to load items'}
        </div>
        <Button onClick={refetch}>Retry</Button>
      </Card>
    );
  }

  return (
    <Card>
      <div className="card-header">
        <h3>Items ({items?.length || 0})</h3>
        <Button onClick={refetch} variant="secondary">
          Refresh
        </Button>
      </div>

      {items && items.length > 0 ? (
        <div className="items-grid">
          {items.map((item) => (
            <div key={item.id} className="item-card">
              <h4>{item.name}</h4>
              <p>{item.description}</p>
              <small>ID: {item.id}</small>
            </div>
          ))}
        </div>
      ) : (
        <p>No items found. Create your first item!</p>
      )}
    </Card>
  );
};
