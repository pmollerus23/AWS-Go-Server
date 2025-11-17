import { apiClient } from './client';
import type {
  S3BucketsResponse,
  DynamoDBTablesResponse,
  DynamoDBRecordsResponse,
  UpsertRecordRequest,
  UpsertRecordResponse,
} from '../types/aws.types';

export const awsApi = {
  // S3 Operations
  listS3Buckets: async (): Promise<S3BucketsResponse> => {
    return apiClient.get<S3BucketsResponse>('/api/v1/aws/s3/buckets');
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
