import { useCallback, useEffect } from 'react';
import { Card, Button } from '../ui';
import { useQuery, useMutation } from '../hooks';
import { awsApi } from '../api';
import type { S3ObjectsResponse } from '../types';

interface S3FilesListProps {
  bucketName: string | null;
  refreshTrigger?: number;
}

export const S3FilesList: React.FC<S3FilesListProps> = ({ bucketName, refreshTrigger }) => {
  const fetchObjects = useCallback(async () => {
    if (!bucketName) {
      return { objects: [], count: 0 };
    }
    return awsApi.listS3Objects(bucketName);
  }, [bucketName]);

  const { data, isLoading, isError, error, refetch } = useQuery<S3ObjectsResponse>(
    fetchObjects,
    { enabled: !!bucketName }
  );

  const deleteMutation = useMutation(
    async (key: string) => {
      if (!bucketName) throw new Error('No bucket selected');
      return awsApi.deleteS3Object(bucketName, key);
    },
    {
      onSuccess: () => {
        refetch();
      },
    }
  );

  // Trigger refetch when refreshTrigger changes
  useEffect(() => {
    if (refreshTrigger !== undefined && refreshTrigger > 0 && bucketName) {
      refetch();
    }
  }, [refreshTrigger, refetch, bucketName]);

  const handleDelete = (key: string) => {
    if (window.confirm(`Are you sure you want to delete "${key}"?`)) {
      deleteMutation.mutate(key);
    }
  };

  const handleDownload = (key: string) => {
    if (!bucketName) return;
    const downloadUrl = awsApi.downloadS3Object(bucketName, key);
    window.open(downloadUrl, '_blank');
  };

  if (!bucketName) {
    return (
      <Card>
        <h3>Files</h3>
        <p>Select a bucket to view files</p>
      </Card>
    );
  }

  if (isLoading) {
    return (
      <Card>
        <h3>Files in {bucketName}</h3>
        <p>Loading files...</p>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card>
        <h3>Files in {bucketName}</h3>
        <div className="error-message">
          Error: {error?.message || 'Failed to load files'}
        </div>
        <Button onClick={refetch}>Retry</Button>
      </Card>
    );
  }

  const objects = data?.objects || [];

  return (
    <Card>
      <div className="card-header">
        <h3>Files in {bucketName} ({data?.count || 0})</h3>
        <Button onClick={refetch} variant="secondary">
          Refresh
        </Button>
      </div>

      {deleteMutation.isError && (
        <div className="error-message">
          Error: {deleteMutation.error?.message || 'Failed to delete file'}
        </div>
      )}

      {objects.length > 0 ? (
        <div className="table-container">
          <table className="data-table">
            <thead>
              <tr>
                <th>File Name</th>
                <th>Size</th>
                <th>Last Modified</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {objects.map((obj) => (
                <tr key={obj.key}>
                  <td>{obj.key}</td>
                  <td>{(obj.size / 1024).toFixed(2)} KB</td>
                  <td>{new Date(obj.lastModified).toLocaleString()}</td>
                  <td>
                    <div className="action-buttons">
                      <Button
                        size="small"
                        variant="secondary"
                        onClick={() => handleDownload(obj.key)}
                      >
                        Download
                      </Button>
                      <Button
                        size="small"
                        variant="danger"
                        onClick={() => handleDelete(obj.key)}
                        disabled={deleteMutation.isLoading}
                      >
                        Delete
                      </Button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <p>No files in this bucket. Upload your first file!</p>
      )}
    </Card>
  );
};
