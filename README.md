# Web Sportwear Shop 🏀

A production-ready microservices-based e-commerce platform built with Go, featuring JWT authentication, role-based access control, and a modern JavaScript frontend. This project demonstrates full-stack development capabilities with enterprise-grade architecture patterns.

[![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-316192?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?style=flat&logo=redis)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker)](https://www.docker.com/)

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Project Structure](#project-structure)
- [Security Features](#security-features)
- [Performance Optimizations](#performance-optimizations)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

## Overview

Web Sportwear Shop is a comprehensive e-commerce platform designed for athletic gear sales. It showcases modern backend development practices including microservices architecture, distributed systems, caching strategies, and secure authentication flows.

**Key Highlights:**
- Microservices architecture with gRPC communication
-  Secure JWT authentication with OAuth 2.0 integration
-  Redis caching with ETag-based HTTP caching
-  Fully containerized with Docker Compose
-  HTTPS/TLS with Nginx reverse proxy
-  Digital wallet system with atomic transactions
-  Three-tier RBAC (Customer, Seller, Admin)

## Features

### User Management
-  **Secure Authentication**: JWT-based auth with HttpOnly cookies
-  **OAuth 2.0 Integration**: Google & GitHub social login via OpenID Connect
-  **Role-Based Access Control (RBAC)**: Customer, Seller, and Admin roles
-  **Digital Wallet System**: Built-in payment processing with transaction history
-  **Password Security**: bcrypt hashing with cost factor 12

### Product Management
-  **Full CRUD Operations**: Create, read, update, delete products
- ️ **Category System**: Organized product categorization with hierarchy
- ️ **Image Support**: Product image URLs with placeholder fallbacks
-  **Advanced Filtering**: Filter by category, seller, price, and search
-  **Seller Dashboard**: Personal inventory management interface

### Shopping Experience
-  **Shopping Cart**: Distributed cart service via gRPC microservice
-  **Order Management**: Complete order workflow with delivery tracking
-  **Shipping Options**: Multiple delivery speed tiers (Ordinary, Medium, Fastest)
-  **Payment Processing**: Wallet-based transactions with atomic database operations
-  **Order Notifications**: Real-time order status updates

### Performance & Scalability
-  **Redis Caching**: Application-level caching for frequently accessed data
- ️ **HTTP Caching**: ETags and cache-control headers for bandwidth optimization
-  **Database Connection Pooling**: Optimized PostgreSQL connections (20 max, 10 min idle)
-  **Load Balancing**: Nginx reverse proxy with upstream configuration
-  **HTTPS/TLS**: Secure communications with SSL certificates
-  **Horizontal Scaling**: Docker Compose scale-out capability

## Technology Stack

### Backend Technologies
- **Language**: Go 1.24 
- **Web Framework**: Gorilla Mux (Routing, Middleware)
- **Database**: PostgreSQL 17 (ACID compliance, Transactions)
- **Cache**: Redis 7 (for caching)
- **Authentication**:
    - JWT: golang-jwt/jwt v4
    - OAuth 2.0: Goth library
    - OpenID Connect: Google & GitHub providers
- **Security**:
    - bcrypt (Password hashing)
    - CORS: rs/cors
    - CSP Headers
- **API Protocols**:
    - REST (HTTP/2)
    - gRPC (Protocol Buffers)
- **Database Drivers**:
    - lib/pq (PostgreSQL driver)
    - jmoiron/sqlx (Extended database/sql)
- **Migrations**: golang-migrate/migrate v4

### Microservices Architecture
- **Cart Service**: gRPC-based distributed cart management
    - Protocol: gRPC (github.com/Asylann/grpc-demo)
    - Port: 50051
    - Database: Separate PostgreSQL instance
- **Order Service**: gRPC-based order processing
    - Protocol: gRPC (github.com/Asylann/orderservicegrpc)
    - Port: 50052
    - Database: Separate PostgreSQL instance
- **Communication**: Protocol Buffers (protobuf) for efficient serialization

### Frontend Technologies
- **Core**: Vanilla JavaScript (ES6+, Async/Await, Fetch API)
- **Styling**: Custom CSS with Adidas-inspired design system
- **Architecture**:
    - Module-based with separation of concerns
    - MVC-like pattern (Models, Views, Controllers)
- **State Management**:
    - LocalStorage for persistent client state
    - In-memory state for session data
- **HTTP Client**: Fetch API with credential management
- **Build**: No build step (native ES6 modules)

### Infrastructure & DevOps
- **Containerization**:
    - Docker 24+
    - Docker Compose v2
    - Multi-stage builds
- **Reverse Proxy**:
    - Nginx 1.27.2
    - SSL/TLS termination
    - HTTP/2 support
    - gzip compression
- **Service Orchestration**:
    - Docker Compose multi-service setup
    - Service discovery via Docker networking
    - Health check integration
- **Environment Management**:
    - dotenv configuration
    - Environment-specific configs
- **Monitoring & Logging**:
    - Nginx access/error logs
    - Structured application logging
    - Docker logs aggregation

### Development Tools
- **Version Control**: Git
- **Testing**:
    - Go testing package
    - httptest for HTTP handlers
    - Integration tests
- **API Testing**:
    - Built-in test suite
- **Code Quality**:
    - Go fmt
    - Go vet
    - Race detector

### Database Schema
- **PostgreSQL Databases**:
    - Main application database (web_sportwear_db)
    - Cart microservice database (grpc_db)
    - Order microservice database (order_db)
- **Tables**:
    - users, roles, products, categories
    - wallets, transactions
    - etag_versions (cache invalidation)
    - Cart items (microservice)
    - Orders (microservice)

## Architecture

### System Architecture Diagram
```
┌─────────────────────────────────────────────────────────────────┐
│                         Nginx (8081)                            │
│         Reverse Proxy | SSL/TLS | Load Balancer                │
│              HTTP/2 | gzip | Static Files                      │
└────────────┬────────────────────────────────┬───────────────────┘
             │                                │
             │ HTTPS                          │ Static Content
             │                                │
    ┌────────▼────────┐              ┌────────▼────────┐
    │   Main API      │              │   Frontend      │
    │   (Go/REST)     │              │ (HTML/JS/CSS)   │
    │   Port 8080     │              │   Static Files  │
    └────────┬────────┘              └─────────────────┘
             │
             ├──── gRPC ────┐
             │              │
    ┌────────▼────────┐   ┌▼─────────────┐
    │  Cart Service   │   │Order Service │
    │    (gRPC)       │   │   (gRPC)     │
    │   Port 50051    │   │  Port 50052  │
    └────────┬────────┘   └──────┬───────┘
             │                    │
             │                    │
    ┌────────▼────────────────────▼───────────────────┐
    │           PostgreSQL Cluster                    │
    │  ┌──────────┐  ┌──────────┐  ┌──────────┐     │
    │  │ Main DB  │  │ Cart DB  │  │Order DB  │     │
    │  │  (5432)  │  │  (5432)  │  │  (5432)  │     │
    │  └──────────┘  └──────────┘  └──────────┘     │
    └─────────────────────────────────────────────────┘
             │
    ┌────────▼────────┐
    │      Redis      │
    │   Port 6379     │
    │  (Caching Layer)│
    └─────────────────┘
```

### Key Architectural Patterns

1. **Microservices Architecture**
    - Separated cart and order services for independent scaling
    - Service-to-service communication via gRPC
    - Database per service pattern

2. **Repository Pattern**
    - Clean separation of data access layer
    - Abstraction over database operations
    - Testable data layer

3. **Middleware Chain**
    - Logging → Authentication → Authorization → Handler
    - Composable and reusable middleware functions
    - Context-based request metadata passing

4. **MVC-like Structure**
    - Handlers (Controllers)
    - Models (Domain entities)
    - Database (Data access)

5. **API Gateway Pattern**
    - Nginx as central entry point
    - Request routing and load balancing
    - SSL/TLS termination

6. **Caching Strategy**
    - Redis for application cache
    - HTTP ETags for client-side caching
    - Cache invalidation on mutations

7. **Transaction Script**
    - Database transactions for financial operations
    - Atomic wallet transfers
    - Rollback on errors

## Getting Started

### Prerequisites

Ensure you have the following installed:

- **Docker** 24.0+ & **Docker Compose** v2.0+
- **Go** 1.24+ (for local development)
- **PostgreSQL** 17+ (for local development without Docker)
- **Redis** 7+ (for local development without Docker)
- **mkcert** (for local SSL certificate generation)

### Installation

#### 1. Clone the repository
```bash
git clone https://github.com/Asylann/WebSportwearShop.git
cd WebSportwearShop
```

#### 2. Generate SSL Certificates
```bash
# Install mkcert (if not already installed)
# macOS
brew install mkcert

# Windows (using Chocolatey)
choco install mkcert

# Linux
apt install mkcert  # or equivalent for your distro

# Generate certificates
mkcert -install
mkcert localhost 127.0.0.1 ::1

# Certificates will be named: localhost+2.pem and localhost+2-key.pem
# Keep them in the project root directory
```

#### 3. Set up environment variables
```bash
# Copy the example env file
cp .env.example .env

# Edit .env with your configuration
nano .env  # or use your preferred editor
```

**Required environment variables:**
```env
# PostgreSQL Database Configuration
PGHOST=db
PGPORT=5432
PGUSER=postgres
PASSWORD=postgres
DATABASE=web_sportwear_db

# Redis Configuration
REDIS_HOST=redis:
REDIS_PORT=6379

# JWT Secret (change in production!)
JWT_SECRET=your_very_secure_secret_key_here

# OAuth 2.0 Credentials (Google)
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret

# OAuth 2.0 Credentials (GitHub)
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret

# Server Configuration
PORT=8080
ADMIN_ID=1
```

**Getting OAuth Credentials:**

<details>
<summary>Google OAuth Setup</summary>

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing
3. Navigate to "APIs & Services" > "Credentials"
4. Create OAuth 2.0 Client ID
5. Add authorized redirect URI: `https://localhost:8080/auth/google/callback`
6. Copy Client ID and Client Secret to `.env`

</details>

<details>
<summary>GitHub OAuth Setup</summary>

1. Go to GitHub Settings > Developer settings > OAuth Apps
2. Click "New OAuth App"
3. Set Homepage URL: `https://localhost:8081`
4. Set Authorization callback URL: `https://localhost:8080/auth/github/callback`
5. Copy Client ID and Client Secret to `.env`

</details>

#### 4. Start the services
```bash
# Start all services with Docker Compose
docker-compose up --build

# View logs
docker-compose logs -f

# Check service health
docker-compose ps
```

**Services will start in this order:**
1. PostgreSQL (3 instances: main, cart, order)
2. Redis
3. Database migrations (automatic)
4. Cart Service (gRPC)
5. Order Service (gRPC)
6. Main API (REST)
7. Nginx

#### 5. Access the application

- **Frontend**: https://localhost:8081
- **Main API**: https://localhost:8080
- **Nginx Status**: https://localhost:8081/nginx_status (from localhost only)

**Note**: You may see SSL warnings in your browser. This is expected with self-signed certificates. Click "Advanced" and proceed.

#### 8. Create your first user

Navigate to https://localhost:8081 and click "Sign up here". Create an account with:
- Email: admin@example.com
- Password: admin123
- Role: Admin (3)

## API Documentation

### Authentication Endpoints

| Method | Endpoint | Description | Access | Request Body | Response |
|--------|----------|-------------|--------|--------------|----------|
| POST | `/api/signup` | Create new user account | Public | `{email, password, roleId}` | `{data: user, err: ""}` |
| POST | `/api/login` | Login with credentials | Public | `{email, password}` | Sets auth_token cookie |
| GET | `/api/auth/google/login` | Initiate Google OAuth | Public | - | Redirects to Google |
| GET | `/api/auth/github/login` | Initiate GitHub OAuth | Public | - | Redirects to GitHub |
| GET | `/api/auth/google/callback` | Google OAuth callback | Public | Query params | Redirects to dashboard |
| GET | `/api/auth/github/callback` | GitHub OAuth callback | Public | Query params | Redirects to dashboard |
| POST | `/api/logout` | Logout current user | Authenticated | - | Clears auth_token |
| GET | `/api/me` | Get current user info | Authenticated | - | `{data: {id, email, roleId}}` |

### Product Endpoints

| Method | Endpoint | Description | Access | Request Body | Response |
|--------|----------|-------------|--------|--------------|----------|
| GET | `/api/products` | List all products | Public | - | `{data: [products]}` |
| GET | `/api/products/{id}` | Get product by ID | Customer+ | - | `{data: product}` |
| POST | `/api/products` | Create new product | Seller+ | `{name, description, price, size, category_id, imageURL, seller_id}` | `{data: product}` |
| PUT | `/api/products/{id}` | Update product | Seller+ | `{name, description, price, size, category_id, imageURL, seller_id}` | `{data: product}` |
| DELETE | `/api/products/{id}` | Delete product | Seller+ | - | `204 No Content` |
| GET | `/api/productsByCategory/{id}` | Products by category | Customer+ | - | `{data: [products]}` |
| GET | `/api/productsBySeller/{id}` | Products by seller | Customer+ | - | `{data: [products]}` |

### Category Endpoints

| Method | Endpoint | Description | Access | Request Body | Response |
|--------|----------|-------------|--------|--------------|----------|
| GET | `/api/categories` | List all categories | Customer+ | - | `{data: [categories]}` |
| GET | `/api/categories/{id}` | Get category by ID | Customer+ | - | `{data: category}` |
| POST | `/api/categories` | Create new category | Seller+ | `{name, description}` | `{data: category}` |
| PUT | `/api/categories/{id}` | Update category | Seller+ | `{name, description}` | `{data: category}` |
| DELETE | `/api/categories/{id}` | Delete category | Seller+ | - | `204 No Content` |

### Cart Endpoints

| Method | Endpoint | Description | Access | Request Body | Response |
|--------|----------|-------------|--------|--------------|----------|
| POST | `/api/carts` | Create new cart | Authenticated | - | `{data: "Cart is created"}` |
| GET | `/api/myCart` | Get user's cart items | Authenticated | - | `{data: [products]}` (with ETag) |
| POST | `/api/addToCart/{id}` | Add product to cart | Authenticated | - | `{data: "Added to your cart"}` |
| DELETE | `/api/myCart/{id}` | Remove item from cart | Authenticated | - | `{data: deletedProduct}` |

### Order Endpoints

| Method | Endpoint | Description | Access | Request Body | Response |
|--------|----------|-------------|--------|--------------|----------|
| POST | `/api/orders` | Create new order | Authenticated | `{transport_type, address}` | `{data: deliveryDate}` |
| GET | `/api/orders` | Get user's orders | Authenticated | - | `{data: [orders]}` |
| GET | `/api/orders/{id}` | Get order items | Authenticated | - | `{data: [products]}` |

### User Management Endpoints (Admin)

| Method | Endpoint | Description | Access | Request Body | Response |
|--------|----------|-------------|--------|--------------|----------|
| GET | `/api/users` | List all users | Admin | - | `{data: [users]}` (with ETag) |
| GET | `/api/users/{id}` | Get user by ID | Admin | - | `{data: user}` |
| PUT | `/api/users/{id}` | Update user | Admin | `{email, password, roleId}` | `{data: user}` |
| DELETE | `/api/users/{id}` | Delete user | Admin | - | `204 No Content` |
| GET | `/api/users/email/{id}` | Get user email | Customer+ | - | `{data: "email"}` |

### Wallet Endpoints

| Method | Endpoint | Description | Access | Request Body | Response |
|--------|----------|-------------|--------|--------------|----------|
| GET | `/api/wallet` | Get user's wallet info | Authenticated | - | `{data: {id, userId, balance, currency}}` |

## 📁 Project Structure
```
WebSportwearShop/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── cache/
│   │   └── redis.go               # Redis client initialization
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── db/
│   │   ├── initDB.go              # Database connection pool
│   │   ├── migration.go           # Migration runner
│   │   ├── product_store.go       # Product repository
│   │   ├── category_store.go      # Category repository
│   │   ├── user_store.go          # User repository
│   │   └── wallet_store.go        # Wallet repository
│   ├── handlers/
│   │   ├── auth.go                # Authentication handlers
│   │   ├── products.go            # Product CRUD handlers
│   │   ├── categories.go          # Category handlers
│   │   ├── users.go               # User management handlers
│   │   ├── carts.go               # Cart gRPC client
│   │   └── orders.go              # Order gRPC client
│   ├── middleware/
│   │   ├── jwt.go                 # JWT authentication
│   │   ├── roles.go               # RBAC middleware
│   │   └── logging.go             # Request logging
│   ├── models/
│   │   ├── product.go             # Product domain model
│   │   ├── category.go            # Category domain model
│   │   ├── user.go                # User domain model
│   │   └── wallet.go              # Wallet domain model
│   ├── httpresponse/
│   │   └── response.go            # Standardized API responses
│   └── tests/
│       ├── main_test.go           # Test setup
│       └── handlers_test.go       # Handler tests
├── frontend/
│   ├── css/
│   │   ├── style.css              # Global styles
│   │   ├── auth.css               # Authentication pages
│   │   ├── dashboard.css          # Dashboard styles
│   │   ├── components.css         # Reusable components
│   │   ├── orders.css             # Orders page
│   │   ├── product-detail.css     # Product detail page
│   │   └── seller.css             # Seller panel
│   ├── js/
│   │   ├── api.js                 # API client module
│   │   ├── auth.js                # Authentication logic
│   │   ├── utils.js               # Utility functions
│   │   ├── dashboard.js           # Dashboard controller
│   │   ├── seller.js              # Seller panel logic
│   │   ├── admin.js               # Admin panel logic
│   │   ├── orders.js              # Orders management
│   │   └── components.js          # Reusable UI components
│   ├── pages/
│   │   ├── dashboard.html         # User dashboard
│   │   ├── products.html          # Products listing
│   │   ├── product-detail.html    # Product details
│   │   ├── cart.html              # Shopping cart
│   │   ├── orders.html            # Order history
│   │   ├── seller.html            # Seller management
│   │   ├── admin.html             # Admin panel
│   │   └── signup.html            # User registration
│   ├── index.html                 # Landing/Login page
│   ├── server.js                  # Frontend dev server
│   └── package.json               # Frontend dependencies
├── migrations/
│   ├── 000001_create_products_table.up.sql
│   ├── 000002_create_categories_table.up.sql
│   ├── 000003_add_category_id_to_products.up.sql
│   ├── 000004_create_role_table.up.sql
│   ├── 000005_create_user_table.up.sql
│   ├── 000006_add_seller_id_to_products.up.sql
│   ├── 000007_create_wallets_table.up.sql
│   ├── 000008_create_transactions_table.up.sql
│   └── 000009_create_etag_table.up.sql
├── nginx/
│   ├── nginx.conf                 # Nginx configuration
│   ├── docker-compose.nginx.yml   # Nginx service definition
│   └── logs/                      # Access and error logs
├── docker-compose.yml             # Main orchestration file
├── docker-compose.db.yml          # Database services
├── Dockerfile                     # Main API container
├── go.mod                         # Go module dependencies
├── go.sum                         # Dependency checksums
├── .env.example                   # Environment template
└── README.md                      # This file
```


## Testing
Running Tests
```
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test file
go test ./internal/tests/handlers_test.go

# Run with race detection
go test -race ./...
```

## Service Health Checks
All services include health checks:

- PostgreSQL: pg_isready every 5 seconds
- Redis: redis-cli ping every 5 seconds
- Main API: HTTP health endpoint
- Nginx: Status page at /nginx_status

## Screenshots
Landing Page
<img src="https://github.com/Asylann/Web_Sportwear_Shop/blob/master/screenshots/loginPage.png">

User Dashboard
<img src="https://github.com/Asylann/Web_Sportwear_Shop/blob/master/screenshots/dashboardPage.png">

Products 
<img src="https://github.com/Asylann/Web_Sportwear_Shop/blob/master/screenshots/products.png">

Shopping Cart & Checkout
<img src="https://github.com/Asylann/Web_Sportwear_Shop/blob/master/screenshots/cartpage.png">

Admin Panel
<img src="https://github.com/Asylann/Web_Sportwear_Shop/blob/master/screenshots/adminPage.png">
## Author

GitHub: @Asylann