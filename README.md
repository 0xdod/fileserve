# File Upload Service

This is a Golang application serving as a File Upload Service. It handles HTTP requests for file uploads and downloads.

## Requirements

- Go 1.21 or higher

## Installation

1. Clone the repository: `git clone https://github.com/your-username/file-upload-service.git`
2. Change to the project directory: `cd file-upload-service`
3. Install dependencies: `go mod download`

## Usage

1. Start the server: `go run main.go`
2. Open your browser and navigate to `http://localhost:8080` to access the application.

## Endpoints

### Upload a File

- **URL**: `/upload`
- **Method**: POST
- **Request Body**: Form data with a file field named `file`
- **Response**: JSON object with the uploaded file details

### Get List of Uploaded Files

- **URL**: `/files`
- **Method**: GET
- **Response**: JSON array containing the list of uploaded files

### Download a File

- **URL**: `/download/{filename}`
- **Method**: GET
- **Response**: The file will be downloaded

## Testing

To run the unit tests, use the following command:
