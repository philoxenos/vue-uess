# Management Information System (MIS)

A Management Information System built with Vue.js frontend and Go backend that serves as the single source of truth for an Android application.

## Project Structure

The project is divided into two main parts:

1. **Backend**: A Go API server with user management functionality
2. **Frontend**: A Vue.js admin interface for managing users

## Features

- User authentication (login/logout)
- User management (create, read, update, delete)
- Integration with Android app via Google authentication
- Secure password handling with bcrypt
- JWT token-based authentication

## Setup Instructions

### Prerequisites

- Go 1.16+ 
- Node.js 14+
- npm or yarn

### Backend Setup

1. Navigate to the backend directory:
   ```
   cd backend
   ```

2. Install Go dependencies:
   ```
   go get
   ```

3. Run the backend server:
   ```
   go run main.go
   ```

   The API server will run on http://localhost:8080

### Frontend Setup

1. Navigate to the frontend directory:
   ```
   cd frontend
   ```

2. Install dependencies:
   ```
   npm install
   ```

3. Start the development server:
   ```
   npm run dev
   ```

   The admin interface will be available at http://localhost:5173

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login a user and get JWT token

### Users

- `GET /api/v1/users` - Get all users (requires authentication)
- `GET /api/v1/users/:id` - Get a specific user (requires authentication)
- `PUT /api/v1/users/:id` - Update a user (requires authentication)
- `DELETE /api/v1/users/:id` - Delete a user (requires authentication)

## Android Integration

This MIS serves as the backend for an Android application, where users:
1. Select a Google/Gmail account in the mobile app
2. Create an account by setting a password in this system
3. The account credentials are stored locally in the Android app using Room Database with SQLCipher
4. Users can then authenticate with the system using their credentials

## Security Features

- Password hashing with bcrypt
- JWT token authentication
- Route protection middleware
- Role-based access control (admin vs regular users)

## Production Deployment

For a production deployment, consider the following:

1. Set proper environment variables for sensitive values
2. Use HTTPS with proper SSL certificates
3. Implement proper logging and monitoring
4. Set up a production-ready database (PostgreSQL or MySQL)
5. Configure proper CORS settings
