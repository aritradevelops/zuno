# github.com/aritradevelops/zuno

github.com/aritradevelops/zuno is cli first, framework agnostic, code generator for golang web applications based on Hexagonal Architecture. github.com/aritradevelops/zuno not only generates the boiler plate but also generates the production ready CRUD.

# Supported Features
[x] Working CRUD Generation
[x] Multi Transport Support (HTTP for now)
[x] Multi Database Adapter Support (PostgreSQL with Bun And MongoDB for now)
[x] Internationalization
[x] Swagger Generation
[x] Validation
[x] Error Handling
[x] Logging
[x] Migration
[x] Hot Reloading with Air

# Why

github.com/aritradevelops/zuno provides the base for human or AI to build upon. github.com/aritradevelops/zuno is simple unlike other tools in the market

# Installation

```bash
# Clone the repository
git clone https://github.com/aritradevelops/github.com/aritradevelops/zuno.git

# Navigate to the project directory
cd github.com/aritradevelops/zuno

# Install dependencies
go mod tidy

# Run the application
go install
```

# Usage

```bash
# Initialize the project
github.com/aritradevelops/zuno init 

# Add a new module
github.com/aritradevelops/zuno add module <module-name-in-pascal>


# Add a new field to a module
github.com/aritradevelops/zuno add field <module-name-in-pascal>
```


# Supported Transports

[x] HTTP 
[ ] gRPC
[ ] WebSocket
[ ] GraphQL

# Supported Databases Adapters

[x] PostgreSQL with Bun
[x] MongoDB
[ ] MySQL

# Supported Migration Tools
[X] Goose

# Docker Support
[x] Docker Compose files for Dependencies.


# Contributions
Contributions are welcome! Please feel free to open an issue or submit a Pull Request.
