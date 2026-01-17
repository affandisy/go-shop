### Belajar membuat GoShop
Official Repository: <br>
[![GitHub](https://img.shields.io/badge/GitHub-GoShop%20Project-181717?style=for-the-badge&logo=github&logoColor=white)](https://github.com/quangdangfit/goshop)

## ✨ Features

- 🔐 **JWT Authentication** - Secure user authentication
- 👥 **User Management** - Registration, login, profile management
- 📦 **Product Management** - CRUD with filtering, search, and pagination
- 🏷️ **Category Management** - Organize products by categories
- 🛒 **Order Management** - Complete checkout and order workflow
- 📊 **Stock Management** - Automatic stock tracking
- 🔄 **Order Status Workflow** - pending → paid → processing → shipped → delivered
- 🎫 **Role-Based Access** - Customer and Admin roles
- 🗑️ **Soft Delete** - Safe data deletion with audit trail
- 📄 **Pagination** - Efficient data retrieval
- 🔍 **Advanced Filtering** - Search by name, category, price range

## 🏗️ Architecture

Built with **Clean Architecture (Domain-Driven Design)**:

```
Handler → Service → Repository → Database
```

- **Handler**: HTTP request handling & validation
- **Service**: Business logic & use cases
- **Repository**: Data access layer
- **Domain**: Business models & entities

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (optional)

## 📚 API Documentation

### Base URL
```
http://localhost:8888/api/v1
```

### Authentication Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Register new user |
| POST | `/auth/login` | Login user |

### User Endpoints
| Method | Endpoint | Auth | Admin | Description |
|--------|----------|------|-------|-------------|
| GET | `/users/profile` | ✅ | ❌ | Get my profile |
| PUT | `/users/profile` | ✅ | ❌ | Update profile |
| GET | `/users` | ✅ | ✅ | Get all users |

### Category Endpoints
| Method | Endpoint | Auth | Admin | Description |
|--------|----------|------|-------|-------------|
| GET | `/categories` | ❌ | ❌ | Get all categories |
| GET | `/categories/:id` | ❌ | ❌ | Get category |
| POST | `/categories` | ✅ | ✅ | Create category |
| PUT | `/categories/:id` | ✅ | ✅ | Update category |
| DELETE | `/categories/:id` | ✅ | ✅ | Delete category |

### Product Endpoints
| Method | Endpoint | Auth | Admin | Description |
|--------|----------|------|-------|-------------|
| GET | `/products` | ❌ | ❌ | Get products (with filters) |
| GET | `/products/:id` | ❌ | ❌ | Get product |
| POST | `/products` | ✅ | ✅ | Create product |
| PUT | `/products/:id` | ✅ | ✅ | Update product |
| DELETE | `/products/:id` | ✅ | ✅ | Delete product |
| PATCH | `/products/:id/stock` | ✅ | ✅ | Update stock |

### Order Endpoints
| Method | Endpoint | Auth | Admin | Description |
|--------|----------|------|-------|-------------|
| POST | `/orders` | ✅ | ❌ | Create order |
| GET | `/orders` | ✅ | ❌ | Get my orders |
| GET | `/orders/:id` | ✅ | ❌ | Get order detail |
| GET | `/orders/all` | ✅ | ✅ | Get all orders |
| PATCH | `/orders/:id/status` | ✅ | ✅ | Update status |
| POST | `/orders/:id/cancel` | ✅ | ❌ | Cancel order |

**Product Filters:**
- `?name=iPhone` - Search by name
- `?category_id=uuid` - Filter by category
- `?min_price=1000000&max_price=5000000` - Price range
- `?page=1&limit=10` - Pagination