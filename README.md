# HireFlow API

A scalable backend service for job posting and job application management built with Go.

## Features

- JWT Authentication
- Refresh Tokens
- Role_Based Authorization
- Background worker pool
- PostgreSQL Database
- Database Migrations
- Dockerized Infrastructure
- Layered Architecture
- Concurrency Handling
- Transaction Support

---

## Tech Stack

- Go
- PostgreSQL
- Docker
- Chi Router
- JWT
- golang-migrate

---

## Architecture

Handler -> Service -> Store -> PostgreSQL

---

## Run Locally
go run ./cmd/api/main.go

---

## Deployment

git clone ...
docker compose up -d --build

---

### Start with Docker

```bash
docker compose up --build

