import { useState } from 'react';
import { Button, Input, Card } from '../ui';
import { useMutation } from '../hooks';
import { awsApi } from '../api';
import type { CreateBucketRequest } from '../types';

interface S3BucketCreateFormProps {
  onSuccess?: () => void;
}

export const S3BucketCreateForm: React.FC<S3BucketCreateFormProps> = ({ onSuccess }) => {
  const [bucketName, setBucketName] = useState('');
  const [region, setRegion] = useState('us-east-1');

  const { mutate, isLoading, isError, error, isSuccess } = useMutation(
    async (data: CreateBucketRequest) => awsApi.createS3Bucket(data),
    {
      onSuccess: () => {
        setBucketName('');
        setRegion('us-east-1');
        onSuccess?.();
      },
    }
  );

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    mutate({ bucketName, region });
  };

  return (
    <Card>
      <h3>Create S3 Bucket</h3>
      <form onSubmit={handleSubmit} className="form">
        <Input
          label="Bucket Name"
          type="text"
          value={bucketName}
          onChange={(e) => setBucketName(e.target.value)}
          required
          placeholder="my-unique-bucket-name"
        />
        <div className="input-wrapper">
          <label className="input-label">Region</label>
          <select
            className="input"
            value={region}
            onChange={(e) => setRegion(e.target.value)}
          >
            <option value="us-east-1">US East (N. Virginia)</option>
            <option value="us-east-2">US East (Ohio)</option>
            <option value="us-west-1">US West (N. California)</option>
            <option value="us-west-2">US West (Oregon)</option>
            <option value="eu-west-1">EU (Ireland)</option>
            <option value="eu-central-1">EU (Frankfurt)</option>
            <option value="ap-southeast-1">Asia Pacific (Singapore)</option>
            <option value="ap-northeast-1">Asia Pacific (Tokyo)</option>
          </select>
        </div>

        {isError && (
          <div className="error-message">
            Error: {error?.message || 'Failed to create bucket'}
          </div>
        )}

        {isSuccess && (
          <div className="success-message">
            Bucket created successfully!
          </div>
        )}

        <Button type="submit" disabled={isLoading || !bucketName}>
          {isLoading ? 'Creating...' : 'Create Bucket'}
        </Button>
      </form>
    </Card>
  );
};
