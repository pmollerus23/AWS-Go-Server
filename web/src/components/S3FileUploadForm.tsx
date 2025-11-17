import { useState, useRef } from 'react';
import { Button, Card } from '../ui';
import { useMutation } from '../hooks';
import { awsApi } from '../api';

interface S3FileUploadFormProps {
  bucketName: string;
  onSuccess?: () => void;
}

export const S3FileUploadForm: React.FC<S3FileUploadFormProps> = ({
  bucketName,
  onSuccess
}) => {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [customKey, setCustomKey] = useState('');
  const fileInputRef = useRef<HTMLInputElement>(null);

  const { mutate, isLoading, isError, error, isSuccess } = useMutation(
    async (file: File) => awsApi.uploadS3Object(bucketName, file, customKey || undefined),
    {
      onSuccess: () => {
        setSelectedFile(null);
        setCustomKey('');
        if (fileInputRef.current) {
          fileInputRef.current.value = '';
        }
        onSuccess?.();
      },
    }
  );

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setSelectedFile(file);
      if (!customKey) {
        setCustomKey(file.name);
      }
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (selectedFile) {
      mutate(selectedFile);
    }
  };

  return (
    <Card>
      <h3>Upload File to {bucketName}</h3>
      <form onSubmit={handleSubmit} className="form">
        <div className="file-input-wrapper">
          <label className="input-label">Select File</label>
          <input
            ref={fileInputRef}
            type="file"
            onChange={handleFileChange}
            className="input"
            required
          />
          {selectedFile && (
            <small>Selected: {selectedFile.name} ({(selectedFile.size / 1024).toFixed(2)} KB)</small>
          )}
        </div>

        <div className="input-wrapper">
          <label className="input-label">
            Key (filename in S3)
            <small> - Leave empty to use original filename</small>
          </label>
          <input
            type="text"
            value={customKey}
            onChange={(e) => setCustomKey(e.target.value)}
            className="input"
            placeholder={selectedFile?.name || 'filename.ext'}
          />
        </div>

        {isError && (
          <div className="error-message">
            Error: {error?.message || 'Failed to upload file'}
          </div>
        )}

        {isSuccess && (
          <div className="success-message">
            File uploaded successfully!
          </div>
        )}

        <Button type="submit" disabled={isLoading || !selectedFile}>
          {isLoading ? 'Uploading...' : 'Upload File'}
        </Button>
      </form>
    </Card>
  );
};
