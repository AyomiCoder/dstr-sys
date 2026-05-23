# Distributed Event-Driven Notification System

This is a backend systems workshop disguised as a project.

Not a rush job.
Not a portfolio speedrun.
Not tutorial sludge.

This repo is where a notification system gets built slowly enough for the ideas to stick:

- HTTP servers
- clean package boundaries
- PostgreSQL
- Redis
- async workers
- retries
- dead-letter queues
- websockets
- observability
- scaling pressure
- distributed systems tradeoffs

## Why This Exists

The goal is not to finish fast.

The goal is to understand what serious backend systems are made of when you stop hiding behind abstractions:

- requests becoming work
- work becoming state
- state becoming events
- events becoming delivery
- failures becoming retries
- pressure becoming architecture

This starts as a simple Go service and grows, one hard-earned layer at a time.

## Current State

Day 1 foundation is in place:

- Go module setup
- thin application entrypoint
- config loading from environment
- basic HTTP router using `net/http`
- `GET /health` endpoint

## Project Layout

```text
cmd/
  api/          # application entrypoint
internal/
  config/       # environment-backed config
  http/         # HTTP router and handlers
configs/        # config examples and supporting files
```

## Getting Started

Requirements:

- Go 1.22+

Run the API:

```bash
go run ./cmd/api
```

Check the health endpoint:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{"status":"ok","service":"notification-service"}
```

## Configuration

Environment variables supported right now:

- `PORT` defaults to `8080`
- `SERVICE_NAME` defaults to `notification-service`

Example:

```bash
PORT=9090 SERVICE_NAME=dstr-sys go run ./cmd/api
```

## Build Philosophy

This codebase follows a few simple rules:

- do not over-engineer
- prefer clear boundaries over clever abstractions
- use standard error handling
- check existing structure before adding new structure
- avoid wasteful database loops when batching or better designs exist
- keep the code maintainable enough that future complexity does not turn it into junk

## What Comes Next

The early stages are intentionally plain.

That is the point.

Plain code is where you learn to see:

- where transport concerns end
- where business logic should begin
- where persistence belongs
- where concurrency helps
- where it lies to you

Today it is a health endpoint.
Later it becomes a system.
