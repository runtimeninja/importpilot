Shopify Product Import SaaS — System Design Scope

1) Actors

1.1 Super Admin
The Super Admin controls the entire SaaS platform.

Capabilities:
- Login
- View all client accounts
- Create new clients
- Activate/deactivate clients
- View system-level scraping logs
- Monitor failed jobs
- View overall usage

1.2 Client User
A Client User is a Shopify store owner or team member.

Capabilities:
- Login
- Password reset
- Manage profile
- Connect Shopify store
- Add scrape source URLs
- Run manual scrape jobs
- Create/update schedules
- View scraping history
- Download CSV exports
- View system notifications

1.3 System
Non-human actor responsible for internal operations.

Responsibilities:
- Trigger scheduled jobs
- Process scraping queue
- Handle retries
- Sync with Shopify
- Save logs
- Generate notifications


2) Core Modules

2.1 Authentication & Authorization
- Login
- Password hashing
- JWT-based authentication
- Role-Based Access Control (RBAC)
- Password reset

2.2 Tenant / Client Management
- Create client accounts
- Store shop information
- Activate/deactivate clients
- Manage usage limits
- Plan-based separation

2.3 Shopify Connection
- Connect Shopify store
- Store credentials/tokens
- Manage sync metadata

2.4 Scraper Management
- Add source URLs
- Save scraping config
- Trigger manual scraping
- Parse and store results

2.5 CSV Export Engine
- Map scraped data to Shopify CSV
- Generate CSV files
- Maintain export history

2.6 Scheduler & Job Queue
- Run jobs every 24 hours
- Retry on failure
- Prevent duplicate jobs
- Handle background processing

2.7 Shopify Product Sync
- Create/update Shopify products
- Maintain sync logs
- Track errors

2.8 Logs & Notifications
- Scraping history
- Sync history
- Failure tracking
- User/admin notifications

2.9 Billing-Ready Foundation
- Plans
- Usage quota
- Feature flags
- Ready for future Stripe integration


3) Out of Scope (MVP)
- Multi-language UI
- Advanced analytics dashboard
- Team collaboration features
- Full billing integration
- Real-time webhook sync
- AI-based product enrichment
- Multi-shop per client
- Visual scraping builder
- Public GraphQL API


4) External Systems
- Shopify Admin API
- Email provider (SendGrid / Resend / SES)
- PostgreSQL
- Redis
- Object storage (S3-compatible)


5) Non-Functional Requirements

Performance:
- Fast API response for manual requests
- Heavy tasks run in background
- Pagination for lists

Scalability:
- Handle thousands of jobs
- Horizontal worker scaling
- Multi-tenant isolation

Security:
- Hashed passwords
- Encrypted secrets
- Rate limiting
- RBAC
- Audit logs

Reliability:
- Retry mechanism
- Inspect failed jobs
- Idempotent design
- Prevent duplicate execution

Maintainability:
- Modular architecture
- Separation of concerns
- Microservice-ready
- Testable code


6) Tech Stack

Backend:
- Golang

Frontend:
- Angular

Database:
- PostgreSQL

Cache/Queue:
- Redis

Job Processing:
- Asynq

API:
- REST

Auth:
- JWT + Refresh Token

Scraping:
- HTTP + HTML parsing
- JSON-LD extraction
- Optional headless browser
- AI fallback (future)

Deployment:
- Docker
- Nginx
- GitHub Actions
- VPS / AWS / DigitalOcean


7) Architecture

Components:
- API Server
- Worker Service
- Scheduler (Cron)
- Redis Queue
- PostgreSQL

Manual Flow:
Client → API → Queue → Worker → Scrape → Save → CSV → Done

Scheduled Flow:
Cron → Fetch schedules → Queue → Worker → Process

Shopify Sync:
Worker → Shopify API → Update products → Log result

Rules:
- API stays thin
- Everything async
- Idempotent jobs
- Retry safe
- Multi-tenant isolation
- Observability ready


8) Database Schema

Core Tables:
- clients
- users
- roles
- user_roles

Business Tables:
- shopify_connections
- scrape_sources
- scrape_schedules
- scrape_jobs
- csv_exports
- shopify_sync_logs

Operational Tables:
- notifications
- usage_counters
- audit_logs


9) Data Integrity Rules

- A user belongs to one client
- Scrape job must match client of source
- Schedule must match client of source
- Sync logs must match job and connection client
- Inactive clients cannot create jobs
- Inactive users cannot act
