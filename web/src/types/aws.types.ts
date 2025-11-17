export interface S3Bucket {
  name: string;
  creationDate: string;
}

export interface S3BucketsResponse {
  buckets: S3Bucket[];
  count: number;
}

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
