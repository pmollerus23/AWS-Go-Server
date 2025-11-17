# S3 Bucket and File Management Features

Complete implementation of S3 bucket and file management functionality.

## Backend (Go)

### New Handlers in `internal/handlers/aws.go`

1. **HandleS3CreateBucket** - Create new S3 buckets
   - POST `/api/v1/aws/s3/buckets`
   - Supports region specification
   - Handles LocationConstraint for non-us-east-1 regions

2. **HandleS3DeleteBucket** - Delete S3 buckets
   - DELETE `/api/v1/aws/s3/buckets/{bucketName}`
   - Bucket must be empty before deletion

3. **HandleS3ListObjects** - List files in a bucket
   - GET `/api/v1/aws/s3/buckets/{bucketName}/objects`
   - Returns key, size, and lastModified for each object

4. **HandleS3UploadObject** - Upload files to S3
   - POST `/api/v1/aws/s3/buckets/{bucketName}/objects`
   - Accepts multipart/form-data
   - Supports custom key (filename) or uses uploaded filename
   - Max file size: 32MB

5. **HandleS3DeleteObject** - Delete files from S3
   - DELETE `/api/v1/aws/s3/buckets/{bucketName}/objects/{key}`
   - Handles URL-encoded keys with slashes

6. **HandleS3GetObject** - Download files from S3
   - GET `/api/v1/aws/s3/buckets/{bucketName}/objects/{key}/download`
   - Streams file to browser
   - Sets proper Content-Disposition headers

### Routes Added in `internal/server/routes.go`

All routes are protected with authentication middleware.

## Frontend (React/TypeScript)

### New TypeScript Types in `types/aws.types.ts`

- `CreateBucketRequest` / `CreateBucketResponse`
- `DeleteBucketResponse`
- `S3Object` / `S3ObjectsResponse`
- `UploadObjectResponse`
- `DeleteObjectResponse`

### New API Client Functions in `api/aws.api.ts`

- `createS3Bucket(data)` - Create bucket
- `deleteS3Bucket(bucketName)` - Delete bucket
- `listS3Objects(bucketName)` - List files
- `uploadS3Object(bucketName, file, key?)` - Upload file
- `deleteS3Object(bucketName, key)` - Delete file
- `downloadS3Object(bucketName, key)` - Get download URL

### New UI Components

#### 1. **S3BucketCreateForm** (`components/S3BucketCreateForm.tsx`)
- Form to create new S3 buckets
- Region selector with common AWS regions
- Success/error messaging

#### 2. **S3BucketManager** (`components/S3BucketManager.tsx`)
- Lists all S3 buckets
- Click to select bucket
- Delete bucket with confirmation
- Shows bucket creation date
- Visual indication of selected bucket

#### 3. **S3FileUploadForm** (`components/S3FileUploadForm.tsx`)
- File input with size display
- Optional custom key/filename
- Upload progress indication
- Success/error messaging

#### 4. **S3FilesList** (`components/S3FilesList.tsx`)
- Table view of files in selected bucket
- Shows file name, size (KB), last modified date
- Download button (opens in new tab)
- Delete button with confirmation
- Auto-refresh on upload

### Updated Pages

#### **AWSPage** (`pages/AWSPage.tsx`)
Reorganized into sections:

**S3 Section:**
1. Bucket creation form + Bucket manager (grid layout)
2. File upload form + Files list (shown only when bucket selected)

**DynamoDB Section:**
- Unchanged, still includes tables list and records management

## User Workflow

### Managing Buckets
1. Navigate to `/aws`
2. Enter bucket name and select region
3. Click "Create Bucket"
4. Bucket appears in the list
5. Click bucket to select it (highlights in blue)
6. Click "Delete" to remove bucket (requires confirmation)

### Managing Files
1. Select a bucket from the bucket list
2. File management section appears below
3. Choose a file to upload
4. Optionally specify a custom filename (key)
5. Click "Upload File"
6. File appears in the files table
7. Use "Download" to get the file
8. Use "Delete" to remove the file (requires confirmation)

## Styling

### New CSS Classes in `styles/globals.css`

- `.bucket-list` - Container for bucket items
- `.bucket-item` - Individual bucket with hover and selected states
- `.bucket-item.selected` - Blue border for selected bucket
- `.bucket-info` - Bucket name and metadata layout
- `.file-input-wrapper` - File input styling
- `.action-buttons` - Button group layout in tables

## Features

✅ Create S3 buckets with region selection
✅ Delete S3 buckets with confirmation
✅ List buckets with creation dates
✅ Select bucket to view files
✅ Upload files to selected bucket
✅ Download files from S3
✅ Delete files with confirmation
✅ Auto-refresh after operations
✅ Error handling and user feedback
✅ Responsive design
✅ Loading states

## Security

- All routes protected with JWT authentication
- File size limited to 32MB
- Confirmation dialogs for destructive operations
- Authorization header on all requests
- URL encoding for file keys with special characters
