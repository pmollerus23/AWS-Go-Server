import { useState } from 'react';
import { ItemCreateForm, ItemsList } from '../components';

export const ItemsPage: React.FC = () => {
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const handleItemCreated = () => {
    // Trigger refresh of items list
    setRefreshTrigger(prev => prev + 1);
  };

  return (
    <div className="items-page">
      <h1>Items Management</h1>
      <p>Create and manage items in your application.</p>

      <div className="page-grid">
        <div className="grid-item">
          <ItemCreateForm onSuccess={handleItemCreated} />
        </div>
        <div className="grid-item grid-span-2">
          <ItemsList refreshTrigger={refreshTrigger} />
        </div>
      </div>
    </div>
  );
};
