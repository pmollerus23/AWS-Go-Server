import { useCallback } from 'react';
import { Card, Button } from '../ui';
import { useQuery } from '../hooks';
import { awsApi } from '../api';
import type { DynamoDBTablesResponse } from '../types';

export const DynamoDBTablesList: React.FC = () => {
  const fetchTables = useCallback(async () => {
    return awsApi.listDynamoDBTables();
  }, []);

  const { data, isLoading, isError, error, refetch } = useQuery<DynamoDBTablesResponse>(
    fetchTables,
    { enabled: true }
  );

  if (isLoading) {
    return (
      <Card>
        <h3>DynamoDB Tables</h3>
        <p>Loading DynamoDB tables...</p>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card>
        <h3>DynamoDB Tables</h3>
        <div className="error-message">
          Error: {error?.message || 'Failed to load DynamoDB tables'}
        </div>
        <Button onClick={refetch}>Retry</Button>
      </Card>
    );
  }

  const tables = data?.tables || [];

  return (
    <Card>
      <div className="card-header">
        <h3>DynamoDB Tables ({data?.count || 0})</h3>
        <Button onClick={refetch} variant="secondary">
          Refresh
        </Button>
      </div>

      {tables.length > 0 ? (
        <ul className="tables-list">
          {tables.map((table) => (
            <li key={table} className="table-item">
              {table}
            </li>
          ))}
        </ul>
      ) : (
        <p>No DynamoDB tables found.</p>
      )}
    </Card>
  );
};
