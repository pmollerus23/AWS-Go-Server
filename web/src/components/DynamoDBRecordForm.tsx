import { useState } from 'react';
import { Button, Input, Card } from '../ui';
import { useMutation } from '../hooks';
import { awsApi } from '../api';
import type { UpsertRecordRequest } from '../types';

interface DynamoDBRecordFormProps {
  onSuccess?: () => void;
}

export const DynamoDBRecordForm: React.FC<DynamoDBRecordFormProps> = ({ onSuccess }) => {
  const [id, setId] = useState('');
  const [name, setName] = useState('');

  const { mutate, isLoading, isError, error, isSuccess } = useMutation(
    async (data: UpsertRecordRequest) => awsApi.upsertDynamoDBRecord(data),
    {
      onSuccess: () => {
        setId('');
        setName('');
        onSuccess?.();
      },
    }
  );

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    mutate({
      id: parseInt(id),
      name,
      updated_at: Math.floor(Date.now() / 1000), // Unix timestamp in seconds
    });
  };

  return (
    <Card>
      <h3>Upsert DynamoDB Record</h3>
      <form onSubmit={handleSubmit} className="form">
        <Input
          label="ID"
          type="number"
          value={id}
          onChange={(e) => setId(e.target.value)}
          required
          placeholder="Enter record ID"
        />
        <Input
          label="Name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          placeholder="Enter record name"
        />

        {isError && (
          <div className="error-message">
            Error: {error?.message || 'Failed to upsert record'}
          </div>
        )}

        {isSuccess && (
          <div className="success-message">
            Record upserted successfully!
          </div>
        )}

        <Button type="submit" disabled={isLoading || !id || !name}>
          {isLoading ? 'Saving...' : 'Upsert Record'}
        </Button>
      </form>
    </Card>
  );
};
