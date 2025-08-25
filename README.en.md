# goBlog

#### Description
goBlog is a modern blog system based on Go + Vue3, featuring a frontend-backend separation architecture with complete blog management functionality.

**Tech Stack:**
- **Backend**: Go + Gin + GORM + MySQL + Redis + Elasticsearch
- **Frontend**: Vue3 + TypeScript + Element Plus + Vite
- **AI Features**: Integrated DeepSeek AI Chat Assistant

#### Software Architecture

```
goBlog/
├── server/                    # Backend Service
│   ├── cmd/server/           # Application Entry Point
│   │   └── main.go          # Main Program Entry
│   ├── internal/             # Private Application and Library Code
│   │   ├── api/             # API Layer - HTTP Handlers
│   │   ├── service/         # Business Logic Layer
│   │   ├── model/           # Data Models
│   │   │   ├── database/    # Database Models
│   │   │   ├── request/     # Request Models
│   │   │   ├── response/    # Response Models
│   │   │   └── ...
│   │   ├── middleware/      # Middleware
│   │   ├── router/          # Route Configuration
│   │   ├── initialize/      # Initialization Modules
│   │   └── task/            # Scheduled Tasks
│   ├── pkg/                 # Reusable Library Code
│   │   ├── config/          # Configuration Management
│   │   ├── utils/           # Utility Functions
│   │   ├── global/          # Global Variables
│   │   └── core/            # Core Functionality
│   ├── configs/             # Configuration Files
│   ├── scripts/             # Script Files
│   ├── log/                 # Log Files
│   └── uploads/             # Upload Files
├── web/                     # Frontend Application
│   ├── src/
│   │   ├── components/      # Vue Components
│   │   ├── views/           # Page Views
│   │   ├── api/             # API Interfaces
│   │   └── utils/           # Utility Functions
│   └── ...
└── docs/                    # Project Documentation
```

#### Main Features

- **User Management**: User registration, login, permission management
- **Article Management**: Article publishing, editing, categorization, tagging
- **Comment System**: Article comments and reply functionality
- **AI Assistant**: Integrated DeepSeek AI chat functionality
- **Content Management**: Advertisements, friend links, website configuration
- **Data Statistics**: Article view count, like count statistics
- **Search Functionality**: Full-text search based on Elasticsearch

#### Installation

1. **Clone the project**
   ```bash
   git clone https://gitee.com/your-repo/go_blog.git
   cd go_blog
   ```

2. **Backend configuration**
   ```bash
   cd server
   # Install dependencies
   go mod tidy
   # Configure database
   cp configs/config.yaml.example configs/config.yaml
   # Modify database connection information in config file
   # Run the project
   go run main.go
   ```

3. **Frontend configuration**
   ```bash
   cd web
   # Install dependencies
   npm install
   # Run in development mode
   npm run dev
   # Build production version
   npm run build
   ```

#### Instructions

1. **Development environment startup**
   - Backend: `cd server && go run main.go`
   - Frontend: `cd web && npm run dev`

2. **Production environment deployment**
   - Backend: `cd server && go build -o main . && ./main`
   - Frontend: `cd web && npm run build`

3. **Database migration**
   - The project will automatically create database table structures

#### Project Highlights

- **Modern Architecture**: Adopts Go standard project structure with clear code organization
- **AI Integration**: Built-in AI chat assistant to enhance user experience
- **High Performance**: Uses Redis caching and Elasticsearch search
- **Responsive Design**: Frontend uses Vue3 + Element Plus
- **Complete Functionality**: Covers all core features of a blog system

#### Technical Highlights

1. **Go Standard Project Structure**: Follows Go project best practices
2. **Frontend-Backend Separation**: Clear API interface design
3. **AI Feature Integration**: Intelligent chat assistant
4. **High-Performance Architecture**: Caching + Search Engine
5. **Modern Frontend**: Vue3 + TypeScript + Element Plus
