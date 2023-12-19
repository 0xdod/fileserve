# File Upload Service

This is a Golang application serving as a File Upload Service. It handles HTTP requests for file uploads and downloads. It uses AWS S3 for file storage and sqlite3 as a database 

## Requirements

- Go 1.21 or higher
- Goose SQL database migration tool (https://github.com/pressly/goose)
- GNU Make (https://www.gnu.org/software/make/) 
- An AWS account, user credentials and an S3 Bucket that is configured for public read.

## Installation

1. Clone the repository: `git clone https://github.com/0xdod/fileserve.git`
2. Change to the project directory: `cd fileserve`
3. Install dependencies: `go mod download`

## Usage

1. Run migrations: `make up-migrate`
2. Set the environment variables: `cp .env.sample .env`c, setting valid S3 access and secret keys
3. Start the server: `make run`
4. Open your browser and navigate to `http://localhost:7000` to access the application.

## Endpoints

### Upload a File

- **URL**: `/api/v1/files/upload`
- **Method**: POST
- **Request Body**: Multipart Form data with a file field named `file`
- **Response**: JSON object with the uploaded file details

### Get List of Uploaded Files

- **URL**: `/api/v1/files`
- **Method**: GET
- **Response**: JSON array containing the list of uploaded files

### Download a File

- **URL**: `/download/{fileId}`
- **Method**: GET
- **Response**: The file will be downloaded

## Testing

To run the unit tests, use the following command:
```
go test ./...
```
