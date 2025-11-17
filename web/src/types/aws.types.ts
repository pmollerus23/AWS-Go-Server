// S3 Bucket Types
export interface S3Bucket {
  name: string;
  creationDate: string;
}

export interface S3BucketsResponse {
  buckets: S3Bucket[];
  count: number;
}

export interface CreateBucketRequest {
  bucketName: string;
  region?: string;
}

export interface CreateBucketResponse {
  success: boolean;
  bucketName: string;
}

export interface DeleteBucketResponse {
  success: boolean;
  bucketName: string;
}

// S3 Object Types
export interface S3Object {
  key: string;
  size: number;
  lastModified: string;
}

export interface S3ObjectsResponse {
  objects: S3Object[];
  count: number;
}

export interface UploadObjectResponse {
  success: boolean;
  key: string;
  bucket: string;
}

export interface DeleteObjectResponse {
  success: boolean;
  key: string;
  bucket: string;
}

// DynamoDB Types
export interface DynamoDBTable {
  name: string;
}

export interface DynamoDBTablesResponse {
  tables: string[];
  count: number;
}

export interface DynamoDBRecord {
  id: number;
  name: string;
  updated_at: number;
}

export interface DynamoDBRecordsResponse {
  records: DynamoDBRecord[];
  count: number;
}

export interface UpsertRecordRequest {
  id: number;
  name: string;
  updated_at: number;
}

export interface UpsertRecordResponse {
  success: boolean;
  result_attributes?: Record<string, any>;
}
