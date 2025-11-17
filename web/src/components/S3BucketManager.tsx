import { useCallback, useState } from 'react';
import { Card, Button } from '../ui';
import { useQuery, useMutation } from '../hooks';
import { awsApi } from '../api';
import type { S3BucketsResponse } from '../types';

interface S3BucketManagerProps {
  onSelectBucket?: (bucketName: string) => void;
  refreshTrigger?: number;
}

export const S3BucketManager: React.FC<S3BucketManagerProps> = ({
  onSelectBucket,
  refreshTrigger
}) => {
  const [selectedBucket, setSelectedBucket] = useState<string | null>(null);

  const fetchBuckets = useCallback(async () => {
    return awsApi.listS3Buckets();
  }, []);

  const { data, isLoading, isError, error, refetch } = useQuery<S3BucketsResponse>(
    fetchBuckets,
    { enabled: true }
  );

  const deleteMutation = useMutation(
    async (bucketName: string) => awsApi.deleteS3Bucket(bucketName),
    {
      onSuccess: () => {
        setSelectedBucket(null);
        refetch();
      },
    }
  );

  const handleSelectBucket = (bucketName: string) => {
    setSelectedBucket(bucketName);
    onSelectBucket?.(bucketName);
  };

  const handleDeleteBucket = (bucketName: string, e: React.MouseEvent) => {
    e.stopPropagation();
    if (window.confirm(`Are you sure you want to delete bucket "${bucketName}"? This action cannot be undone.`)) {
      deleteMutation.mutate(bucketName);
    }
  };

  if (isLoading) {
    return (
      <Card>
        <h3>S3 Buckets</h3>
        <p>Loading buckets...</p>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card>
        <h3>S3 Buckets</h3>
        <div className="error-message">
          Error: {error?.message || 'Failed to load buckets'}
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

      {deleteMutation.isError && (
        <div className="error-message">
          Error: {deleteMutation.error?.message || 'Failed to delete bucket'}
        </div>
      )}

      {buckets.length > 0 ? (
        <div className="bucket-list">
          {buckets.map((bucket) => (
            <div
              key={bucket.name}
              className={`bucket-item ${selectedBucket === bucket.name ? 'selected' : ''}`}
              onClick={() => handleSelectBucket(bucket.name)}
            >
              <div className="bucket-info">
                <strong>{bucket.name}</strong>
                <small>Created: {new Date(bucket.creationDate).toLocaleString()}</small>
              </div>
              <Button
                variant="danger"
                size="small"
                onClick={(e) => handleDeleteBucket(bucket.name, e)}
                disabled={deleteMutation.isLoading}
              >
                Delete
              </Button>
            </div>
          ))}
        </div>
      ) : (
        <p>No buckets found. Create your first bucket!</p>
      )}
    </Card>
  );
};
