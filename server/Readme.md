# Flamingo Tutorial "Example Helloworld"

## Quickie: Run the (final) example

```bash
(cd frontend && npx flamingo-carotene build)
STYX_SERVICES=localhost go run main.go serve
```
Open http://localhost:3322 to access the example application.

## Docker

Docker build is located at `ghcr.io/saitho/styx:master`.
Run it with `docker run -p 3322:3322 ghcr.io/saitho/styx:master`