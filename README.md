# Kai Migration Demo - Kubernetes API Upgrade

Minimal Go project demonstrating Kai's ability to migrate deprecated Kubernetes APIs from v1.22 to v1.33.

## APIs Used (Deprecated in 1.22, Removed in 1.25)

- `autoscaling/v2beta1` HorizontalPodAutoscaler → `autoscaling/v2`
- `batch/v1beta1` CronJob → `batch/v1`
- `events/v1beta1` Event → `events/v1`

## Usage

```bash
# Install dependencies
go mod tidy

# Run the application
go run main.go

# Apply manifests (requires k8s cluster)
kubectl apply -f manifests/
```

## Migration Process

1. **Pre-migration**: App uses deprecated APIs from k8s 1.22
2. **Kai Analysis**: Scan for deprecated API usage
3. **Migration**: Update to k8s 1.33 compatible APIs
4. **Validation**: Verify functionality preserved

## Files

- `main.go` - Go app using deprecated client-go APIs
- `manifests/` - YAML examples with deprecated APIs
- `go.mod` - Dependencies locked to k8s 1.25 (last version supporting deprecated APIs)

## Expected Kai Changes

- `autoscaling/v2beta1` → `autoscaling/v2`
- `batch/v1beta1` → `batch/v1`  
- `events/v1beta1` → `events/v1`
- Client imports updated accordingly