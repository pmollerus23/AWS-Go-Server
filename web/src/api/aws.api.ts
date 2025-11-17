import { apiClient } from './client';
import { storage, STORAGE_KEYS } from '../lib';
import type {
  S3BucketsResponse,
  CreateBucketRequest,
  CreateBucketResponse,
  DeleteBucketResponse,
  S3ObjectsResponse,
  UploadObjectResponse,
  DeleteObjectResponse,
  DynamoDBTablesResponse,
  DynamoDBRecordsResponse,
  UpsertRecordRequest,
  UpsertRecordResponse,
} from '../types/aws.types';

export const awsApi = {
  // S3 Bucket Operations
  listS3Buckets: async (): Promise<S3BucketsResponse> => {
    return apiClient.get<S3BucketsResponse>('/api/v1/aws/s3/buckets');
  },

  createS3Bucket: async (data: CreateBucketRequest): Promise<CreateBucketResponse> => {
    return apiClient.post<CreateBucketResponse>('/api/v1/aws/s3/buckets', data);
  },

  deleteS3Bucket: async (bucketName: string): Promise<DeleteBucketResponse> => {
    return apiClient.delete<DeleteBucketResponse>(`/api/v1/aws/s3/buckets/${bucketName}`);
  },

  // S3 Object Operations
  listS3Objects: async (bucketName: string): Promise<S3ObjectsResponse> => {
    return apiClient.get<S3ObjectsResponse>(`/api/v1/aws/s3/buckets/${bucketName}/objects`);
  },

  uploadS3Object: async (bucketName: string, file: File, key?: string): Promise<UploadObjectResponse> => {
    const formData = new FormData();
    formData.append('file', file);
    if (key) {
      formData.append('key', key);
    }

    const token = storage.get<string>(STORAGE_KEYS.ACCESS_TOKEN);

    const response = await fetch(`/api/v1/aws/s3/buckets/${bucketName}/objects`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token || ''}`,
      },
      body: formData,
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Upload failed: ${response.statusText} - ${errorText}`);
    }

    return response.json();
  },

  deleteS3Object: async (bucketName: string, key: string): Promise<DeleteObjectResponse> => {
    const encodedKey = encodeURIComponent(key);
    return apiClient.delete<DeleteObjectResponse>(
      `/api/v1/aws/s3/buckets/${bucketName}/objects/${encodedKey}`
    );
  },

  downloadS3Object: async (bucketName: string, key: string): Promise<void> => {
    const encodedKey = encodeURIComponent(key);
    const token = storage.get<string>(STORAGE_KEYS.ACCESS_TOKEN);

    const response = await fetch(`/api/v1/aws/s3/buckets/${bucketName}/download/${encodedKey}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token || ''}`,
      },
    });

    if (!response.ok) {
      throw new Error(`Download failed: ${response.statusText}`);
    }

    // Create a blob from the response
    const blob = await response.blob();

    // Create a temporary URL for the blob
    const url = window.URL.createObjectURL(blob);

    // Create a temporary anchor element and trigger download
    const a = document.createElement('a');
    a.href = url;
    a.download = key.split('/').pop() || 'download'; // Use filename from key
    document.body.appendChild(a);
    a.click();

    // Cleanup
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
  },

  // DynamoDB Operations
  listDynamoDBTables: async (): Promise<DynamoDBTablesResponse> => {
    return apiClient.get<DynamoDBTablesResponse>('/api/v1/aws/dynamodb/tables');
  },

  listDynamoDBRecords: async (): Promise<DynamoDBRecordsResponse> => {
    return apiClient.get<DynamoDBRecordsResponse>('/api/v1/aws/dynamodb/records');
  },

  upsertDynamoDBRecord: async (data: UpsertRecordRequest): Promise<UpsertRecordResponse> => {
    return apiClient.post<UpsertRecordResponse>('/api/v1/aws/dynamodb/tables', data);
  },
};
