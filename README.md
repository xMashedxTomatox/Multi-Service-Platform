# Feedback Platform – Multi-Service Backend

A backend platform showcasing a **scalable multi-service architecture** with user authentication and feedback services.  
Designed to be **extensible**, allowing new services to be added alongside existing ones as the system grows.

---

## 🚀 Features
- **Auth Service** – User signup/login with JWT authentication  
- **Feedback Service** – Allows authenticated users to submit feedback  
- **Redis Caching** – Fast session and cross-service data sharing  
- **PostgreSQL** – Persistent datastore for auth and feedback  
- **Kubernetes (EKS)** – Service orchestration, scaling, and deployment  
- **Ingress (ALB)** – Unified access through `/auth` and `/feedback` routes  
- **CI/CD (planned)** – Automated builds & deployments with GitHub Actions → Amazon ECR → EKS  

---

## 🏗️ Architecture

```mermaid
flowchart TB
  user(User / Client)
  alb(AWS ALB Ingress)

  user --> alb

  alb --> auth(Auth Service - Go)
  alb --> fb(Feedback Service - Go)

  auth --> redis[Redis Cache]
  fb --> redis

  auth --> authpg[(Postgres - Auth)]
  fb --> fbpg[(Postgres - Feedback)]
