# Hope - Carpooling Platform

## Table of Contents
1. [Project Overview](#project-overview)
2. [Architecture](#architecture)
3. [Features](#features)
4. [Technology Stack](#technology-stack)
5. [Project Structure](#project-structure)
6. [API Documentation](#api-documentation)
7. [Database Schema](#database-schema)
8. [Authentication & Security](#authentication--security)
9. [Installation & Setup](#installation--setup)
10. [Development](#development)
11. [Deployment](#deployment)
12. [Contributing](#contributing)

## Project Overview

**Hope** is a comprehensive carpooling platform designed specifically for students and university communities. The platform facilitates ride-sharing by connecting drivers offering rides with passengers seeking transportation, promoting sustainable mobility and community building.

### Key Objectives
- **Student-Focused**: Restricted to authorized educational domain emails
- **Location-Based Matching**: Uses geohash-based proximity search for optimal ride matching
- **Real-Time Communication**: Built-in chat system for rider-driver coordination
- **Trust & Safety**: Review system and user verification through Google OAuth
- **Scalable Architecture**: Microservices-based design with gRPC communication

### Target Users
- **Students**: Seeking affordable and convenient transportation
- **Drivers**: Students offering rides to offset travel costs
- **Universities**: Promoting sustainable transportation initiatives

## Architecture

### System Design
The platform follows a **Clean Architecture** pattern with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   gRPC Client  │    │   gRPC Server   │    │   MySQL DB      │
│   (Frontend)   │◄──►│   (Backend)     │◄──►│   (Data Layer)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   Google OAuth  │
                       │   (Auth Layer)  │
                       └─────────────────┘
```

### Architecture Layers

#### 1. **Transport Layer (gRPC)**
- **Protocol**: gRPC with Protocol Buffers
- **Communication**: Unary RPC calls
- **Interceptors**: JWT authentication middleware
- **Reflection**: Enabled for development and testing

#### 2. **Application Layer (Handlers)**
- **API Handlers**: Convert gRPC requests to domain operations
- **Validation**: Input sanitization and business rule enforcement
- **Error Handling**: Consistent error responses across services

#### 3. **Domain Layer (Services)**
- **Business Logic**: Core carpooling algorithms and rules
- **Transaction Management**: Ride matching and status updates
- **Domain Validation**: Business rule enforcement

#### 4. **Infrastructure Layer (Repositories)**
- **Data Access**: GORM-based database operations
- **Query Optimization**: Geohash-based spatial queries
- **Connection Management**: Database connection pooling

#### 5. **Data Layer (Database)**
- **Storage**: MySQL with GORM ORM
- **Spatial Indexing**: Geohash-based location queries
- **Auto-migration**: Schema management and versioning

### Design Patterns

#### **Dependency Injection (Wire)**
- **Purpose**: Loose coupling and testability
- **Implementation**: Google Wire for compile-time DI
- **Benefits**: Reduced coupling, easier testing, maintainable code

#### **Repository Pattern**
- **Purpose**: Abstract data access layer
- **Implementation**: Interface-based repository contracts
- **Benefits**: Testable, swappable data sources

#### **Service Layer Pattern**
- **Purpose**: Encapsulate business logic
- **Implementation**: Domain-specific service interfaces
- **Benefits**: Reusable business logic, clear separation

## Features

### Core Functionality

#### 1. **User Management**
- **Google OAuth Integration**: Secure authentication via Google accounts
- **Domain Restriction**: Limited to authorized educational institutions
- **Profile Management**: User profiles with photos and preferences
- **Last Seen Tracking**: User activity monitoring

#### 2. **Ride Management**
- **Ride Offers**: Drivers can create ride offers with details
- **Ride Requests**: Passengers can request rides to destinations
- **Dynamic Pricing**: Fare negotiation and seat management
- **Status Tracking**: Active, matched, completed ride states

#### 3. **Location Services**
- **Geohash Encoding**: Efficient spatial indexing and queries
- **Proximity Search**: Find nearby rides and passengers
- **Real-time Updates**: Location tracking and updates
- **Spatial Optimization**: Route-based matching algorithms

#### 4. **Matching System**
- **Request Management**: Rider join requests for specific rides
- **Driver Approval**: Accept/reject passenger requests
- **Status Workflow**: Requested → Accepted → Completed
- **Conflict Resolution**: Handle multiple requests per ride

#### 5. **Communication**
- **Ride-Specific Chat**: Private messaging between matched users
- **Message Persistence**: Chat history and conversation threads
- **Access Control**: Only matched users can communicate
- **Real-time Messaging**: Instant communication for coordination

#### 6. **Review System**
- **Bidirectional Reviews**: Both riders and drivers can rate each other
- **Rating Scale**: 1-5 star rating system
- **Comment System**: Detailed feedback and suggestions
- **Quality Assurance**: Maintain platform trust and safety

### Advanced Features

#### **Geospatial Intelligence**
- **Geohash Implementation**: Efficient location-based queries
- **Proximity Algorithms**: Smart ride matching based on distance
- **Route Optimization**: Optimal pickup and drop-off coordination

#### **Security & Privacy**
- **JWT Authentication**: Secure session management
- **Domain Verification**: Educational institution email validation
- **Data Encryption**: Sensitive information protection
- **Access Control**: Role-based permissions and restrictions

## Technology Stack

### Backend Technologies

#### **Core Framework**
- **Language**: Go 1.24.4
- **Runtime**: Compiled binary with high performance
- **Concurrency**: Goroutines and channels for scalability

#### **Communication Protocol**
- **gRPC**: High-performance RPC framework
- **Protocol Buffers**: Efficient serialization format
- **HTTP/2**: Modern transport protocol

#### **Database & ORM**
- **Database**: MySQL 8.0+
- **ORM**: GORM v1.30.1
- **Driver**: MySQL driver with connection pooling
- **Migration**: Auto-migration with schema versioning

#### **Authentication & Security**
- **OAuth Provider**: Google OAuth 2.0
- **JWT Library**: golang-jwt/jwt v3.2.2
- **Token Validation**: Google token info endpoint
- **Secret Management**: Environment-based configuration

#### **Dependency Management**
- **Module System**: Go modules (go.mod)
- **Dependency Injection**: Google Wire v0.6.0
- **Package Management**: go.sum for dependency verification

### Development Tools

#### **Code Quality**
- **Linting**: Go vet and static analysis
- **Formatting**: gofmt for consistent code style
- **Testing**: Built-in Go testing framework
- **Documentation**: Comprehensive code comments

#### **Configuration Management**
- **Environment Variables**: .env file support via godotenv
- **Configuration Structs**: Type-safe configuration management
- **Validation**: Runtime configuration validation

## Project Structure

```
hope/
├── api/                    # gRPC API handlers
│   ├── auth_handlers.go   # Authentication endpoints
│   ├── chat_handler.go    # Chat service endpoints
│   ├── location_handler.go # Location service endpoints
│   ├── match_handler.go   # Matching service endpoints
│   ├── review_handler.go  # Review service endpoints
│   ├── ride_handler.go    # Ride service endpoints
│   └── user_handler.go    # User service endpoints
├── cmd/                   # Application entry points
├── config/                # Configuration management
│   ├── config.go         # Domain and JWT configuration
│   └── connection.go     # Database connection setup
├── db/                    # Database models and schemas
│   ├── user.go           # User entity model
│   ├── ride_offer.go     # Ride offer model
│   ├── ride_request.go   # Ride request model
│   ├── match.go          # Match entity model
│   ├── chat_message.go   # Chat message model
│   ├── review.go         # Review entity model
│   └── user_location.go  # User location model
├── di/                    # Dependency injection
│   ├── wire.go           # Wire dependency definitions
│   └── wire_gen.go       # Generated dependency code
├── middleware/            # HTTP/gRPC middleware
│   └── auth.go           # JWT authentication interceptor
├── proto/                 # Protocol Buffer definitions
│   └── v1/               # API version 1
│       ├── auth.proto    # Authentication service
│       ├── user.proto    # User management service
│       ├── ride.proto    # Ride management service
│       ├── match.proto   # Matching service
│       ├── chat.proto    # Chat service
│       ├── review.proto  # Review service
│       └── location.proto # Location service
├── repository/            # Data access layer
│   ├── user_repository.go        # User data operations
│   ├── ride_offer_repository.go  # Ride offer operations
│   ├── ride_request_repository.go # Ride request operations
│   ├── match_repository.go       # Match operations
│   ├── chat_message_repository.go # Chat operations
│   ├── review_repository.go      # Review operations
│   └── user_location_repository.go # Location operations
├── service/               # Business logic layer
│   ├── auth_service.go   # Authentication business logic
│   ├── user_service.go   # User management logic
│   ├── ride_service.go   # Ride management logic
│   ├── match_service.go  # Matching algorithms
│   ├── chat_service.go   # Chat business logic
│   ├── review_service.go # Review management
│   └── location_service.go # Location services
├── utils/                 # Utility functions and helpers
├── main.go               # Application entry point
├── go.mod                # Go module dependencies
├── go.sum                # Dependency checksums
├── Class_Diagram.pdf     # System architecture diagram
└── SRS.pdf              # Software Requirements Specification
```

### Directory Responsibilities

#### **`api/` - Transport Layer**
- **Purpose**: Handle gRPC requests and responses
- **Responsibilities**: Request validation, response formatting, error handling
- **Pattern**: Handler interfaces implementing gRPC service contracts

#### **`config/` - Configuration Management**
- **Purpose**: Centralized configuration and environment management
- **Responsibilities**: Database config, JWT secrets, domain restrictions
- **Pattern**: Environment-based configuration with validation

#### **`db/` - Data Models**
- **Purpose**: Define database schema and entity relationships
- **Responsibilities**: GORM tags, field validation, relationship mapping
- **Pattern**: Domain-driven design with clear entity boundaries

#### **`di/` - Dependency Injection**
- **Purpose**: Wire application dependencies at compile time
- **Responsibilities**: Service composition, repository wiring, handler creation
- **Pattern**: Constructor injection with Google Wire

#### **`middleware/` - Cross-cutting Concerns**
- **Purpose**: Handle authentication and request processing
- **Responsibilities**: JWT validation, user context injection, access control
- **Pattern**: Interceptor pattern for gRPC middleware

#### **`proto/` - API Contracts**
- **Purpose**: Define service interfaces and message formats
- **Responsibilities**: RPC definitions, message schemas, service contracts
- **Pattern**: Contract-first API design with Protocol Buffers

#### **`repository/` - Data Access**
- **Purpose**: Abstract database operations and queries
- **Responsibilities**: CRUD operations, query optimization, data persistence
- **Pattern**: Repository pattern with interface contracts

#### **`service/` - Business Logic**
- **Purpose**: Implement core business rules and workflows
- **Responsibilities**: Ride matching, validation, business processes
- **Pattern**: Service layer with domain-specific interfaces

## API Documentation

### Service Overview

The platform exposes 7 core gRPC services, each handling specific domain functionality:

#### **1. AuthService** - Authentication & Authorization
```protobuf
service AuthService {
    rpc Login(LoginRequest) returns (LoginResponse);
}
```

**Purpose**: Handle Google OAuth authentication and JWT token generation
**Features**: Domain validation, user provisioning, secure token issuance

#### **2. UserService** - User Profile Management
```protobuf
service UserService {
    rpc GetMe(GetMeRequest) returns (GetMeResponse);
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
    rpc UpdateMe(UpdateMeRequest) returns (UpdateMeResponse);
    rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}
```

**Purpose**: Manage user profiles, updates, and user discovery
**Features**: Partial updates with FieldMask, user listing, profile management

#### **3. RideService** - Ride Management
```protobuf
service RideService {
    // Ride Offers
    rpc CreateOffer(CreateOfferRequest) returns (CreateOfferResponse);
    rpc GetOffer(GetOfferRequest) returns (GetOfferResponse);
    rpc UpdateOffer(UpdateOfferRequest) returns (UpdateOfferResponse);
    rpc DeleteOffer(DeleteOfferRequest) returns (DeleteOfferResponse);
    rpc ListNearbyOffers(ListNearbyOffersRequest) returns (ListNearbyOffersResponse);
    
    // Ride Requests
    rpc CreateRequest(CreateRequestRequest) returns (CreateRequestResponse);
    rpc GetRequest(GetRequestRequest) returns (GetRequestResponse);
    rpc UpdateRequestStatus(UpdateRequestStatusRequest) returns (UpdateRequestStatusResponse);
    rpc DeleteRequest(DeleteRequestRequest) returns (DeleteRequestResponse);
    rpc ListNearbyRequests(ListNearbyRequestsRequest) returns (ListNearbyRequestsResponse);
}
```

**Purpose**: Manage ride offers and requests with location-based matching
**Features**: Geohash proximity search, status management, CRUD operations

#### **4. MatchService** - Ride Matching
```protobuf
service MatchService {
    rpc RequestToJoin(RequestToJoinRequest) returns (RequestToJoinResponse);
    rpc AcceptRequest(AcceptRequestRequest) returns (AcceptRequestResponse);
    rpc RejectRequest(RejectRequestRequest) returns (RejectRequestResponse);
    rpc CompleteMatch(CompleteMatchRequest) returns (CompleteMatchResponse);
    rpc GetMatch(GetMatchRequest) returns (GetMatchResponse);
    rpc ListMatchesByRide(ListMatchesByRideRequest) returns (ListMatchesByRideResponse);
    rpc ListMatchesByRider(ListMatchesByRiderRequest) returns (ListMatchesByRiderResponse);
}
```

**Purpose**: Handle ride matching workflow and status management
**Features**: Request lifecycle, driver approval, match completion

#### **5. ChatService** - Communication
```protobuf
service ChatService {
    rpc SendMessage(SendMessageRequest) returns (SendMessageResponse);
    rpc ListMessagesByRide(ListMessagesByRideRequest) returns (ListMessagesByRideResponse);
    rpc ListMessagesBySender(ListMessagesBySenderRequest) returns (ListMessagesBySenderResponse);
    rpc ListChatsForUser(ListChatsForUserRequest) returns (ListChatsForUserResponse);
    rpc DeleteMessage(DeleteMessageRequest) returns (DeleteMessageResponse);
}
```

**Purpose**: Enable communication between matched riders and drivers
**Features**: Ride-specific chats, message persistence, access control

#### **6. ReviewService** - Rating & Feedback
```protobuf
service ReviewService {
    rpc SubmitReview(SubmitReviewRequest) returns (SubmitReviewResponse);
    rpc ListReviewsByUser(ListReviewsByUserRequest) returns (ListReviewsByUserResponse);
    rpc ListReviewsByRide(ListReviewsByRideRequest) returns (ListReviewsByRideResponse);
    rpc DeleteReview(DeleteReviewRequest) returns (DeleteReviewResponse);
}
```

**Purpose**: Manage user reviews and ratings for completed rides
**Features**: Bidirectional reviews, rating system, feedback management

#### **7. LocationService** - Location Management
```protobuf
service LocationService {
    rpc UpsertLocation(UpsertLocationRequest) returns (UpsertLocationResponse);
    rpc GetLocationByUser(GetLocationByUserRequest) returns (GetLocationByUserResponse);
    rpc ListNearby(ListNearbyRequest) returns (ListNearbyResponse);
    rpc DeleteLocation(DeleteLocationRequest) returns (DeleteLocationResponse);
}
```

**Purpose**: Handle user location updates and proximity queries
**Features**: Real-time location tracking, geohash encoding, spatial queries

### Authentication Flow

#### **Login Process**
1. **Client Request**: Send Google ID token to `/proto.v1.AuthService/Login`
2. **Token Validation**: Verify token with Google OAuth endpoint
3. **Domain Verification**: Check email domain against allowed list
4. **User Provisioning**: Create or retrieve user profile
5. **JWT Generation**: Issue backend JWT with user claims
6. **Response**: Return JWT token and user information

#### **Request Authentication**
1. **Token Extraction**: Extract Bearer token from Authorization header
2. **JWT Validation**: Verify signature and expiration
3. **Claims Extraction**: Extract user ID and email from token
4. **Context Injection**: Inject user identity into gRPC context
5. **Request Processing**: Handle authenticated request

### Error Handling

#### **gRPC Status Codes**
- **OK**: Successful operation
- **InvalidArgument**: Invalid request parameters
- **Unauthenticated**: Missing or invalid authentication
- **NotFound**: Requested resource not found
- **Internal**: Server-side errors

#### **Error Response Format**
```protobuf
message ErrorResponse {
    string error_code = 1;
    string message = 2;
    string details = 3;
}
```

## Database Schema

### Entity Relationships

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐
│    User     │    │ RideOffer    │    │ RideRequest │
│             │    │              │    │             │
│ - ID (PK)   │◄───┤ - DriverID   │    │ - UserID    │
│ - Name      │    │ - FromGeo    │    │ - FromGeo   │
│ - Email     │    │ - ToGeo      │    │ - ToGeo     │
│ - PhotoURL  │    │ - Fare       │    │ - Time      │
│ - Geohash   │    │ - Time       │    │ - Seats     │
│ - LastSeen  │    │ - Seats      │    │ - Status    │
└─────────────┘    │ - Status     │    └─────────────┘
                   └──────────────┘
                          │
                          ▼
                   ┌──────────────┐
                   │    Match     │
                   │              │
                   │ - ID (PK)    │
                   │ - RiderID    │
                   │ - DriverID   │
                   │ - RideId     │
                   │ - Status     │
                   │ - CreatedAt  │
                   └──────────────┘
                          │
                          ▼
              ┌─────────────────────────┐
              │     ChatMessage         │
              │                         │
              │ - ID (PK)               │
              │ - RideID                │
              │ - SenderID              │
              │ - Content               │
              │ - Timestamp             │
              └─────────────────────────┘
```

### Table Definitions

#### **Users Table**
```sql
CREATE TABLE users (
    id VARCHAR(191) PRIMARY KEY,
    name VARCHAR(191) NOT NULL,
    email VARCHAR(191) UNIQUE NOT NULL,
    photo_url VARCHAR(191),
    geohash VARCHAR(64),
    last_seen TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### **Ride Offers Table**
```sql
CREATE TABLE ride_offers (
    id VARCHAR(191) PRIMARY KEY,
    driver_id VARCHAR(191) NOT NULL,
    from_geo VARCHAR(64) NOT NULL,
    to_geo VARCHAR(64) NOT NULL,
    fare VARCHAR(191),
    time TIMESTAMP NOT NULL,
    seats INT NOT NULL,
    status VARCHAR(191) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (driver_id) REFERENCES users(id)
);
```

#### **Ride Requests Table**
```sql
CREATE TABLE ride_requests (
    id VARCHAR(191) PRIMARY KEY,
    user_id VARCHAR(191) NOT NULL,
    from_geo VARCHAR(64) NOT NULL,
    from_geo VARCHAR(64) NOT NULL,
    time TIMESTAMP NOT NULL,
    seats INT NOT NULL,
    status VARCHAR(191) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

#### **Matches Table**
```sql
CREATE TABLE matches (
    id VARCHAR(191) PRIMARY KEY,
    rider_id VARCHAR(191) NOT NULL,
    driver_id VARCHAR(191) NOT NULL,
    ride_id VARCHAR(191) NOT NULL,
    status VARCHAR(191) DEFAULT 'requested',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (rider_id) REFERENCES users(id),
    FOREIGN KEY (driver_id) REFERENCES users(id),
    FOREIGN KEY (ride_id) REFERENCES ride_offers(id)
);
```

#### **Chat Messages Table**
```sql
CREATE TABLE chat_messages (
    id VARCHAR(191) PRIMARY KEY,
    ride_id VARCHAR(191) NOT NULL,
    sender_id VARCHAR(191) NOT NULL,
    content TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ride_id) REFERENCES ride_offers(id),
    FOREIGN KEY (sender_id) REFERENCES users(id)
);
```

#### **Reviews Table**
```sql
CREATE TABLE reviews (
    id VARCHAR(191) PRIMARY KEY,
    ride_id VARCHAR(191) NOT NULL,
    from_user_id VARCHAR(191) NOT NULL,
    to_user_id VARCHAR(191) NOT NULL,
    score INT NOT NULL CHECK (score >= 1 AND score <= 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ride_id) REFERENCES ride_offers(id),
    FOREIGN KEY (from_user_id) REFERENCES users(id),
    FOREIGN KEY (to_user_id) REFERENCES users(id)
);
```

#### **User Locations Table**
```sql
CREATE TABLE user_locations (
    user_id VARCHAR(191) PRIMARY KEY,
    latitude DOUBLE NOT NULL,
    longitude DOUBLE NOT NULL,
    geohash VARCHAR(64) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Indexing Strategy

#### **Primary Indexes**
- **Users**: `id` (PRIMARY), `email` (UNIQUE)
- **Ride Offers**: `id` (PRIMARY), `driver_id` (FOREIGN KEY)
- **Ride Requests**: `id` (PRIMARY), `user_id` (FOREIGN KEY)
- **Matches**: `id` (PRIMARY), `ride_id` (FOREIGN KEY)

#### **Spatial Indexes**
- **Geohash Prefix**: `from_geo LIKE 'prefix%'` for proximity search
- **Location Updates**: `updated_at` for real-time location tracking
- **Time-based Queries**: `time` index for ride scheduling

#### **Performance Optimizations**
- **Composite Indexes**: `(from_geo, status)` for active ride queries
- **Covering Indexes**: Include frequently accessed fields
- **Query Optimization**: Geohash-based spatial queries

## Authentication & Security

### Security Architecture

#### **Multi-Layer Security**
1. **Transport Security**: gRPC over HTTP/2 with TLS
2. **Authentication**: Google OAuth 2.0 with JWT validation
3. **Authorization**: Domain-based access control
4. **Data Protection**: Input validation and sanitization

#### **OAuth 2.0 Flow**
1. **Client Authentication**: Google OAuth client ID validation
2. **Token Verification**: Google token info endpoint validation
3. **Domain Validation**: Educational institution email verification
4. **User Provisioning**: Automatic account creation for new users

#### **JWT Implementation**
- **Algorithm**: HMAC-SHA256 (HS256)
- **Claims**: User ID, email, name, issued at, expiration
- **Expiration**: 24-hour token lifetime
- **Secret Management**: Environment-based JWT secret

### Access Control

#### **Public Endpoints**
- **Login**: `/proto.v1.AuthService/Login`
- **Health Check**: Server status endpoints

#### **Protected Endpoints**
- **User Management**: Profile updates, user discovery
- **Ride Operations**: Create, update, delete rides
- **Matching**: Request management, status updates
- **Communication**: Chat messages, reviews
- **Location**: Location updates, proximity queries

#### **Domain Restrictions**
- **Email Validation**: Restricted to allowed domains
- **Institution Verification**: Educational email verification
- **Access Control**: Domain-based user provisioning

### Security Best Practices

#### **Input Validation**
- **Request Sanitization**: Validate all input parameters
- **Type Safety**: Strong typing with Protocol Buffers
- **Business Rules**: Enforce domain-specific validation

#### **Error Handling**
- **Information Leakage**: Generic error messages
- **Logging**: Secure logging without sensitive data
- **Monitoring**: Security event tracking

#### **Data Protection**
- **Encryption**: Sensitive data encryption at rest
- **Access Logging**: Audit trail for security events
- **Privacy**: Minimal data collection and retention

## Installation & Setup

### Prerequisites

#### **System Requirements**
- **Operating System**: Linux, macOS, or Windows
- **Go Version**: 1.24.4 or higher
- **Memory**: Minimum 2GB RAM
- **Storage**: 10GB available disk space

#### **Software Dependencies**
- **Go**: Official Go distribution
- **MySQL**: 8.0 or higher
- **Git**: Version control system
- **Make**: Build automation (optional)

#### **Development Tools**
- **Protocol Buffers**: protoc compiler
- **gRPC Tools**: Go gRPC plugins
- **Wire**: Dependency injection tool

### Environment Setup

#### **1. Go Environment**
```bash
# Install Go 1.24.4+
wget https://golang.org/dl/go1.24.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.4.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Verify installation
go version
```

#### **2. Protocol Buffers**
```bash
# Install protoc compiler
sudo apt-get install protobuf-compiler

# Install Go gRPC plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/google/wire/cmd/wire@latest
```

#### **3. MySQL Database**
```bash
# Install MySQL
sudo apt-get install mysql-server

# Start MySQL service
sudo systemctl start mysql
sudo systemctl enable mysql

# Secure installation
sudo mysql_secure_installation
```

### Project Setup

#### **1. Clone Repository**
```bash
git clone https://github.com/your-username/hope.git
cd hope
```

#### **2. Install Dependencies**
```bash
go mod download
go mod tidy
```

#### **3. Environment Configuration**
```bash
# Copy environment template
cp .env.example .env

# Configure environment variables
cat > .env << EOF
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=hope_user
DB_PASSWORD=secure_password
DB_NAME=hope_db

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here

# Google OAuth
GOOGLE_CLIENT_ID=your-google-oauth-client-id

# Allowed Domains
ALLOWED_DOMAINS=university.edu,college.edu

# Server Configuration
GRPC_PORT=8080
EOF
```

#### **4. Database Setup**
```bash
# Create database and user
mysql -u root -p << EOF
CREATE DATABASE hope_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'hope_user'@'localhost' IDENTIFIED BY 'secure_password';
GRANT ALL PRIVILEGES ON hope_db.* TO 'hope_user'@'localhost';
FLUSH PRIVILEGES;
EOF
```

#### **5. Generate Code**
```bash
# Generate Protocol Buffer code
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/v1/*.proto

# Generate dependency injection code
wire ./di
```

#### **6. Build Application**
```bash
# Build binary
go build -o hope main.go

# Verify build
./hope --help
```

### Configuration Options

#### **Environment Variables**
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_HOST` | MySQL host address | localhost | Yes |
| `DB_PORT` | MySQL port number | 3306 | Yes |
| `DB_USER` | Database username | - | Yes |
| `DB_PASSWORD` | Database password | - | Yes |
| `DB_NAME` | Database name | - | Yes |
| `JWT_SECRET` | JWT signing secret | - | Yes |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID | - | Yes |
| `ALLOWED_DOMAINS` | Comma-separated allowed domains | - | Yes |
| `GRPC_PORT` | gRPC server port | 8080 | No |

#### **Database Configuration**
```go
type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Database string
}
```

#### **JWT Configuration**
```go
type Config struct {
    JWTSecret     []byte
    PublicMethods map[string]bool
}
```

## Development

### Development Workflow

#### **1. Code Structure**
- **Clean Architecture**: Clear separation of concerns
- **Interface Design**: Contract-first development
- **Error Handling**: Consistent error patterns
- **Testing**: Comprehensive test coverage

#### **2. Development Tools**
- **Go Modules**: Dependency management
- **Wire**: Dependency injection
- **Protocol Buffers**: API contract definition
- **GORM**: Database ORM

#### **3. Code Quality**
- **Linting**: Go vet and static analysis
- **Formatting**: gofmt for consistent style
- **Documentation**: Comprehensive code comments
- **Testing**: Unit and integration tests

### Testing Strategy

#### **Unit Testing**
```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./service

# Run with coverage
go test -cover ./...
```

#### **Integration Testing**
```bash
# Run with database integration
go test -tags=integration ./...

# Run specific integration tests
go test -tags=integration ./repository
```

#### **Test Structure**
```
service/
├── auth_service_test.go
├── user_service_test.go
├── ride_service_test.go
└── match_service_test.go
```

### API Development

#### **Protocol Buffer Updates**
1. **Define Service**: Add new RPC methods to .proto files
2. **Generate Code**: Run protoc to generate Go code
3. **Implement Service**: Add business logic to service layer
4. **Add Handler**: Implement gRPC handler
5. **Update Wire**: Add new dependencies to DI configuration

#### **Example: Adding New Endpoint**
```protobuf
// proto/v1/user.proto
service UserService {
    rpc GetUserProfile(GetUserProfileRequest) returns (GetUserProfileResponse);
}

message GetUserProfileRequest {
    string user_id = 1;
}

message GetUserProfileResponse {
    UserProfile profile = 1;
}
```

```go
// service/user_service.go
func (s *userService) GetUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
    // Implementation
}

// api/user_handler.go
func (h *UserHandler) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
    // Handler implementation
}
```

### Database Development

#### **Schema Changes**
1. **Model Updates**: Modify Go structs in `db/` package
2. **Migration**: GORM auto-migration handles schema updates
3. **Repository Updates**: Update repository methods if needed
4. **Testing**: Verify changes with integration tests

#### **Example: Adding New Field**
```go
// db/user.go
type User struct {
    ID       string    `gorm:"primaryKey"`
    Name     string    `gorm:"size:191"`
    Email    string    `gorm:"uniqueIndex;size:191"`
    Phone    string    `gorm:"size:20"` // New field
    PhotoURL string    `gorm:"size:191"`
    Geohash  string    `gorm:"size:64"`
    LastSeen time.Time
}
```

### Performance Optimization

#### **Database Optimization**
- **Indexing**: Strategic index placement for common queries
- **Query Optimization**: Efficient geohash-based spatial queries
- **Connection Pooling**: Optimized database connection management
- **Caching**: Redis integration for frequently accessed data

#### **Application Optimization**
- **Goroutines**: Concurrent request processing
- **Connection Reuse**: Efficient gRPC connection management
- **Memory Management**: Optimized data structures
- **Profiling**: Performance monitoring and optimization

## Deployment

### Production Deployment

#### **1. Environment Preparation**
```bash
# Production environment variables
cat > .env.prod << EOF
DB_HOST=prod-mysql.example.com
DB_PORT=3306
DB_USER=hope_prod_user
DB_PASSWORD=super_secure_prod_password
DB_NAME=hope_prod_db
JWT_SECRET=production-jwt-secret-key
GOOGLE_CLIENT_ID=prod-google-client-id
ALLOWED_DOMAINS=university.edu
GRPC_PORT=8080
EOF
```

#### **2. Build Production Binary**
```bash
# Build for production
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o hope main.go

# Verify binary
file hope
./hope --version
```

#### **3. Database Migration**
```bash
# Run production migrations
./hope migrate --env=production

# Verify schema
mysql -h prod-mysql.example.com -u hope_prod_user -p hope_prod_db -e "SHOW TABLES;"
```

#### **4. Service Deployment**
```bash
# Create systemd service
sudo tee /etc/systemd/system/hope.service << EOF
[Unit]
Description=Hope Carpooling Platform
After=network.target

[Service]
Type=simple
User=hope
WorkingDirectory=/opt/hope
ExecStart=/opt/hope/hope
Restart=always
RestartSec=5
EnvironmentFile=/opt/hope/.env.prod

[Install]
WantedBy=multi-user.target
EOF

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable hope
sudo systemctl start hope
```

### Container Deployment

#### **Docker Configuration**
```dockerfile
# Dockerfile
FROM golang:1.24.4-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hope .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/hope .
COPY --from=builder /app/.env.prod .env

EXPOSE 8080
CMD ["./hope"]
```

#### **Docker Compose**
```yaml
# docker-compose.yml
version: '3.8'

services:
  hope:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=hope_user
      - DB_PASSWORD=hope_password
      - DB_NAME=hope_db
    depends_on:
      - mysql
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: hope_db
      MYSQL_USER: hope_user
      MYSQL_PASSWORD: hope_password
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "3306:3306"

volumes:
  mysql_data:
```

### Monitoring & Observability

#### **Health Checks**
```go
// Health check endpoint
func (s *Server) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
    return &pb.HealthCheckResponse{
        Status: "healthy",
        Timestamp: timestamppb.Now(),
    }, nil
}
```

#### **Logging Strategy**
- **Structured Logging**: JSON-formatted logs
- **Log Levels**: Debug, Info, Warn, Error
- **Context Information**: Request ID, user ID, operation
- **Performance Metrics**: Response time, throughput

#### **Metrics Collection**
- **Prometheus**: Application metrics
- **Grafana**: Visualization and alerting
- **Custom Metrics**: Business-specific KPIs
- **Alerting**: Automated issue detection

## Contributing

### Development Guidelines

#### **Code Standards**
- **Go Format**: Use `gofmt` for code formatting
- **Naming Conventions**: Follow Go naming conventions
- **Error Handling**: Consistent error handling patterns
- **Documentation**: Comprehensive code comments

#### **Commit Guidelines**
```
feat: add new user profile endpoint
fix: resolve geohash proximity search issue
docs: update API documentation
test: add integration tests for matching service
refactor: improve error handling in auth service
```

#### **Pull Request Process**
1. **Fork Repository**: Create personal fork
2. **Feature Branch**: Create feature-specific branch
3. **Implementation**: Implement feature with tests
4. **Testing**: Ensure all tests pass
5. **Documentation**: Update relevant documentation
6. **Submit PR**: Create pull request with description

### Testing Requirements

#### **Test Coverage**
- **Unit Tests**: Minimum 80% coverage
- **Integration Tests**: Database and service integration
- **API Tests**: gRPC endpoint testing
- **Performance Tests**: Load and stress testing

#### **Test Structure**
```
service/
├── auth_service.go
├── auth_service_test.go
├── user_service.go
└── user_service_test.go
```

### Documentation Standards

#### **Code Documentation**
- **Package Comments**: Describe package purpose
- **Function Comments**: Explain function behavior
- **Type Comments**: Document struct fields
- **Example Usage**: Provide usage examples

#### **API Documentation**
- **Protocol Buffers**: Clear service definitions
- **Request/Response**: Document message formats
- **Error Codes**: Standardized error responses
- **Examples**: Request/response examples

### Review Process

#### **Code Review Checklist**
- [ ] **Functionality**: Feature works as expected
- [ ] **Code Quality**: Follows Go best practices
- [ ] **Testing**: Adequate test coverage
- [ ] **Documentation**: Updated documentation
- [ ] **Performance**: No performance regressions
- [ ] **Security**: No security vulnerabilities

#### **Review Guidelines**
- **Constructive Feedback**: Provide helpful suggestions
- **Code Quality**: Focus on maintainability
- **Testing**: Ensure adequate test coverage
- **Documentation**: Verify documentation updates

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support and questions:
- **Issues**: [GitHub Issues](https://github.com/your-username/hope/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-username/hope/discussions)
- **Documentation**: [Project Wiki](https://github.com/your-username/hope/wiki)

## Acknowledgments

- **Google**: OAuth 2.0 and Protocol Buffers
- **GORM**: Database ORM framework
- **gRPC**: High-performance RPC framework
- **Go Community**: Language and ecosystem support

---

*Built with ❤️ for sustainable transportation and community building.* 