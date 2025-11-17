import { useCallback, useEffect } from 'react';
import { Card, Button } from '../ui';
import { useQuery } from '../hooks';
import { awsApi } from '../api';
import type { DynamoDBRecordsResponse } from '../types';

interface DynamoDBRecordsListProps {
  refreshTrigger?: number;
}

export const DynamoDBRecordsList: React.FC<DynamoDBRecordsListProps> = ({ refreshTrigger }) => {
  const fetchRecords = useCallback(async () => {
    return awsApi.listDynamoDBRecords();
  }, []);

  const { data, isLoading, isError, error, refetch } = useQuery<DynamoDBRecordsResponse>(
    fetchRecords,
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
        <h3>DynamoDB Records</h3>
        <p>Loading records...</p>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card>
        <h3>DynamoDB Records</h3>
        <div className="error-message">
          Error: {error?.message || 'Failed to load records'}
        </div>
        <Button onClick={refetch}>Retry</Button>
      </Card>
    );
  }

  const records = data?.records || [];

  return (
    <Card>
      <div className="card-header">
        <h3>DynamoDB Records ({data?.count || 0})</h3>
        <Button onClick={refetch} variant="secondary">
          Refresh
        </Button>
      </div>

      {records.length > 0 ? (
        <div className="table-container">
          <table className="data-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Name</th>
                <th>Updated At</th>
              </tr>
            </thead>
            <tbody>
              {records.map((record) => (
                <tr key={record.id}>
                  <td>{record.id}</td>
                  <td>{record.name}</td>
                  <td>{new Date(record.updated_at * 1000).toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <p>No records found. Create your first record!</p>
      )}
    </Card>
  );
};
