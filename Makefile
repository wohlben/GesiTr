.PHONY: build build-web build-go dev dev-api dev-web docker clean dolt generate seed test test-go test-web test-e2e lint lint-go lint-web update-screenshots update-screenshots-web update-screenshots-e2e

# Generate TypeScript types from Go structs
generate:
	go run github.com/gzuidhof/tygo@latest generate

# Full production build: generate types → Angular → Go binary with embedded SPA
build: generate build-web build-go

build-web:
	cd web && npm install && npx ng build --configuration=production

build-go: build-web
	go build -o gesitr .

# Development: start Go API + Angular dev server together
dev: generate
	cd web && node dev-server.mjs

dev-api:
	DEV=true AUTH_FALLBACK_USER=anon go run .

dev-web: generate
	cd web && CI=true npx ng serve

dolt:
	dolt sql-server --host 127.0.0.1 --port 3307 --data-dir .beads/dolt

# Lint and format check
lint: lint-go lint-web

lint-go:
	@test -z "$$(gofmt -l .)" || (echo "Go files not formatted:" && gofmt -l . && exit 1)

lint-web:
	cd web && npm run lint && npm run format:check

# Run all tests
test: lint test-go test-web test-e2e

test-go:
	go test ./...

test-web:
	cd web && npm test

# E2E uses dedicated ports to avoid conflicts with dev servers (8080/4200).
# API on :9876, Angular on :4300 (configured in angular.json serve:e2e + proxy.e2e.conf.json).
E2E_API_PORT := 9876

# E2E tests: builds Go binary, starts API on dedicated port, runs ng e2e (which manages
# its own Angular dev server via the "isolated" configuration), then cleans up.
test-e2e:
	@rm -f gesitr-e2e.db .gesitr-e2e
	@go build -o .gesitr-e2e .
	@echo "Starting API server on :$(E2E_API_PORT)..."
	@DATABASE_PATH=gesitr-e2e.db DEV=true AUTH_FALLBACK_USER=anon PORT=$(E2E_API_PORT) ./.gesitr-e2e & \
		API_PID=$$!; \
		for i in 1 2 3 4 5 6 7 8 9 10; do \
			curl -sf http://localhost:$(E2E_API_PORT)/api/exercises > /dev/null 2>&1 && break; \
			sleep 1; \
		done; \
		cd web && npx ng e2e --configuration=isolated; \
		TEST_EXIT=$$?; \
		kill $$API_PID 2>/dev/null; wait $$API_PID 2>/dev/null; \
		rm -f ../gesitr-e2e.db ../.gesitr-e2e; \
		exit $$TEST_EXIT

# Update screenshot baselines — same server setup as test-e2e.
update-screenshots: update-screenshots-web update-screenshots-e2e

update-screenshots-web:
	find web/src/app -path '*__screenshots__*' -name '*.png' -delete
	cd web && npx ng run web:test-screenshot || true
	cd web && npx ng run web:test-screenshot

update-screenshots-e2e:
	@rm -f gesitr-e2e.db .gesitr-e2e
	@go build -o .gesitr-e2e .
	@echo "Starting API server on :$(E2E_API_PORT)..."
	@DATABASE_PATH=gesitr-e2e.db DEV=true AUTH_FALLBACK_USER=anon PORT=$(E2E_API_PORT) ./.gesitr-e2e & \
		API_PID=$$!; \
		for i in 1 2 3 4 5 6 7 8 9 10; do \
			curl -sf http://localhost:$(E2E_API_PORT)/api/exercises > /dev/null 2>&1 && break; \
			sleep 1; \
		done; \
		cd web && npx ng e2e --configuration=isolated --update-snapshots; \
		TEST_EXIT=$$?; \
		kill $$API_PID 2>/dev/null; wait $$API_PID 2>/dev/null; \
		rm -f ../gesitr-e2e.db ../.gesitr-e2e; \
		exit $$TEST_EXIT

docker:
	docker build -t gesitr .

seed:
	rm -f gesitr.db
	go run ./cmd/seed

clean:
	rm -f gesitr gesitr.db
	rm -rf web/dist
	rm -f web/src/app/generated/models.ts
