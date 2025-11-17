import { useCallback } from 'react';
import { Card, Button } from '../ui';
import { useQuery } from '../hooks';
import { awsApi } from '../api';
import type { S3BucketsResponse } from '../types';

export const S3BucketsList: React.FC = () => {
  const fetchBuckets = useCallback(async () => {
    return awsApi.listS3Buckets();
  }, []);

  const { data, isLoading, isError, error, refetch } = useQuery<S3BucketsResponse>(
    fetchBuckets,
    { enabled: true }
  );

  if (isLoading) {
    return (
      <Card>
        <h3>S3 Buckets</h3>
        <p>Loading S3 buckets...</p>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card>
        <h3>S3 Buckets</h3>
        <div className="error-message">
          Error: {error?.message || 'Failed to load S3 buckets'}
        </div>
        <Button onClick={refetch}>Retry</Button>
      </Card>
    );
  }

  const buckets = data?.buckets || [];

  return (
    <Card>
      <div className="card-header">
        <h3>S3 Buckets ({data?.count || 0})</h3>
        <Button onClick={refetch} variant="secondary">
          Refresh
        </Button>
      </div>

      {buckets.length > 0 ? (
        <div className="table-container">
          <table className="data-table">
            <thead>
              <tr>
                <th>Bucket Name</th>
                <th>Creation Date</th>
              </tr>
            </thead>
            <tbody>
              {buckets.map((bucket) => (
                <tr key={bucket.name}>
                  <td>{bucket.name}</td>
                  <td>{new Date(bucket.creationDate).toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <p>No S3 buckets found.</p>
      )}
    </Card>
  );
};
