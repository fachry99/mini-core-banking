# Mini Core Banking API

## Overview

Mini Core Banking API is a lightweight backend service that simulates core banking functionalities, with a strong focus on **data consistency, concurrency safety, and backend engineering best practices**.

This project is designed as a **backend engineering portfolio**, demonstrating how financial transaction systems are typically handled in real-world environments.

### Key Features

- User Management
- Account Management
- Deposit
- Fund Transfer between accounts
  - Atomic (all-or-nothing)
  - Race-condition safe
  - Idempotent

---

## Architecture

### High-Level Architecture

[ Client / API Consumer ]
            |
            v
      HTTP Handler Layer
            |
            v
      Service (Business Logic)
            |
            v
      Repository (Database Access)
            |
            v
        PostgreSQL

Each layer has a clear responsibility:
- Handler: HTTP request & response handling
- Service: business rules and transaction orchestration
- Repository: database access logic

---

## Transfer Flow (Atomic & Concurrent-Safe)

Client
  |
  | POST /transfer (Idempotency-Key)
  v
TransferHandler
  - Request validation
  - Idempotency check
  v
TransferService
  - Business rule validation
  - Begin database transaction
  - Lock accounts (SELECT ... FOR UPDATE)
  - Update balances
  - Insert transaction record
  - Commit transaction
  - Write audit log
  v
PostgreSQL

---

## Concurrency & Race Condition Handling

### Problem

Without proper locking, concurrent transfers may cause:
- Lost updates
- Negative balances
- Inconsistent account data

### Solution

- Row-level locking using `SELECT ... FOR UPDATE`
- Consistent lock ordering based on account ID
- Single database transaction per transfer

Example execution order:

Tx A: Lock Account A → Lock Account B → Transfer  
Tx B: Wait Account A → Lock Account A → Lock Account B → Transfer  

Result:
- No race condition
- No deadlock
- Consistent balances

---

## Idempotent Transfer Handling

To prevent duplicate transfers caused by retries (network issues, timeouts, etc.):

- Client must send an `Idempotency-Key` header
- Transfer response is stored
- Requests with the same key return the previously stored response

Example:

Idempotency-Key: abc-123

---

## Audit Logging

Every transfer generates an audit log containing:
- request_id
- from_account_id
- to_account_id
- amount
- status (SUCCESS / FAILED)

Audit logging is executed exactly once using `defer` in the service layer, ensuring it runs even if an error occurs.

---

## System Guarantees

This system guarantees:
- Atomic fund transfers (all-or-nothing)
- No double spending under concurrent requests
- Idempotent transfer execution
- Consistent account balances
- Auditable transaction history

---

## API Example

### Transfer Funds

Request:
POST /transfer  
Content-Type: application/json  
Idempotency-Key: abc-123  

{
  "from_account_id": 1,
  "to_account_id": 2,
  "amount": 50000
}

Response:
{
  "transaction_id": "tx_20240101_001",
  "status": "SUCCESS"
}

---

## Project Structure

cmd/api/main.go  
internal/  
 ├── handler/      // HTTP handlers  
 ├── service/      // Business logic & transactions  
 ├── repository/   // Database access  
 ├── middleware/   // Request ID, CORS, etc  
 ├── audit/        // Audit logging  
 ├── dto/          // Request/response models  
 └── config/       // Application & database configuration  

---

## Error Handling Strategy

Handler layer  
- Responsible for HTTP status codes and response formatting

Service layer  
- Handles business rules and domain-level errors

Repository layer  
- Handles database-related errors

---

## How to Run Locally

### Prerequisites

- Go 1.20 or newer
- PostgreSQL
- Properly configured environment variables for database connection

### Run the Application

go run cmd/api/main.go

The server will start on:

http://localhost:8080

The API can be tested using:
- Postman
- cURL
- Frontend applications
- Any API testing tool

---

## Production Readiness

Implemented:
- Transaction safety
- Concurrency control
- Idempotency
- Audit logging
- Request ID middleware
- Clean architecture layering

Planned improvements:
- Authentication & authorization (JWT)
- Rate limiting
- Observability (metrics & tracing)
- CI/CD pipeline

---

## Author

Fachry  
Backend Developer
