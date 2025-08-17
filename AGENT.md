# AGENT.md - k8s-workshop Project Guide

## Project Overview
This is a Kubernetes workshop project that combines a Hugo static website with a
Go web server. The website content is in `web/` and served by a Go binary built
from `cmd/k8s-workshop`. The project serves as basis for a course teaching
Kubernetes. The course has emphasis on simplicity and hand-on learning through
exercises.

## Common Commands

### Build & Run
- **Build everything**: `task build` (builds static files + Go binary)
- **Run the application**: `task run` (builds and runs the binary)
- **Build static files only**: `task static-build`
- **Run Hugo dev server**: `task static-run`
- **Build Go binary only**: `go build -o bin/k8s-workshop ./cmd/k8s-workshop`

### Testing & Validation
- **Test Go code**: `go test ./...`
- **Check Go modules**: `go mod tidy && go mod verify`
- **Lint (if available)**: Check for golangci-lint or similar tools

### Docker
- **Build Docker image**: `task docker`
- **Push Docker image**: `task docker-push`
- **Run Docker container**: `task docker-run`

## Project Structure
- `cmd/k8s-workshop/` - Go web server that serves the Hugo site
- `web/` - Hugo static site with workshop content
  - `web/content/` - Markdown content files
  - `web/layouts/` - Hugo templates
  - `web/static/` - Static assets
  - `web/public/` - Generated static files (build output)
- `exercises/` - Exercise files referenced by the website
- `bin/` - Built binaries
- `build/` - Build configuration (Docker, etc.)

## Development Workflow
1. Make changes to `web/content/` for website content
2. Make changes to `cmd/k8s-workshop/` for server logic
3. Use `task run` to build and test locally
4. Use `task static-run` for rapid website development with Hugo's live reload

## Dependencies
- **Go**: 1.24.3
- **Hugo**: Static site generator
- **Chi router**: go-chi/chi/v5 for HTTP routing
- **Task**: Task runner for build automation

## Docker Configuration
- **Registry**: registry.digitalocean.com/blacksails/systematic-k8s-workshop
- **Multi-arch**: Builds for linux/amd64 and linux/arm64
- **Port**: Application runs on port 8080

## Code Style
- Markdown is written with a text width of 80 characters
