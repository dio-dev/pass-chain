# GitHub Actions CI/CD for GKE

Complete CI/CD pipeline for deploying Pass Chain to Google Kubernetes Engine (GKE).

## Workflows Created

### 1. `deploy-backend.yml`
- Builds Go backend Docker image
- Pushes to Google Container Registry (GCR)
- Deploys to GKE with Helm
- Runs health checks

### 2. `deploy-frontend.yml`
- Builds Next.js frontend Docker image
- Pushes to GCR
- Deploys to GKE with Helm
- Gets LoadBalancer IP

### 3. `deploy-infrastructure.yml`
- Deploys PostgreSQL (Bitnami Helm chart)
- Deploys Redis (Bitnami Helm chart)
- Deploys Vault (HashiCorp Helm chart)
- Configures Ingress controller

### 4. `deploy-full.yml`
- Complete pipeline: Test → Deploy Infra → Deploy Apps → Verify
- Runs tests first
- Orchestrates all deployments
- Health checks everything

---

## Setup Instructions

### 1. Create GKE Cluster

```bash
# Set variables
export PROJECT_ID="your-gcp-project"
export CLUSTER_NAME="pass-chain-cluster"
export REGION="us-central1"

# Create cluster
gcloud container clusters create $CLUSTER_NAME \
  --project=$PROJECT_ID \
  --zone=$REGION-a \
  --num-nodes=3 \
  --machine-type=n1-standard-2 \
  --enable-autoscaling \
  --min-nodes=3 \
  --max-nodes=10 \
  --enable-autorepair \
  --enable-autoupgrade \
  --disk-size=50GB \
  --disk-type=pd-standard

# Get credentials
gcloud container clusters get-credentials $CLUSTER_NAME --zone=$REGION-a
```

### 2. Create Service Account

```bash
# Create service account for GitHub Actions
gcloud iam service-accounts create github-actions \
  --display-name="GitHub Actions"

# Grant permissions
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/container.developer"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/storage.admin"

# Create key
gcloud iam service-accounts keys create key.json \
  --iam-account=github-actions@$PROJECT_ID.iam.gserviceaccount.com

# Convert to base64 for GitHub secret
cat key.json | base64 > key-base64.txt
```

### 3. Configure GitHub Secrets

Go to your repo → Settings → Secrets and add:

| Secret Name | Value | Description |
|-------------|-------|-------------|
| `GKE_PROJECT` | `your-project-id` | GCP Project ID |
| `GKE_SA_KEY` | `<contents of key-base64.txt>` | Service account key (base64) |
| `DB_PASSWORD` | `<strong-password>` | PostgreSQL password |
| `REDIS_PASSWORD` | `<strong-password>` | Redis password |
| `API_URL` | `http://api.passchain.io` | Backend API URL |

### 4. Update values-gke.yaml

Edit `infrastructure/k8s/helm/passchain/values-gke.yaml`:

```yaml
frontend:
  image:
    repository: gcr.io/YOUR_PROJECT_ID/pass-chain-frontend  # Change this

backend:
  image:
    repository: gcr.io/YOUR_PROJECT_ID/pass-chain-backend  # Change this

ingress:
  hosts:
    - host: app.passchain.io  # Your domain
    - host: api.passchain.io  # Your domain
```

### 5. Deploy!

```bash
# Push to main branch
git add .
git commit -m "Add CI/CD workflows"
git push origin main

# Or trigger manually
gh workflow run deploy-full.yml
```

---

## Workflow Triggers

### Automatic Triggers
- **Backend**: Pushes to `backend/**`
- **Frontend**: Pushes to `frontend/**`
- **Infrastructure**: Pushes to `infrastructure/**`
- **Full Pipeline**: Pushes to `main`

### Manual Triggers
All workflows have `workflow_dispatch` - trigger from GitHub Actions tab.

---

## Deployment Flow

```
1. Push to main
   ↓
2. Run tests (backend + frontend)
   ↓
3. Deploy infrastructure (DB, Redis, Vault)
   ↓
4. Build & push Docker images
   ↓
5. Deploy with Helm
   ↓
6. Health checks
   ↓
7. Display URLs
```

---

## Monitoring Deployment

### Check Workflow Status
```bash
gh run list
gh run watch
```

### Check GKE Pods
```bash
kubectl get pods -n passchain
kubectl logs -f deployment/passchain-backend -n passchain
```

### Get Service URLs
```bash
kubectl get services -n passchain

# Frontend URL
echo "http://$(kubectl get service passchain-frontend -n passchain -o jsonpath='{.status.loadBalancer.ingress[0].ip}')"

# Backend URL
echo "http://$(kubectl get service passchain-backend -n passchain -o jsonpath='{.status.loadBalancer.ingress[0].ip}')"
```

---

## Rollback

```bash
# Rollback backend
helm rollback passchain-backend -n passchain

# Rollback frontend
helm rollback passchain-frontend -n passchain

# Or use kubectl
kubectl rollout undo deployment/passchain-backend -n passchain
```

---

## Cost Optimization

### Development
```bash
# Scale down when not in use
kubectl scale deployment --all --replicas=0 -n passchain

# Or delete cluster
gcloud container clusters delete $CLUSTER_NAME --zone=$REGION-a
```

### Production
- Use Preemptible nodes for non-critical workloads
- Enable cluster autoscaler
- Use committed use discounts
- Set up budget alerts

---

## Security Checklist

- ✅ Service account with minimal permissions
- ✅ Secrets stored in GitHub Secrets (not code)
- ✅ Network policies enabled
- ✅ Pod security policies enabled
- ✅ Private GKE cluster (optional)
- ✅ Workload Identity (recommended)
- ✅ Binary Authorization (optional)

---

## Troubleshooting

### Image pull errors
```bash
# Check GCR permissions
gcloud projects get-iam-policy $PROJECT_ID

# Ensure service account has storage.admin role
```

### Helm deployment fails
```bash
# Debug Helm
helm list -n passchain
helm get values passchain -n passchain

# Check pod events
kubectl describe pod <pod-name> -n passchain
```

### Health check fails
```bash
# Check pod logs
kubectl logs -f deployment/passchain-backend -n passchain

# Check service endpoints
kubectl get endpoints -n passchain
```

---

**Your CI/CD is ready! Push to main and watch it deploy!** 🚀

**AUUUUFFFF!** 🔥

