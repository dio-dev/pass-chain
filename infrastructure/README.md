---
hidden: true
---

# Infrastructure - Terraform & Kubernetes

Infrastructure as Code for Pass Chain using Terraform and Kubernetes.

## Overview

This directory contains all infrastructure configuration for deploying Pass Chain to production.

## Architecture

```
AWS Cloud
├── VPC
│   ├── Public Subnets (NAT, Load Balancers)
│   └── Private Subnets (EKS, RDS, Redis)
├── EKS Cluster
│   ├── Backend Pods
│   ├── Frontend Pods
│   ├── Vault Pods
│   └── Blockchain Pods
├── RDS PostgreSQL
├── ElastiCache Redis
├── S3 (Backups, Assets)
└── CloudFront (CDN)
```

## Directory Structure

```
infrastructure/
├── terraform/           # Terraform configurations
│   ├── main.tf         # Main configuration
│   ├── variables.tf    # Input variables
│   ├── outputs.tf      # Output values
│   ├── modules/        # Reusable modules
│   │   ├── vpc/        # VPC module
│   │   ├── eks/        # EKS cluster module
│   │   ├── rds/        # RDS module
│   │   └── vault/      # Vault module
│   └── environments/   # Environment configs
│       ├── dev/
│       ├── staging/
│       └── production/
├── k8s/                # Kubernetes manifests
│   ├── backend/        # Backend deployment
│   ├── frontend/       # Frontend deployment
│   ├── vault/          # Vault deployment
│   └── monitoring/     # Prometheus, Grafana
└── docker/             # Docker Compose for local dev
    └── docker-compose.yml
```

## Prerequisites

* AWS CLI configured
* Terraform 1.6+
* kubectl
* Helm 3+
* Docker

## Getting Started

### 1. Configure AWS Credentials

```bash
aws configure
```

### 2. Initialize Terraform

```bash
cd infrastructure/terraform
terraform init
```

### 3. Review Plan

```bash
terraform plan
```

### 4. Apply Infrastructure

```bash
terraform apply
```

This will create:

* VPC with public/private subnets
* EKS cluster
* RDS PostgreSQL
* ElastiCache Redis
* HashiCorp Vault
* All necessary security groups and IAM roles

### 5. Configure kubectl

```bash
aws eks update-kubeconfig --region us-east-1 --name passchain-eks
```

### 6. Deploy Applications

```bash
# Deploy backend
kubectl apply -f k8s/backend/

# Deploy frontend
kubectl apply -f k8s/frontend/

# Deploy Vault
helm install vault hashicorp/vault -f k8s/vault/values.yaml
```

## Terraform Modules

### VPC Module

Creates VPC with:

* 3 public subnets (for load balancers)
* 3 private subnets (for workloads)
* Internet Gateway
* NAT Gateways
* Route tables

### EKS Module

Creates EKS cluster with:

* Managed node groups
* Auto-scaling configured
* IRSA for pod IAM roles
* VPC CNI plugin
* CoreDNS
* kube-proxy

### RDS Module

Creates RDS PostgreSQL with:

* Multi-AZ deployment
* Automated backups
* Encryption at rest
* Security group rules

### Vault Module

Deploys HashiCorp Vault with:

* HA mode (3 replicas)
* Auto-unseal with AWS KMS
* Integrated storage (Raft)
* TLS certificates

## Kubernetes Deployments

### Backend Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: passchain-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: passchain-backend
  template:
    spec:
      containers:
      - name: backend
        image: passchain/backend:latest
        ports:
        - containerPort: 8080
```

### Frontend Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: passchain-frontend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: passchain-frontend
  template:
    spec:
      containers:
      - name: frontend
        image: passchain/frontend:latest
        ports:
        - containerPort: 3000
```

## Local Development

Use Docker Compose for local testing:

```bash
cd infrastructure/docker
docker-compose up
```

Services available:

* Backend: http://localhost:8080
* Frontend: http://localhost:3000
* PostgreSQL: localhost:5432
* Redis: localhost:6379
* Vault: http://localhost:8200

## Secrets Management

All secrets are stored in HashiCorp Vault:

```bash
# Login to Vault
vault login

# Write secret
vault kv put secret/passchain/db username=admin password=secret

# Read secret
vault kv get secret/passchain/db
```

Kubernetes pods use Vault Agent Injector to fetch secrets.

## Monitoring

### Prometheus & Grafana

```bash
# Install Prometheus
helm install prometheus prometheus-community/prometheus

# Install Grafana
helm install grafana grafana/grafana
```

Access Grafana:

```bash
kubectl port-forward svc/grafana 3000:80
```

### CloudWatch

All logs are sent to CloudWatch Log Groups:

* `/aws/eks/passchain/backend`
* `/aws/eks/passchain/frontend`
* `/aws/rds/passchain`

## Scaling

### Horizontal Pod Autoscaling

```bash
kubectl autoscale deployment passchain-backend --cpu-percent=70 --min=3 --max=10
```

### Cluster Autoscaling

EKS node groups auto-scale based on pod resource requests.

## Backup & Disaster Recovery

### RDS Backups

* Automated daily backups (7-day retention)
* Manual snapshots before major changes

### Vault Backups

```bash
# Take snapshot
vault operator raft snapshot save backup.snap

# Restore snapshot
vault operator raft snapshot restore backup.snap
```

### Blockchain Backups

Ledger data backed up to S3 daily.

## Security

### Network Security

* Private subnets for all workloads
* Security groups with minimal access
* WAF on CloudFront/ALB

### Encryption

* RDS encryption at rest
* S3 bucket encryption
* EKS secrets encryption with KMS
* TLS for all traffic

### IAM & RBAC

* Least privilege IAM roles
* Kubernetes RBAC configured
* Pod security policies

## CI/CD

GitHub Actions workflows deploy to environments:

```yaml
# .github/workflows/deploy.yml
- name: Deploy to EKS
  run: |
    kubectl set image deployment/passchain-backend backend=$IMAGE_TAG
    kubectl rollout status deployment/passchain-backend
```

## Cost Optimization

* Use Spot instances for non-critical workloads
* Auto-scaling to match demand
* S3 lifecycle policies for old backups
* Reserved instances for production

## Troubleshooting

### EKS Connection Issues

```bash
# Update kubeconfig
aws eks update-kubeconfig --name passchain-eks

# Check nodes
kubectl get nodes

# Check pods
kubectl get pods --all-namespaces
```

### RDS Connection Issues

```bash
# Test connection
psql -h <rds-endpoint> -U passchain -d passchain

# Check security groups
aws ec2 describe-security-groups --group-ids <sg-id>
```

### Vault Issues

```bash
# Check status
kubectl exec vault-0 -- vault status

# Unseal manually
kubectl exec vault-0 -- vault operator unseal <key>
```

## Destroy Infrastructure

⚠️ **Warning**: This will delete everything!

```bash
terraform destroy
```

## Production Checklist

* [ ] Change all default passwords
* [ ] Enable MFA on AWS account
* [ ] Set up CloudTrail logging
* [ ] Configure backups
* [ ] Set up monitoring alerts
* [ ] Review security groups
* [ ] Enable GuardDuty
* [ ] Set up WAF rules
* [ ] Configure auto-scaling
* [ ] Test disaster recovery

## License

MIT
