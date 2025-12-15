# Go CRUD API with MongoDB

A simple CRUD (Create, Read, Update, Delete) REST API built with Go and MongoDB.

## Features

- RESTful API endpoints for user management
- MongoDB integration using the official Go driver
- JSON request/response handling
- Proper error handling and status codes
- Graceful server shutdown
- Health check endpoint

## Prerequisites

- Go 1.21 or higher
- MongoDB running locally or accessible via connection string

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd go-backend-project
```

2. Install dependencies:
```bash
go mod tidy
```

3. Make sure MongoDB is running. By default, the application connects to `mongodb://localhost:27017`.

4. Run the application:

### Option A: Manual run
```bash
go run main.go
```

### Option B: With file watcher (recommended for development)
```bash
# Using Air (recommended)
air

# Or use the development script
./dev.sh

# Or manually choose your watcher
./watch.sh
```

The server will start on port 8081 by default. You can change the port by setting the `PORT` environment variable.

## API Endpoints

### Health Check
- **GET** `/health` - Check if the server is running

### Users

#### Create User
- **POST** `/api/users`
- **Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30
}
```

#### Get All Users
- **GET** `/api/users`

#### Get User by ID
- **GET** `/api/users/{id}`

#### Update User
- **PUT** `/api/users/{id}`
- **Body:**
```json
{
  "name": "John Updated",
  "email": "john.updated@example.com",
  "age": 31
}
```

#### Delete User
- **DELETE** `/api/users/{id}`

## Response Format

All API responses follow this format:

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {
    // Response data here
  }
}
```

Error responses:
```json
{
  "success": false,
  "error": "Error message here"
}
```

## Project Structure

```
go-backend-project/
├── main.go              # Application entry point
├── go.mod               # Go module file
├── README.md            # This file
├── models/
│   └── user.go          # User model and request structs
├── database/
│   └── mongodb.go       # MongoDB connection and configuration
├── handlers/
│   └── user_handler.go  # HTTP handlers for user operations
└── routes/
    └── routes.go        # Route configuration
```

## Configuration

### MongoDB Connection

The application connects to MongoDB using the connection string in `database/mongodb.go`. You can modify the connection string to point to your MongoDB instance:

```go
clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
```

### Environment Variables

- `PORT` - Server port (default: 8080)

## Testing the API

You can test the API using curl or any HTTP client like Postman:

### Create a user:
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "age": 30}'
```

### Get all users:
```bash
curl http://localhost:8080/api/users
```

### Get a specific user:
```bash
curl http://localhost:8080/api/users/{user_id}
```

### Update a user:
```bash
curl -X PUT http://localhost:8080/api/users/{user_id} \
  -H "Content-Type: application/json" \
  -d '{"name": "John Updated", "email": "john.updated@example.com", "age": 31}'
```

### Delete a user:
```bash
curl -X DELETE http://localhost:8080/api/users/{user_id}
```

## Development Tools

### File Watchers
- **Air** (recommended) - Go-specific live reload tool
- **Reflex** - Simple file watcher
- **Nodemon** - Node.js file watcher (if you have Node.js installed)
- **Custom shell script** - Basic file watching with shell commands

### Installation
```bash
# Install Air (recommended)
go install github.com/air-verse/air@latest

# Install Reflex (alternative)
go install github.com/cespare/reflex@latest

# Install Nodemon (if you have Node.js)
npm install -g nodemon
```

## Dependencies

- `github.com/gorilla/mux` - HTTP router and URL matcher
- `go.mongodb.org/mongo-driver` - Official MongoDB Go driver

## License

This project is open source and available under the [MIT License](LICENSE).
