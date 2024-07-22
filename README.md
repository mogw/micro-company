# Company Microservice

This is a microservice for managing company information, built with Golang, Gin, MongoDB, Kafka, and Docker. It supports creating, updating, deleting, and retrieving company information. The service uses JWT for authentication and produces Kafka events for each mutating operation.

## Table of Contents

- [Company Microservice](#company-microservice)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Setup](#setup)
    - [Prerequisites](#prerequisites)
    - [Running with Docker Compose](#running-with-docker-compose)
    - [Running Locally](#running-locally)
  - [Running the Tests](#running-the-tests)
  - [Linting the Code](#linting-the-code)
  - [API Endpoints](#api-endpoints)

## Features

- Create, update, delete, and retrieve companies.
- JWT-based authentication for protected routes.
- Kafka integration for producing events on mutating operations.
- Integration tests with database cleanup.
- Linting with `golangci-lint`.

## Configuration

The application configuration is stored in `config.yaml`. Here is an example:

```yaml
mongo_uri: "mongodb://root:example@mongodb:27017"
kafka_broker: "kafka:9092"
jwt_secret: "my_secret_key"
```

## Setup

### Prerequisites

- Docker
- Docker Compose
- Go (if running locally)
- MongoDB (if running locally)
- Kafka (if running locally)

### Running with Docker Compose

1. Clone the repository:

    ```sh
    git clone https://github.com/mogw/micro-company.git
    cd micro-company
    ```

2. Build and run the services:

    ```sh
    docker-compose up --build
    ```

    This will start the application along with MongoDB and Kafka services.

3. Access the application:

    The application will be available at http://localhost:8080.

### Running Locally

1. Clone the repository:

    ```sh
    git clone https://github.com/mogw/micro-company.git
    cd micro-company
    ```

2. Start MongoDB and Kafka:

    Ensure MongoDB and Kafka are running. You can use Docker for this:

    ```sh
    docker run -d -p 27017:27017 --name mongodb mongo
    docker run -d -p 2181:2181 -p 9092:9092 --name zookeeper wurstmeister/zookeeper
    docker run -d -p 9092:9092 --name kafka --env KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 --env KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 --link zookeeper wurstmeister/kafka
    ```

3. Run the applicaiton:

    ```sh
    go mod tidy
    go run cmd/main.go
    ```

    The application will be available at http://localhost:8080.


## Running the Tests

To run the integration tests, use the following command:

```sh
go test ./...
```

This will run all tests in the project.


## Linting the Code

To lint the code using golangci-lint, use the following command:

1. Install golangci-lint:

    ```sh
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    ```

2. Run the linter:

    ```sh
    golangci-lint run
    ```

## API Endpoints

### Public Endpoints

- GET `/companies/`: Retrieve a company by ID.

### Protected Endpoints

- POST `/companies`: Create a new company.
- PATCH `/companies/{id}`: Update an existing company.
- DELETE `/companies/{id}`: Delete a company.

### Example Requests

- Create a Company

    ```sh
    curl -X POST http://localhost:8080/companies \
        -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Example Company",
            "description": "This is an example company.",
            "amount_of_employees": 50,
            "registered": true,
            "type": "Corporations"
            }'

    ```


- Get a Company

    ```sh
    curl http://localhost:8080/companies/<COMPANY_ID>
    ```

- Update a Company

    ```sh
    curl -X PATCH http://localhost:8080/companies/<COMPANY_ID> \
        -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Updated Company",
            "amount_of_employees": 100
            }'

    ```

- Delete a Company

    ```sh
    curl -X DELETE http://localhost:8080/companies/<COMPANY_ID> \
        -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
    ```

Replace <YOUR_JWT_TOKEN> with a valid JWT token and <COMPANY_ID> with the ID of the company you want to operate on.