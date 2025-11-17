import { useState } from 'react';
import {
  DynamoDBRecordForm,
  DynamoDBRecordsList,
  DynamoDBTablesList,
  S3BucketCreateForm,
  S3BucketManager,
  S3FileUploadForm,
  S3FilesList,
} from '../components';

export const AWSPage: React.FC = () => {
  const [dynamoRefreshTrigger, setDynamoRefreshTrigger] = useState(0);
  const [s3RefreshTrigger, setS3RefreshTrigger] = useState(0);
  const [selectedBucket, setSelectedBucket] = useState<string | null>(null);

  const handleRecordUpserted = () => {
    setDynamoRefreshTrigger(prev => prev + 1);
  };

  const handleBucketCreated = () => {
    setS3RefreshTrigger(prev => prev + 1);
  };

  const handleFileUploaded = () => {
    setS3RefreshTrigger(prev => prev + 1);
  };

  return (
    <div className="aws-page">
      <h1>AWS Resources</h1>
      <p>Manage your AWS resources including S3 buckets and DynamoDB tables.</p>

      <section className="aws-section">
        <h2>Amazon S3</h2>

        {/* Bucket Management */}
        <div className="page-grid">
          <div className="grid-item">
            <S3BucketCreateForm onSuccess={handleBucketCreated} />
          </div>
          <div className="grid-item grid-span-2">
            <S3BucketManager
              onSelectBucket={setSelectedBucket}
              refreshTrigger={s3RefreshTrigger}
            />
          </div>
        </div>

        {/* File Management - Only show when bucket is selected */}
        {selectedBucket && (
          <>
            <h3 style={{ marginTop: '2rem', marginBottom: '1rem' }}>
              File Management for {selectedBucket}
            </h3>
            <div className="page-grid">
              <div className="grid-item">
                <S3FileUploadForm
                  bucketName={selectedBucket}
                  onSuccess={handleFileUploaded}
                />
              </div>
              <div className="grid-item grid-span-2">
                <S3FilesList
                  bucketName={selectedBucket}
                  refreshTrigger={s3RefreshTrigger}
                />
              </div>
            </div>
          </>
        )}
      </section>

      <section className="aws-section">
        <h2>Amazon DynamoDB</h2>
        <div className="page-grid">
          <div className="grid-item grid-span-full">
            <DynamoDBTablesList />
          </div>
          <div className="grid-item">
            <DynamoDBRecordForm onSuccess={handleRecordUpserted} />
          </div>
          <div className="grid-item grid-span-2">
            <DynamoDBRecordsList refreshTrigger={dynamoRefreshTrigger} />
          </div>
        </div>
      </section>
    </div>
  );
};
