# synapsis-inventory-service

## Overview
This repository contains the Synapsis Inventory Service, a Go-based application for managing inventory. It uses PostgreSQL as the database and supports gRPC for communication.

## Prerequisites
- Docker
- Docker Compose

## How to Run

1. Clone the repository:
   ```bash
   git clone https://github.com/0x5w4/synapsis-inventory-service.git
   cd synapsis-inventory-service
   ```

2. Build and run the application using Docker Compose:
   ```bash
   docker-compose up --build
   ```

3. The application will be available at `http://localhost:8080`.

## Environment Variables
The application uses a `.env` file to configure environment variables. Below are the key variables:

- `APP_NAME`: Name of the application
- `APP_ENV`: Application environment (e.g., development, production)
- `APP_DEBUG`: Enable debug mode
- `APP_PORT`: Port the application runs on
- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name

## Database Migrations
Database migrations are located in the `migration/` directory. These are automatically applied when the application starts.

## Stopping the Application
To stop the application, run:
```bash
docker-compose down
```