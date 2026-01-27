# Description

This is an example project for my [blog](https://sh4nnongoh.github.io/blog/4-csrf-magic-links/)!

# Quick Start

```
go get

go tool air
```

# Docker

```
docker build -t myapp:latest .

# Without swap (recommended for containers)
# Same as memory = no swap allowed
# Disable swapping to disk
docker run -d \
  --name myapp \
  --memory=64m \
  --memory-swap=64m \
  --memory-swappiness=0 \
  --cpus="0.5" \
  -p 8080:8080 \
  -p 6060:6060 \
  -e GIN_MODE=release \
  myapp:latest
```

# Load Testing

```bash
cd ./cmd/loadtest/
./start.sh

# Evaluate the respective pprof output
go tool pprof -http=:8081 ./<metric>.pprof
```
