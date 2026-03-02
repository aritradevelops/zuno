# Zuno

Zuno is cli first, framework agnostic, code generator for golang web applications based on Hexagonal Architecture. Zuno not only generates the boiler plate but also generates the production ready CRUD.

## Supported Features
- [x] Working CRUD Generation
- [x] Multi Transport Support (HTTP for now)
- [x] Multi Database Adapter Support (PostgreSQL with Bun And MongoDB for now)
- [x] Internationalization
- [x] Swagger Generation
- [x] Validation
- [x] Error Handling
- [x] Logging
- [x] Migration
- [x] Hot Reloading with Air

## Why

Zuno provides the base for human or AI to build upon. Zuno is simple unlike other tools in the market

## Installation

```bash
# Clone the repository
git clone https://github.com/aritradevelops/zuno.git

# Navigate to the project directory
cd zuno

# Install dependencies
go mod tidy

# Run the application
go install
```

## Usage

```bash
# Initialize the project
zuno init 

# Add a new module
zuno add module <module-name-in-pascal>

# Add a new field to a module
zuno add field <module-name-in-pascal>
```

## Supported Transports

* [x] HTTP
* [ ] gRPC
* [ ] WebSocket
* [ ] GraphQL

## Supported Databases Adapters

* [x] PostgreSQL with Bun
* [x] MongoDB
* [ ] MySQL

## Supported Migration Tools

* [x] Goose

## Docker Support

* [x] Docker Compose files for Dependencies.

## Contributions

Contributions are welcome! Please feel free to open an issue or submit a Pull Request.

