import { useState } from 'react';
import {
  DynamoDBRecordForm,
  DynamoDBRecordsList,
  DynamoDBTablesList,
  S3BucketsList,
} from '../components';

export const AWSPage: React.FC = () => {
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const handleRecordUpserted = () => {
    // Trigger refresh of records list
    setRefreshTrigger(prev => prev + 1);
  };

  return (
    <div className="aws-page">
      <h1>AWS Resources</h1>
      <p>Manage your AWS resources including S3 buckets and DynamoDB tables.</p>

      <section className="aws-section">
        <h2>Amazon S3</h2>
        <div className="page-grid">
          <div className="grid-item grid-span-full">
            <S3BucketsList />
          </div>
        </div>
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
            <DynamoDBRecordsList refreshTrigger={refreshTrigger} />
          </div>
        </div>
      </section>
    </div>
  );
};
