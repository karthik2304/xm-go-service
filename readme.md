# CRUD - Company Microservices Golang (XM)

## Table of Contents
- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Project Setup](#project-setup)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Testing the API](#testing-the-api)

## Introduction
This is a **Golang-based microservice** that provides CRUD operations for handling company records. It includes JWT authentication, Kafka event handling, and MongoDB for data storage.

## Prerequisites
Ensure you have the following installed:

- **Go** (v1.23+)
- **Docker & Docker Compose latest Version**
- **MongoDB**
- **Kafka & Zookeeper**
- **Make** (for running Makefile commands)

## Project Setup

### Clone the Repository
```bash
git clone https://github.com/karthik2304/xm-go-service.git
cd xm-go-service
```

### Install Dependencies
```bash
go mod tidy
```

### Build OpenAPI for any Schema API Request changes
```bash
make openapi-build
```

### Run Docker Containers (MongoDB, Kafka, Zookeeper)
```bash
make run-docker
```

### Stop & Clean Docker Containers
```bash
make stop-docker
```

### Lint the Code
```bash
make lint
```

### Format the Code
```bash
make gofmt
```

### Run GO Application and exposed it http://localhost:9091 accessible outside of container:
```bash
make run-docker
```

## API Endpoints

### Authentication
#### 1. **Signup**
- **URL:** `POST /v1/auth/signup`
- **Payload:**
```json
{
  "username": "testuser",
  "password": "password123"
}
```
- **Response:**
```json
{
  "message": "User created successfully"
}
```

#### 2. **Login**
- **URL:** `POST /v1/auth/login`
- **Payload:**
```json
{
  "username": "testuser",
  "password": "password123"
}
```
- **Response:**
```json
{
  "jwt_token": "eyJhbGciOiJI..."
}
```

### Company Management (Requires Bearer Token)

#### 3. **Create Company**
- **URL:** `POST /v1/create-company`
- **Headers:**
  - Authorization: `Bearer <JWT_TOKEN>`
- **Payload:**
```json
{
  "companyUUID": "12345",
  "companyName": "TechCorp",
  "description": "A tech company",
  "totalEmployees": 100,
  "registered": true,
  "type": "Corporations"
}
```
- **Response:**
```json
{
  "message": "Company created successfully"
}
```

#### 4. **List Companies**
- **URL:** `GET /v1/list-companies`
- **Headers:**
  - Authorization: `Bearer <JWT_TOKEN>`
- **Response:**
```json
[
    {
      "companyUUID": "12345",
      "companyName": "TechCorp",
      "description": "A tech company",
      "totalEmployees": 100,
      "registered": true,
      "type": "Corporations"
    },{
        ...
    }
  
]
```

#### 5. **Get Company Details**
- **URL:** `GET /v1/company-details/{company_uuid}`
- **Headers:**
  - Authorization: `Bearer <JWT_TOKEN>`
- **Response:**
```json
{
  "companyUUID": "12345",
  "companyName": "TechCorp",
  "description": "A tech company",
  "totalEmployees": 100,
  "registered": true,
  "type": "Corporations"
}
```

#### 6. **Update Company Details**
- **URL:** `PATCH /v1/company-details/{company_uuid}`
- **Headers:**
  - Authorization: `Bearer <JWT_TOKEN>`
- **Payload:**
```json
{
  "companyName": "Updated TechCorp",
  "description": "Updated description",
  "totalEmployees": 150,
  "registered": true,
  "type": "Corporations"
}
```
- **Response:**
```json
{
  "message": "Company updated successfully"
}
```

#### 7. **Delete Company**
- **URL:** `DELETE /v1/company-details/{company_uuid}`
- **Headers:**
  - Authorization: `Bearer <JWT_TOKEN>`
- **Response:**
```json
{
  "message": "Company deleted successfully"
}
```

## Testing the API

### Run Unit Tests
```bash
make test
```

### Using cURL for API Calls

Example 1:

```bash
curl -X POST "http://localhost:9091/v1/auth/login" \
     -H "Content-Type: application/json" \
     -d '{ "username": "xmtest@gmail.com", "password": "xmtest" }'

response: {
"jwt_token": "jsxmx.."
}
```

Example 2 for JWT Token Auth:
```bash
curl -X POST "http://localhost:9091/v1/create-company" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <JWT_TOKEN>" \
     -d '{
           "companyUUID": "12345",
           "companyName": "TechCorp",
           "description": "A tech company",
           "totalEmployees": 100,
           "registered": true,
           "type": "Corporations"
         }'
```


### For Testing via Postman
1. Import the openapi-postman-collection.json document in collection
2. It will Import API test routes for testing
3. Include JWT tokens for protected routes inside {{bearerToken}} for Authorization.

---

NOTE:
1) by default 1 login user & 2 companies data created automatically by bootstrap script

This document should help you **set up, run, and test** the microservice easily! 🚀




