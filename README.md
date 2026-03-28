# PDS OAuth

This is an example for OAuth with ATProto.

## Docker Compose

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/)

### Getting Started

```bash
docker compose up -d
```

Caddy listens on port 80/443, proxies `/api/*` and `/oauth/*` to the backend, and serves the frontend as a SPA for all other routes.

### Custom Domain

By default the site address is `localhost`. Set the `SITE_ADDRESS` environment variable to use a custom domain (Caddy will automatically provision HTTPS certificates):

```bash
SITE_ADDRESS=example.com docker compose up -d
```

### Configuration

Backend config files are in `configs/` and are copied into the image at build time. To override at runtime, mount a volume in `docker-compose.yml`:

```yaml
services:
  backend:
    volumes:
      - ./configs:/etc/pds-oauth:ro
```

### Common Commands

```bash
# View logs
docker compose logs -f

# Stop services
docker compose down

# Rebuild images
docker compose up -d --build
```
