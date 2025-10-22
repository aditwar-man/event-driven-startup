# 🚀 SMM Platform - Social Media Management Platform

A modern, scalable Social Media Management platform built with microservices architecture, event-driven design, and advanced user management features.

## 📋 Overview

The SMM Platform helps e-commerce businesses manage their social media presence across multiple platforms with AI-powered content generation, automated posting, and comprehensive analytics.

## 🏗️ Architecture

### Microservices Structure

| Service | Port | Description |
|---------|------|-------------|
| **Auth Service** | 8081 | Authentication, authorization, user management |
| **User Service** | 8082 | User profiles, quota management, role-based access |
| **Product Service** | 8083 | Product catalog management (Coming Soon) |
| **Posting Service** | 8084 | Social media posting & scheduling (Coming Soon) |
| **Analytics Service** | 8085 | Performance analytics & insights (Coming Soon) |
| **AI Content Service** | 8086 | AI-powered content generation (Coming Soon) |

### Technology Stack

- **Backend**: Go 1.21+ (Gin framework)
- **Database**: PostgreSQL (per service)
- **Message Broker**: Apache Kafka
- **Caching**: Redis
- **Containerization**: Docker & Docker Compose
- **API Documentation**: Swagger/OpenAPI
- **Authentication**: JWT with refresh tokens

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.21+ (for local development)

### Running with Docker

1. **Clone and setup the project:**
```bash
git clone <repository-url>
cd smm-platform
```

2. **Start all services:**
```bash
docker-compose up -d
```

3. **Verify services are running:**
```bash
docker-compose ps
```

### Access Points

| Service | URL | Description |
|---------|-----|-------------|
| **Auth Service** | http://localhost:8081 | Authentication API |
| **User Service** | http://localhost:8082 | User management API |
| **Auth Swagger** | http://localhost:8081/swagger/index.html | API Documentation |
| **User Swagger** | http://localhost:8082/swagger/index.html | API Documentation |
| **Kafka UI** | http://localhost:8090 | Kafka management interface |
| **Jaeger UI** | http://localhost:16686 | Distributed tracing |

## 🔧 Development

### Local Development Setup

1. **Install dependencies:**
```bash
# Generate go.sum files
cd shared && go mod tidy && cd ..
cd services/auth && go mod tidy && cd ../..
cd services/user && go mod tidy && cd ../..
```

2. **Build and run services:**
```bash
# Build auth service
cd services/auth
go build -o main ./cmd
./main

# Build user service (in another terminal)
cd services/user  
go build -o main ./cmd
./main
```

### Project Structure

```
smm-platform/
├── shared/                 # Shared libraries
│   ├── pkg/
│   │   ├── domain/        # Shared domain models
│   │   ├── events/        # Event definitions & bus
│   │   └── database/      # Database utilities
├── services/
│   ├── auth/              # Authentication service
│   │   ├── cmd/           # Main application
│   │   ├── internal/
│   │   │   ├── domain/    # Domain models
│   │   │   ├── application/ # Use cases & services
│   │   │   └── infrastructure/ # External implementations
│   │   └── docs/          # Swagger documentation
│   └── user/              # User management service
├── scripts/               # Database initialization
└── docker-compose.yml     # Container orchestration
```

## 🔐 Authentication & Security

### Features

- **JWT-based authentication** with access/refresh tokens
- **Role-Based Access Control (RBAC)** with dynamic roles
- **Session management** with device tracking
- **Password strength validation**
- **Rate limiting** and security headers
- **Secure password hashing** with bcrypt

### Default Roles

| Role | Permissions | Description |
|------|-------------|-------------|
| **super_admin** | All permissions | Full system access |
| **admin** | User management, content management | Administrative access |
| **user** | Basic content creation, analytics | Standard user |

## 📊 User Management & Quotas

### Tier System

| Tier | AI Descriptions | AI Videos | Auto Posts | Price |
|------|-----------------|-----------|------------|-------|
| **Free** | 5/month | 0 | 5/month | Free |
| **Pro** | 100/month | 10/month | Unlimited | $29/month |

### Quota Management

- **Monthly quota reset** automated system
- **Real-time quota tracking**
- **Graceful quota exhaustion handling**
- **Admin quota management**

## 🔄 Event-Driven Architecture

### Key Events

- `user.registered` - New user registration
- `user.tier.upgraded` - User tier change
- `user.quota.updated` - Quota usage updates
- `content.scheduled` - Post scheduling
- `content.published` - Post publication

### Event Flow

```
User Action → Service → Event → Kafka → Consumer Services
```

## 🗄️ Database Schema

### Auth Service
- `users` - User accounts and authentication
- `sessions` - Active user sessions
- `roles` - System roles and permissions
- `user_roles` - Role assignments

### User Service  
- `users` - User profiles and quotas
- `user_preferences` - User settings
- `quota_usage` - Quota tracking

## 🧪 Testing

### Running Tests

```bash
# Test auth service
cd services/auth
go test ./...

# Test user service
cd services/user  
go test ./...
```

### API Testing

Use the Swagger documentation or import the OpenAPI spec into tools like:
- Postman
- Insomnia
- curl

## 📈 Monitoring & Observability

### Logging
- Structured JSON logging
- Correlation IDs for request tracing
- Log levels (DEBUG, INFO, WARN, ERROR)

### Metrics
- Request latency
- Error rates
- Quota usage metrics
- Business metrics

### Tracing
- Distributed tracing with Jaeger
- Kafka message tracing
- Database query tracing

## 🔒 Security Best Practices

- **JWT token expiration** (15min access, 7day refresh)
- **Password hashing** with bcrypt
- **CORS configuration**
- **Rate limiting** per user/IP
- **Input validation** and sanitization
- **SQL injection prevention**

## 🚢 Deployment

### Production Checklist

- [ ] Set proper JWT secrets
- [ ] Configure database connections
- [ ] Set up monitoring and alerting
- [ ] Configure backup strategies
- [ ] Set up SSL/TLS certificates
- [ ] Configure firewall rules
- [ ] Set up log aggregation

### Environment Variables

Key environment variables for each service:

```bash
# Auth Service
JWT_SECRET=your-super-secret-key
DB_HOST=postgres-auth
DB_PASSWORD=secure-password
REDIS_HOST=redis
KAFKA_BROKERS=kafka:9092

# User Service
DB_HOST=postgres-user
DB_PASSWORD=secure-password
KAFKA_BROKERS=kafka:9092
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

### Development Guidelines

- Follow Go best practices and style guide
- Write comprehensive tests
- Update API documentation
- Use conventional commit messages
- Ensure backward compatibility

## 📄 License

MIT License - see LICENSE file for details

## 🆘 Support

- 📚 [API Documentation](http://localhost:8081/swagger/index.html)
- 🐛 [Issue Tracker](https://github.com/aditwar-mann/event-driven-startup/issues)
- 💬 [Discussion Forum](https://github.com/aditwar-mann/event-driven-startup/discussions)

## 🏆 Features Status

| Feature | Status | Version |
|---------|--------|---------|
| User Authentication | ✅ Complete | v1.0 |
| Role-Based Access Control | ✅ Complete | v1.0 |
| Quota Management | ✅ Complete | v1.0 |
| Product Management | 🚧 In Progress | v1.1 |
| Social Media Posting | 🚧 In Progress | v1.1 |
| AI Content Generation | 🚧 In Progress | v1.2 |
| Advanced Analytics | 📋 Planned | v1.3 |

---

**Built with ❤️ using Go, Docker, and modern microservices patterns**