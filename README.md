# XM Microservice - Golang - Mrudhul Guda

## Overview

The **XM Microservice**, developed in **Golang**, handles company-related operations, including the creation, updating, deletion, and retrieval of company data. It ensures secure access through **JWT authentication**, leverages **PostgreSQL** for robust data management, and utilizes **Kafka** for efficient event-driven communication. The service is containerized using **Docker** for easy deployment.

## Features

- **Comprehensive CRUD Operations** for company data
- **JWT-Based Authentication** for secure access
- **PostgreSQL Database Integration** for reliable data storage
- **Kafka Integration** for event-driven architecture
- **Dockerized Environment** for seamless deployment

## Prerequisites

- **Docker & Docker Compose**
- Ensure the following ports are available:
  - `5432` (PostgreSQL)
  - `9092` (Kafka)
  - `8080` (Application Port)

## Environment Variables

Create a `.env` file in the root directory with the following content:

```env
PORT=8080
DATABASE_URL=postgresql://user:password@db:5432/xmdb?sslmode=disable
JWT_SECRET=xmsecretkey
KAFKA_BROKER=kafka:9092
KAFKA_PARTITIONS=3
KAFKA_REPLICATION_FACTOR=1
KAFKA_TOPIC_COMPANY=company-events
```

## Setup and Running the Service

1. **Build and Start the Services:**
   ```bash
   docker-compose up --build
   ```

2. **Verify Service Health:**
   ```bash
   curl http://localhost:8080/health
   ```
   **Expected Response:**
   ```json
   {
     "status": "ok"
   }
   ```

## API Endpoints

### **1. User Registration**

Registers a new user.

```bash
curl --location 'http://localhost:8080/api/users' \
--header 'Content-Type: application/json' \
--data '{
    "username":"admin",
    "password":"password"
}'
```

### **2. User Login (JWT Token Generation)**

Generates a JWT token for authenticated requests.

```bash
curl --location 'http://localhost:8080/api/login' \
--header 'Content-Type: application/json' \
--data '{
    "username": "admin",
    "password": "password"
}'
```

**Response:**
```json
{
  "created_at": 1739111062,
  "expires_at": 1739118262,
  "token": "<JWT_TOKEN>"
}
```

### **3. Create a Company (Authenticated)**

**Request:**
```bash
curl --location 'http://localhost:8080/api/companies' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <JWT_TOKEN>' \
--data '{
    "name": "Corp",
    "amount_of_employees": 150,
    "registered": true,
    "type": "Corporation",
    "description": "A leading technology company"
}'
```

**Constraints:**
- **Name:** Max 15 characters, unique
- **Description:** Optional, up to 3000 characters
- **Type:** One of `Corporation`, `NonProfit`, `Cooperative`, `Sole Proprietorship`
- **Amount of Employees:** Required, integer
- **Registered:** Required, boolean

**Response:**
```json
{
  "id": "e3f1a8b2-9d14-4c2b-8c3f-1a2f3d4e5678",
  "name": "Corp",
  "description": "A leading technology company",
  "amount_of_employees": 150,
  "registered": true,
  "type": "Corporation"
}
```

### **4. Get Company Details (Public Access)**

**Request:**
```bash
curl --location 'http://localhost:8080/api/companies/{id}'
```

**Example:**
```bash
curl --location 'http://localhost:8080/api/companies/e3f1a8b2-9d14-4c2b-8c3f-1a2f3d4e5678'
```

**Response:**
```json
{
  "id": "e3f1a8b2-9d14-4c2b-8c3f-1a2f3d4e5678",
  "name": "Corp",
  "description": "A leading technology company",
  "amount_of_employees": 150,
  "registered": true,
  "type": "Corporation"
}
```

### **5. Update Company Details (Authenticated)**

**Request Format:**
```bash
curl --location 'http://localhost:8080/api/companies/{id}' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <JWT_TOKEN>' \
--request PATCH \
--data '{
    "name": "Corp",
    "amount_of_employees": 200,
    "registered": false,
    "type": "NonProfit",
    "description": "Updated description for the company."
}'
```

**Example with ID:**
```bash
curl --location 'http://localhost:8080/api/companies/e3f1a8b2-9d14-4c2b-8c3f-1a2f3d4e5678' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <JWT_TOKEN>' \
--request PATCH \
--data '{
    "name": "Corp",
    "amount_of_employees": 200,
    "registered": false,
    "type": "NonProfit",
    "description": "Updated description for the company."
}'
```

**Response:**
```json
{
  "id": "e3f1a8b2-9d14-4c2b-8c3f-1a2f3d4e5678",
  "name": "Corp",
  "description": "Updated description for the company.",
  "amount_of_employees": 200,
  "registered": false,
  "type": "NonProfit"
}
```

### **6. Delete a Company (Authenticated)**

**Request Format:**
```bash
curl --location --request DELETE 'http://localhost:8080/api/companies/{id}' \
--header 'Authorization: Bearer <JWT_TOKEN>'
```

**Example with ID:**
```bash
curl --location --request DELETE 'http://localhost:8080/api/companies/e3f1a8b2-9d14-4c2b-8c3f-1a2f3d4e5678' \
--header 'Authorization: Bearer <JWT_TOKEN>'
```

**Response:**
- **Status Code:** `204 No Content`
- No content is returned in the response body.


## Kafka Consumer for Company Events

To consume events from the `company-events` Kafka topic, use the following command:

```bash
docker exec -it xm-microservice-kafka-1 kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic company-events \
  --from-beginning
```

This command will allow you to view all events related to company operations (create, update, delete) from the beginning of the topic.


## Authorization Header Format

For all authenticated requests, include the following header:

```bash
Authorization: Bearer <JWT_TOKEN>
```

## Troubleshooting

- **Verify Docker Containers:**
  ```bash
  docker ps
  ```

- **Review Service Logs:**
  ```bash
  docker-compose logs
  ```

- **Database Connection Issues:** Ensure PostgreSQL is running and accessible with the credentials provided in `.env`.

- **Kafka Issues:** Check Kafka broker logs for errors and ensure it is running on the configured port (`9092`).

