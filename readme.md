# Mini Core Banking API

## 📌 Overview

Mini Core Banking adalah backend service sederhana yang mensimulasikan sistem perbankan inti (core banking) dengan fitur utama:

* User Management
* Account Management
* Deposit
* Transfer antar akun (atomic, race-condition safe)

Project ini difokuskan pada **backend engineering best practices**: transaction safety, concurrency control, idempotency, audit logging, dan clean architecture.

---

## 🏗️ Architecture

### High-Level Architecture

```
[ Client / API Consumer ]
            |
            v
      HTTP Handler Layer
            |
            v
      Service (Business Logic)
            |
            v
      Repository (DB Access)
            |
            v
        PostgreSQL
```

---

### Transfer Flow (Phase 9)

```
Client
  |
  | POST /transfer (Idempotency-Key)
  v
TransferHandler
  |  - Validate request
  |  - Idempotency check
  v
TransferService
  |  - Validate business rule
  |  - Begin TX
  |  - Lock accounts (FOR UPDATE)
  |  - Update balances
  |  - Insert transaction log
  |  - Commit
  |  - Audit log
  v
PostgreSQL
```

---

## 🔐 Concurrency & Race Condition Handling

### Problem

Tanpa locking, concurrent transfer bisa menyebabkan:

* Lost update
* Negative balance
* Inconsistent data

### Solution

* `SELECT ... FOR UPDATE`
* Consistent lock ordering (sorted account ID)
* Single database transaction

```
Tx A: Lock A -> Lock B -> Transfer
Tx B: Wait A -> Lock A -> Lock B -> Transfer
```

Result: **No race condition, no deadlock**

---

## 🔁 Idempotency (Phase 8)

Untuk mencegah double transfer akibat retry:

* Client wajib mengirim `Idempotency-Key`
* Response disimpan
* Request dengan key sama akan return response sebelumnya

```
Idempotency-Key: abc-123
```

---

## 🧾 Audit Logging (Phase 9)

Setiap transfer akan menghasilkan audit log:

* request_id
* from_account_id
* to_account_id
* amount
* status (SUCCESS / FAILED)

Audit dipanggil **1 kali** menggunakan `defer` di service.

---

## 📂 Project Structure

```
cmd/api/main.go
internal/
 ├── handler/
 ├── service/
 ├── repository/
 ├── middleware/
 ├── audit/
 ├── dto/
 └── config/
```

---

## 🧪 Error Handling Strategy

| Layer      | Responsibility               |
| ---------- | ---------------------------- |
| Handler    | HTTP status & response       |
| Service    | Business rule & domain error |
| Repository | DB error                     |

---

## 🚀 Production Readiness

Implemented:

* Transaction safety
* Concurrency control
* Idempotency
* Audit logging
* Request ID middleware

Future improvements:

* Authentication (JWT)
* Rate limiting
* Observability (metrics & tracing)

---

## 🎯 Purpose

Project ini dirancang sebagai:

* Backend portfolio
* Interview showcase (mid–senior ready)
* Learning reference untuk core banking system

---

## Author

Fachry — Backend Developer
