import { useState } from 'react';
import { Button, Input, Card } from '../ui';
import { useMutation } from '../hooks';
import { itemsApi } from '../api';
import type { CreateItemRequest } from '../types';

interface ItemCreateFormProps {
  onSuccess?: () => void;
}

export const ItemCreateForm: React.FC<ItemCreateFormProps> = ({ onSuccess }) => {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');

  const { mutate, isLoading, isError, error, isSuccess } = useMutation(
    async (data: CreateItemRequest) => itemsApi.create(data),
    {
      onSuccess: () => {
        setName('');
        setDescription('');
        onSuccess?.();
      },
    }
  );

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    mutate({ name, description });
  };

  return (
    <Card>
      <h3>Create New Item</h3>
      <form onSubmit={handleSubmit} className="form">
        <Input
          label="Name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          placeholder="Enter item name"
        />
        <Input
          label="Description"
          type="text"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Enter item description"
        />

        {isError && (
          <div className="error-message">
            Error: {error?.message || 'Failed to create item'}
          </div>
        )}

        {isSuccess && (
          <div className="success-message">
            Item created successfully!
          </div>
        )}

        <Button type="submit" disabled={isLoading || !name}>
          {isLoading ? 'Creating...' : 'Create Item'}
        </Button>
      </form>
    </Card>
  );
};
