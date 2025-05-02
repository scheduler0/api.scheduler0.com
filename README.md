# Scheduler0 API Documentation Server

A lightweight HTTP server that serves API documentation for the Scheduler0 project. This server is designed to host and serve OpenAPI documentation files.

## Features

- Serves OpenAPI documentation files
- Built with Go using the Gorilla Mux router
- HTTP request logging
- Simple and lightweight implementation

## Prerequisites

- Go 1.23.7 or later
- Git

## Installation

1. Clone the repository:
```bash
git clone https://github.com/scheduler0/api.scheduler0.com
cd api.scheduler0.com
```

2. Install dependencies:
```bash
go mod download
```

## Usage

1. Start the server:
```bash
go run main.go
```

The server will start on port 3002 and serve the API documentation from the `api-docs` directory.

## API Documentation

The API documentation is served at:
```
http://localhost:3002/api-docs/
```

## Dependencies

- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router and dispatcher
- [go-http-utils/logger](https://github.com/go-http-utils/logger) - HTTP request logging middleware

## License

This project is licensed under the terms specified in the [LICENSE](LICENSE) file.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 