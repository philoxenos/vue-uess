# UESS - User Enrollment and Sign-in System

A complete web application with Vue.js frontend and Go backend that implements a secure authentication system matching the UESS Android application.

## Project Structure

The project is divided into two main parts:

1. **Backend**: A Go API server with user management and authentication functionality
2. **Frontend**: A Vue.js admin interface with Google Sign-In integration

## Features

- Google OAuth integration for seamless authentication
- Secure JWT-based authentication with access and refresh tokens
- Role-based access control
- User management (create, read, update, delete)
- Secure password handling with bcrypt
- Password reset functionality
- Session management
- Audit logging for security events

## Authentication Flow

### First-time Login
1. User selects "Continue with Google"
2. Google authentication takes place
3. Backend checks if user exists in database:
   - If exists, issues access & refresh tokens
   - If not, auto-provisions account tied to Google identity
   - If local password is required, redirects to registration page with email pre-filled

### Returning Login
1. On app launch, checks for valid tokens
2. If token is expired but refresh token exists, silently refreshes
3. If refresh fails or no tokens exist, shows login screen

### Forgot Password
1. User enters email address
2. If account exists with local password, sends reset instructions
3. If account uses Google-only auth, notifies user to use Google Sign-In

## Technical Architecture

### Backend (Go)

- **Framework**: Gin web framework
- **Database**: SQLite with GORM ORM
- **Authentication**: JWT tokens (access + refresh)
- **Google Auth**: OAuth2 integration with Google Identity Services

### Frontend (Vue.js)

- **Framework**: Vue 3 with Composition API
- **State Management**: Pinia
- **UI Components**: Vuetify
- **Routing**: Vue Router with auth guards
- **HTTP Client**: Axios with interceptors for token refresh

## Setup Instructions

### Prerequisites

- Go 1.24+
- Node.js 16+
- npm 7+
- GCC compiler (for SQLite with CGO) or use the pure Go SQLite driver

### Backend Setup

1. Navigate to the backend directory:
   ```
   cd backend
   ```

2. Install Go dependencies:
   ```
   go mod tidy
   ```

3. Run the backend server:
   ```
   $env:CGO_ENABLED=1; go run main.go
   ```
   
   Or if using pure Go SQLite driver (no CGO required):
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

### Easy Start (Both Servers)

Use the provided PowerShell script to start both servers at once:

```
./start-dev.ps1
```

## Configuration

1. Backend configuration:
   - Update Google OAuth credentials in `handlers/google_auth.go`
   - Change JWT secret key in `handlers/auth_handlers.go`

2. Frontend configuration:
   - Update Google Client ID in `src/views/Login.vue`
   - Adjust API base URL in `src/services/api.js` if needed

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login with email and password
- `POST /api/v1/auth/google` - Authenticate with Google ID token
- `GET /api/v1/auth/google/login` - Initiate Google OAuth flow
- `GET /api/v1/auth/google/callback` - Handle Google OAuth callback
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout (revoke refresh token)
- `POST /api/v1/auth/forgot-password` - Request password reset
- `POST /api/v1/auth/reset-password` - Reset password with token

### User Management
- `GET /api/v1/users` - Get all users (requires authentication)
- `GET /api/v1/users/:id` - Get a specific user (requires authentication)
- `PUT /api/v1/users/:id` - Update a user (requires authentication)
- `DELETE /api/v1/users/:id` - Delete a user (requires authentication)
- `GET /api/v1/me` - Get current user info (requires authentication)

## Android Integration

This system serves as the backend for the UESS Android application, implementing the following flows:

### First Time Login Flow:
- App → Google Login → Check Local Database → If not exist → Open browser (MIS) → Enter Password → Account is Created → Redirect to App Automatically

### Next Time Login Flow:
- App → Google Login → Check Local Database → If exist → Redirect to App Automatically

### Forgot Password Flow:
- App → Forgot Password → Select Email → If exist → Open browser (MIS) → Enter Password → Account is Updated → Redirect to App Automatically

## Security Features

- Access tokens are short-lived (15 minutes)
- Refresh tokens are stored securely and rotated on use
- All password hashes use bcrypt
- Comprehensive audit logging for security events
- CORS properly configured
- Role-based access control enforced on both client and server
- Pure SQLite Go driver or CGO-enabled SQLite driver options

## Production Deployment

For a production deployment, consider the following:

1. Set proper environment variables for sensitive values
2. Use HTTPS with proper SSL certificates
3. Implement proper logging and monitoring
4. Set up a production-ready database (PostgreSQL or MySQL)
5. Configure stricter CORS settings
6. Use a reverse proxy (Nginx, Caddy) for serving the application
7. Implement rate limiting for auth endpoints
8. Set up monitoring for failed login attempts
